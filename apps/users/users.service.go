package users

import (
	"errors"
	"waroong-be/apps/constants"
	profileModel "waroong-be/apps/user_profiles/model"
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/interfaces"
	"waroong-be/apps/users/model"
)

type userService struct {
	userRepository interfaces.UserRepository
}

// NewService is used to create a single instance of the service
func NewService(r interfaces.UserRepository) interfaces.UserService {
	return &userService{
		userRepository: r,
	}
}

// SaveBank is a service layer that helps insert user admin data to database
func (s *userService) StoreUser(user *entity.UserRequestDTO) error {
	checkEmail, _ := s.userRepository.GetUserByEmail(user.Email)
	if checkEmail.Email != "" {
		return errors.New("email already existedsss")
	}

	// start to insert the data to database through repository
	errSave := s.userRepository.Save(&model.UserModel{
		Email:      user.Email,
		Password:   user.Password,
		UserTypeID: constants.SUPERADMIN_USER_ROLE,
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

// // StoreCustomer is a service layer that helps insert customer data to database
// func (s *userService) StoreCustomer(user *entity.UserRequestDTO) error {

// 	checkEmail, _ := s.userRepository.GetUserByEmail(user.Email, constants.CUSTOMER_USER_ROLE)
// 	if checkEmail.Email != "" {
// 		return errors.New("email already existed")
// 	}

// 	// start to insert the data to database through repository
// 	errSave := s.userRepository.Save(&model.UserModel{
// 		Email:      user.Email,
// 		Password:   user.Password,
// 		UserTypeID: constants.SUPERADMIN_USER_ROLE,
// 		Profile: profileModel.UserProfileModel{
// 			FirstName: user.FirstName,
// 			LastName:  user.LastName,
// 			Phone:     user.Phone,
// 		},
// 	})

// 	if errSave != nil {
// 		return errSave
// 	}

// 	return nil

// }

// FindAllSuperadminUsers is a service layer that helps fetch all banks in Banks table
func (s *userService) FindAllSuperadminUsers() ([]*model.UserModel, error) {
	getAllUsers, err := s.userRepository.GetAllSuperadminUsers()
	if err != nil {
		return nil, err
	}
	return getAllUsers, nil
}

// func (s *userService) LoginUser(userLogin *entity.UserLoginRequestDTO) (*entity.LoginUserResponse, error) {
// 	user, _ := s.userRepository.GetUserByEmail(userLogin.Email)
// 	if user.Email == "" {
// 		return &entity.LoginUserResponse{}, errors.New("email does not exists")
// 	}
// 	isPasswordCorrect := model.VerifyPassword(user.Password, userLogin.Password)
// 	if isPasswordCorrect != nil && isPasswordCorrect == bcrypt.ErrMismatchedHashAndPassword {
// 		return &entity.LoginUserResponse{}, errors.New("password is wrong")
// 	}

// 	// Create the Claims
// 	claims := jwt.MapClaims{
// 		"id":         user.ID,
// 		"userTypeId": user.UserType.ID,
// 		"email":      user.Email,
// 		"expiresAt":  time.Now().Local().Add(time.Hour * time.Duration(72)).Unix(), //3 days
// 	}
// 	// Create token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
// 	// Generate encoded token and send it as response.
// 	generatedToken, err := token.SignedString([]byte(config.GoDotEnvVariable("SECRET_KEY")))
// 	if err != nil {
// 		return &entity.LoginUserResponse{}, err
// 	}

// 	return &entity.LoginUserResponse{
// 		ID:        user.ID,
// 		Email:     user.Email,
// 		Token:     generatedToken,
// 		ExpiresAt: "3 days",
// 	}, nil

// }

// // GetBankById is a service layer that helps get book data
// func (s *bankService) GetBankById(id string) (*model.BankModel, error) {
// 	getBank, err := s.bankRepository.GetById(id)
// 	if err != nil {
// 		return &model.BankModel{}, err
// 	}
// 	return getBank, nil
// }

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
