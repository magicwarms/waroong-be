package user_profiles

import (
	"waroong-be/apps/user_profiles/model"

	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(gormDB *gorm.DB) *userProfileRepository {
	gormDB.AutoMigrate(&model.UserProfileModel{})
	return &userProfileRepository{
		db: gormDB,
	}
}

func (userProfileRepo *userProfileRepository) Update(user *model.UserProfileModel) error {
	if err := userProfileRepo.db.Select("first_name", "last_name", "phone").Updates(&model.UserProfileModel{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Phone: user.Phone}).Error; err != nil {
		return err
	}
	return nil
}
