package users

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	profileModel "waroong-be/apps/user_profiles/model"
	userTypeInterfaces "waroong-be/apps/user_types/interfaces"
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/interfaces"
	"waroong-be/apps/users/model"
	"waroong-be/apps/utils"
	"waroong-be/config"

	"github.com/golang-jwt/jwt/v5"
)

type userService struct {
	userRepository  interfaces.UserRepository
	userTypeService userTypeInterfaces.UserTypeService
}

// NewService is used to create a single instance of the service
func NewService(r interfaces.UserRepository, userTypeService userTypeInterfaces.UserTypeService) interfaces.UserService {
	return &userService{
		userRepository:  r,
		userTypeService: userTypeService,
	}
}

// SaveBank is a service layer that helps insert user admin data to database
func (s *userService) StoreUser(user *entity.UserRequestDTO) error {

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
		return &entity.LoginUserResponse{}, errors.New("email didn't existed")
	}

	errPasswordIncorrect := utils.VerifyPassword(user.Password, userLogin.Password)
	if errPasswordIncorrect != nil {
		return &entity.LoginUserResponse{}, errors.New("password was incorrect")
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

func (s *userService) UpdateUserPassword(data *entity.ChangePasswordUserDTO) error {
	parseUserId, errParseUserId := strconv.ParseUint(data.UserID, 10, 64)
	if errParseUserId != nil {
		return errors.New(errParseUserId.Error())
	}

	getUser, err := s.userRepository.GetById(uint(parseUserId))
	if err != nil {
		return err
	}

	hashedPassword, errHashPassword := utils.HashPassword(data.Password)
	if errHashPassword != nil {
		return errHashPassword
	}
	newUserPassword := string(hashedPassword)

	fmt.Println(newUserPassword)

	errUpdatePassword := s.userRepository.UpdateUserPassword(getUser.ID, newUserPassword)

	if errUpdatePassword != nil {
		return err
	}

	return nil
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
