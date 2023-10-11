package interfaces

import (
	"waroong-be/apps/user_types/model"
)

// UserService is an interface from which our api module can access our repository of all our models
type UserTypeService interface {
	GetUserTypeById(id uint) (*model.UserTypeModel, error)
}

// UserRepository interface allows us to access the CRUD Operations in repositories.
type UserTypeRepository interface {
	GetById(id uint) (*model.UserTypeModel, error)
}
