package user_types

import (
	"errors"
	"waroong-be/apps/constants"
	"waroong-be/apps/user_types/model"
	"waroong-be/config"

	"gorm.io/gorm"
)

type userTypeRepository struct {
	db *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(gormDB *gorm.DB) *userTypeRepository {
	// execute seed data
	if config.GoDotEnvVariable("APPLICATION_ENV") == constants.DEV_ENV {
		seedUserType(gormDB)
	}

	return &userTypeRepository{
		db: gormDB,
	}
}

func seedUserType(gormDB *gorm.DB) {
	if errMigration := gormDB.AutoMigrate(&model.UserTypeModel{}); errMigration == nil && gormDB.Migrator().HasTable(&model.UserTypeModel{}) {
		userTypes := []*model.UserTypeModel{
			{
				Name: "superadmin",
			},
			{
				Name: "customer",
			},
			{
				Name: "merchant",
			},
		}
		for _, val := range userTypes {
			gormDB.Where(model.UserTypeModel{Name: val.Name}).FirstOrCreate(&model.UserTypeModel{Name: val.Name})
		}
	}
}

func (userTypeRepo *userTypeRepository) GetById(id uint) (*model.UserTypeModel, error) {
	var userType model.UserTypeModel
	result := userTypeRepo.db.First(&userType, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &model.UserTypeModel{}, nil
	}
	if result.Error != nil {
		return &model.UserTypeModel{}, result.Error
	}
	return &userType, nil
}
