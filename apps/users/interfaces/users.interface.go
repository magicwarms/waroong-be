package interfaces

import (
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/model"
)

// UserService is an interface from which our api module can access our repository of all our models
type UserService interface {
	StoreUser(user *entity.AddUserRequestDTO) error
	UpdateUser(user *entity.UpdateUserRequestDTO) error
	FindAllSuperadminUsers() ([]*model.UserModel, error)
	LoginUser(user *entity.UserLoginRequestDTO) (*entity.LoginUserResponse, error)
	UpdateSuperadminPassword(user *entity.ChangePasswordUserDTO) error
	GetUserById(id uint64) (*model.UserModel, error)
	ForgotPassword(user *entity.ForgotPasswordRequestDTO) (bool, error)
	ChangeForgotPassword(user *entity.ChangeForgotPasswordRequestDTO) (bool, error)

	// DeleteUser(id uint) (bool, error)
}

// UserRepository interface allows us to access the CRUD Operations in repositories.
type UserRepository interface {
	Store(user *model.UserModel) error
	GetAllSuperadminUsers() ([]*model.UserModel, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetById(id uint64) (*model.UserModel, error)
	UpdateUserPassword(userId uint64, password string) error
	UpdateForgotPasswordUserToken(userId uint64, password string) error
	GetUserForgotPasswordToken(token string) (*model.UserModel, error)
	UpdateRemoveUserForgotPasswordToken(userId uint64) error

	// Delete(id uint) (bool, error)
}
