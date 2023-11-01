package users

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	userProfileEntity "waroong-be/apps/user_profiles/entity"
	userProfileInterfaces "waroong-be/apps/user_profiles/interfaces"
	profileModel "waroong-be/apps/user_profiles/model"
	userTypeInterfaces "waroong-be/apps/user_types/interfaces"
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/interfaces"
	"waroong-be/apps/users/model"
	"waroong-be/apps/utils"
	"waroong-be/config"

	"github.com/golang-jwt/jwt/v5"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type userService struct {
	userRepository     interfaces.UserRepository
	userTypeService    userTypeInterfaces.UserTypeService
	userProfileService userProfileInterfaces.UserProfileService
}

// NewService is used to create a single instance of the service
func NewService(r interfaces.UserRepository, userTypeService userTypeInterfaces.UserTypeService, userProfileService userProfileInterfaces.UserProfileService) interfaces.UserService {
	return &userService{
		userRepository:     r,
		userTypeService:    userTypeService,
		userProfileService: userProfileService,
	}
}

// StoreUser is a service layer that helps insert user data to database
func (s *userService) StoreUser(user *entity.AddUserRequestDTO) error {

	checkEmail, _ := s.userRepository.GetUserByEmail(user.Email)
	if checkEmail.Email != "" {
		return errors.New("email already existed")
	}

	userTypeId, errParseUint := strconv.ParseUint(user.UserTypeId, 10, 64)
	if errParseUint != nil {
		return errors.New(errParseUint.Error())
	}

	// start to insert the data to database through repository
	errSave := s.userRepository.Store(&model.UserModel{
		Email:      user.Email,
		Password:   user.Password,
		UserTypeID: uint(userTypeId),
		Profile: profileModel.UserProfileModel{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		},
	})

	if errSave != nil {
		return errSave
	}

	return nil
}

// UpdateUser is a service layer that helps update user data to database
func (s *userService) UpdateUser(user *entity.UpdateUserRequestDTO) error {
	parseUserId, errParseUserId := strconv.ParseUint(user.ID, 10, 64)
	if errParseUserId != nil {
		return errParseUserId
	}

	checkUser, _ := s.userRepository.GetById(uint(parseUserId))
	if checkUser.Email == "" {
		return errors.New("the user doesn't exists")
	}

	errUpdate := s.userProfileService.UpdateUserProfile(&userProfileEntity.UpdateUserRequestDTO{
		ID:        uint(checkUser.Profile.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	})

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

// FindAllSuperadminUsers is a service layer that helps fetch all banks in Banks table
func (s *userService) FindAllSuperadminUsers() ([]*model.UserModel, error) {
	getAllUsers, err := s.userRepository.GetAllSuperadminUsers()
	if err != nil {
		return nil, err
	}
	return getAllUsers, nil
}

func (s *userService) LoginUser(userLogin *entity.UserLoginRequestDTO) (*entity.LoginUserResponse, error) {
	user, errGetUserEmail := s.userRepository.GetUserByEmail(userLogin.Email)

	if errGetUserEmail != nil {
		return &entity.LoginUserResponse{}, errors.New("the email didn't exist")
	}

	errPasswordIncorrect := utils.VerifyPassword(user.Password, userLogin.Password)
	if errPasswordIncorrect != nil {
		return &entity.LoginUserResponse{}, errors.New("password is incorrect")
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":         user.ID,
		"userTypeId": user.UserType.ID,
		"expiresAt":  time.Now().Local().Add(time.Hour * time.Duration(72)).Unix(), //3 days
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	// Generate encoded token and send it as response.
	generatedToken, err := token.SignedString([]byte(config.GoDotEnvVariable("SECRET_KEY")))
	if err != nil {
		return &entity.LoginUserResponse{}, err
	}

	userType, _ := s.userTypeService.GetUserTypeById(user.UserType.ID)

	return &entity.LoginUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Token:     generatedToken,
		ExpiresAt: "3 days",
		UserType:  userType,
	}, nil

}

// GetUserById is a service layer that helps get user by id
func (s *userService) GetUserById(userId uint) (*model.UserModel, error) {
	getUser, err := s.userRepository.GetById(userId)
	if err != nil {
		return &model.UserModel{}, err
	}
	return getUser, nil
}

func (s *userService) UpdateSuperadminPassword(data *entity.ChangePasswordUserDTO) error {
	parseUserId, errParseUserId := strconv.ParseUint(data.UserID, 10, 64)
	if errParseUserId != nil {
		return errParseUserId
	}

	getUser, errGetUser := s.userRepository.GetById(uint(parseUserId))
	if errGetUser != nil {
		return errGetUser
	}

	hashedPassword, errHashPassword := utils.HashPassword(data.Password)
	if errHashPassword != nil {
		return errHashPassword
	}
	newUserPassword := string(hashedPassword)

	errUpdatePassword := s.userRepository.UpdateUserPassword(getUser.ID, newUserPassword)

	if errUpdatePassword != nil {
		return errUpdatePassword
	}

	return nil
}

func (s *userService) ForgotPassword(userData *entity.ForgotPasswordRequestDTO) (bool, error) {
	user, errGetUserEmail := s.userRepository.GetUserByEmail(userData.Email)

	if errGetUserEmail != nil {
		return false, errors.New("the email didn't exist")
	}

	forgotPasswordToken := utils.GenerateRandomStringBytes(20)

	go s.userRepository.UpdateForgotPasswordUserToken(user.ID, forgotPasswordToken)

	templateData := entity.ForgotPasswordEmailParams{
		Name:  user.Profile.FirstName + " " + user.Profile.LastName,
		Token: forgotPasswordToken,
	}
	sendEmail, errorSendForgotEmail := sendForgotPasswordEmail(templateData, user)
	if errorSendForgotEmail != nil {
		return false, errorSendForgotEmail
	}

	if sendEmail {
		fmt.Println("======================EMAIL SUCCESSFULLY SENT TO " + user.Email + " ======================")
	}

	return true, nil

}

func sendForgotPasswordEmail(templateData entity.ForgotPasswordEmailParams, user *model.UserModel) (bool, error) {
	userfullName := user.Profile.FirstName + " " + user.Profile.LastName
	htmlFile, errHtmlFile := config.ParseHtmlTemplate("forgot_password.html", templateData)

	if errHtmlFile != nil {
		return false, errHtmlFile
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: config.GoDotEnvVariable("MJ_EMAIL_SENDER"),
				Name:  "No-Reply - Beyond Astro",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: user.Email,
					Name:  userfullName,
				},
			},
			Subject:  "Forgot Password -" + userfullName,
			HTMLPart: htmlFile,
		},
	}

	emailResponseStatus := config.SendEmail(messagesInfo)

	return emailResponseStatus == "success", nil
}

func (s *userService) ChangeForgotPassword(user *entity.ChangeForgotPasswordRequestDTO) (bool, error) {
	userToken, errUserToken := s.userRepository.GetUserForgotPasswordToken(user.Token)

	if errUserToken != nil || userToken == nil {
		return false, errors.New("token forgot password is incorrect")
	}

	hashedPassword, errHashPassword := utils.HashPassword(user.Password)
	if errHashPassword != nil {
		return false, errHashPassword
	}
	newUserPassword := string(hashedPassword)

	errUpdatePassword := s.userRepository.UpdateUserPassword(userToken.ID, newUserPassword)

	go s.userRepository.UpdateRemoveUserForgotPasswordToken(userToken.ID)

	if errUpdatePassword != nil {
		return false, errUpdatePassword
	}

	return true, nil
}

// // UpdateBank is a service layer that helps update bank data to database
// func (s *bankService) UpdateBank(bank *entity.BankUpdateRequestDTO) error {
// 	// check book data by name to validate
// 	result, _ := s.bankRepository.GetById(bank.ID)
// 	if result.ID == "" {
// 		return errors.New("ID not found")
// 	}

// 	// start to insert the data to database through repository
// 	errUpdate := s.bankRepository.Update(&model.BankModel{
// 		ID:            bank.ID,
// 		Name:          bank.Name,
// 		LicenseId:     bank.LicenseId,
// 		AccountName:   bank.AccountName,
// 		AccountNumber: bank.AccountNumber,
// 		IsActive:      bank.IsActive,
// 	})
// 	if errUpdate != nil {
// 		return errUpdate
// 	}
// 	return nil
// }

// // DeleteBank soft deletes user data in the database
// func (s *userService) DeleteUser(id uint) (bool, error) {
// 	// Check if user with given ID exists
// 	userToDelete, err := s.userRepository.GetById(id)
// 	if err != nil {
// 		return false, err
// 	}

// 	if userToDelete.ID < 1 {
// 		return false, errors.New("user not found")
// 	}

// 	// Soft delete the user from the database
// 	deleted, err := s.userRepository.Delete(id)
// 	if err != nil {
// 		return false, err
// 	}

// 	return deleted, nil
// }
