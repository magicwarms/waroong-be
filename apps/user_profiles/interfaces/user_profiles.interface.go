package interfaces

import (
	"waroong-be/apps/user_profiles/entity"
	"waroong-be/apps/user_profiles/model"
)

// UserProfileService is an interface from which our api module can access our repository of all our models
type UserProfileService interface {
	UpdateUserProfile(user *entity.UpdateUserRequestDTO) error
}

// UserProfileRepository interface allows us to access the CRUD Operations in repositories.
type UserProfileRepository interface {
	Update(user *model.UserProfileModel) error
}
