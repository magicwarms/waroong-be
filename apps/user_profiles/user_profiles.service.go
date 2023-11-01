package user_profiles

import (
	"waroong-be/apps/user_profiles/entity"
	"waroong-be/apps/user_profiles/interfaces"
	"waroong-be/apps/user_profiles/model"
)

type userProfileService struct {
	userProfileRepository interfaces.UserProfileRepository
}

// NewService is used to create a single instance of the service
func NewService(r interfaces.UserProfileRepository) interfaces.UserProfileService {
	return &userProfileService{
		userProfileRepository: r,
	}
}

func (userProfile *userProfileService) UpdateUserProfile(user *entity.UpdateUserRequestDTO) error {
	// TODO
	// add validation for phone
	// check phone has been used in another user data or not
	// find user that only deleted_at null
	// start to update the data to database through repository
	errUpdate := userProfile.userProfileRepository.Update(&model.UserProfileModel{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	})

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}
