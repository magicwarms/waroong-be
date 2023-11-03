package users

import (
	"errors"
	"waroong-be/apps/constants"
	"waroong-be/apps/users/model"
	"waroong-be/config"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(gormDB *gorm.DB) *userRepository {
	gormDB.AutoMigrate(&model.UserModel{})
	return &userRepository{
		db: gormDB,
	}
}

// Save is to save user data based on user input
func (userRepo *userRepository) Store(user *model.UserModel) error {
	// begin a transaction
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserByEmail is to get only one user data by Email
func (userRepo *userRepository) GetUserByEmail(email string) (*model.UserModel, error) {
	var user model.UserModel
	result := userRepo.db.Preload("UserType", func(db *gorm.DB) *gorm.DB {
		return db.Select("id")
	}).Preload("Profile", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_id", "first_name", "last_name")
	}).Select("users.id", "email", "password", "users.is_active", "user_type_id").Where("email = ?", email).Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &model.UserModel{}, result.Error
	}

	if result.Error != nil {
		return &model.UserModel{}, result.Error
	}

	return &user, nil
}

// GetAll is to get all user data
func (userRepo *userRepository) GetAllSuperadminUsers() ([]*model.UserModel, error) {
	var users []*model.UserModel

	results := userRepo.db.Preload("UserType", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "is_active")
	}).Preload("Profile", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_id", "first_name", "last_name", "phone")
	}).Where("is_active = ?", true).Where("user_type_id = ?", constants.SUPERADMIN_USER_ROLE).Order("users.created_at DESC").Select("users.id", "email", "users.is_active", "users.created_at", "user_type_id").Find(&users)

	if results.Error != nil {
		return nil, results.Error
	}

	return users, nil
}

// GetById is to get only one user data by ID
func (userRepo *userRepository) GetById(userId uint64) (*model.UserModel, error) {
	var user model.UserModel
	result := userRepo.db.Preload("UserType", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "is_active")
	}).Preload("Profile").Select("users.id", "email", "password", "users.is_active", "user_type_id", "created_at", "updated_at").Where("id = ?", userId).Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &model.UserModel{}, result.Error
	}

	if result.Error != nil {
		return &model.UserModel{}, result.Error
	}

	return &user, nil
}

// UpdateUserPassword is to update user data based on user input
func (userRepo *userRepository) UpdateUserPassword(userId uint64, password string) error {
	if err := userRepo.db.Select("password").Updates(&model.UserModel{ID: userId, Password: password}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateForgotPasswordUserToken is to update user data based on user input
func (userRepo *userRepository) UpdateForgotPasswordUserToken(userId uint64, token string) error {
	if err := userRepo.db.Select("forgot_password_token").Updates(&model.UserModel{ID: userId, ForgotPasswordToken: token}).Error; err != nil {
		return err
	}
	return nil
}

// Delete is to delete user based on user input
func (userRepo *userRepository) Delete(id uint64) (bool, error) {
	userModel := &model.UserModel{
		ID: id,
	}
	if err := userRepo.db.Delete(userModel).Error; err != nil {
		return false, err
	}
	return true, nil
}

// GetUserTokenByEmailAndToken is to get forgot password token and id by Email and Token
func (userRepo *userRepository) GetUserForgotPasswordToken(token string) (*model.UserModel, error) {
	var user model.UserModel
	result := userRepo.db.Select("id", "forgot_password_token").Where("forgot_password_token = ?", token).Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &model.UserModel{}, result.Error
	}

	if result.Error != nil {
		return &model.UserModel{}, result.Error
	}

	return &user, nil
}

// UpdateRemoveUserForgotPasswordToken is to update user data based on user input
func (userRepo *userRepository) UpdateRemoveUserForgotPasswordToken(userId uint64) error {
	if err := userRepo.db.Select("forgot_password_token").Updates(&model.UserModel{ID: userId, ForgotPasswordToken: ""}).Error; err != nil {
		return err
	}
	return nil
}
