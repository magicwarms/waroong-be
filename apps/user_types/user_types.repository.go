package user_types

import (
	"waroong-be/apps/user_types/model"

	"gorm.io/gorm"
)

type userTypeRepository struct {
	db *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(gormDB *gorm.DB) *userTypeRepository {
	// execute seed data
	if errMigration := gormDB.AutoMigrate(&model.UserTypeModel{}); errMigration == nil && gormDB.Migrator().HasTable(&model.UserTypeModel{}) {
		userTypes := []*model.UserTypeModel{
			{
				Name: "superadmin",
			},
			{
				Name: "customer",
			},
		}
		for _, val := range userTypes {
			go gormDB.Where(model.UserTypeModel{Name: val.Name}).FirstOrCreate(&model.UserTypeModel{Name: val.Name})
		}
	}
	return &userTypeRepository{
		db: gormDB,
	}
}
