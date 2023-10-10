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
