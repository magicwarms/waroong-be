package interfaces

import (
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/model"
)

// UserService is an interface from which our api module can access our repository of all our models
type UserService interface {
	StoreUser(user *entity.UserRequestDTO) error
	FindAllSuperadminUsers() ([]*model.UserModel, error)
	LoginUser(user *entity.UserLoginRequestDTO) (*entity.LoginUserResponse, error)
	UpdateUserPassword(user *entity.ChangePasswordUserDTO) error
	GetUserById(id uint) (*model.UserModel, error)

	// DeleteUser(id uint) (bool, error)
}

// UserRepository interface allows us to access the CRUD Operations in repositories.
type UserRepository interface {
	Store(user *model.UserModel) error
	GetAllSuperadminUsers() ([]*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetById(id uint) (*model.UserModel, error)
	UpdateUserPassword(userId uint, password string) error
	// Update(bank *model.UserModel) error

	// Delete(id uint) (bool, error)
}
