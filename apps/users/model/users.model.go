package model

import (
	"fmt"
	"time"
	userProfile "waroong-be/apps/user_profiles/model"
	userType "waroong-be/apps/user_types/model"
	"waroong-be/apps/utils"
	"waroong-be/config"

	"gorm.io/gorm"
)

// UserModel Constructs your UserModel under entities.
type UserModel struct {
	ID                  uint                         `gorm:"primaryKey" json:"id"`
	Email               string                       `gorm:"not null;unique" json:"email"`
	Password            string                       `gorm:"not null" json:"-"`
	IsActive            *bool                        `gorm:"type:boolean;default:true; not null" json:"is_active"`
	UserTypeID          uint                         `gorm:"not null" json:"user_type_id"`
	ForgotPasswordToken string                       `gorm:"null;unique" json:"-"`
	UserType            userType.UserTypeModel       `json:"user_type"`
	Profile             userProfile.UserProfileModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"profile"`
	CreatedAt           time.Time                    `json:"created_at"`
	UpdatedAt           time.Time                    `json:"updated_at"`
	DeletedAt           gorm.DeletedAt               `gorm:"index" json:"deleted_at"`
}

// Set tablename (GORM)
func (UserModel) TableName() string {
	return "users"
}

// DEFINE HOOKS
func (user *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before create data", config.PrettyPrint(user))
	hashedPassword, err := utils.HashPassword(user.Password)
	// if hash password error, return the error
	if err != nil {
		return err
	}
	// set password before create
	user.Password = string(hashedPassword)
	return
}

func (user *UserModel) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("After create data", config.PrettyPrint(user))
	return
}

func (user *UserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Before update data", config.PrettyPrint(user))
	return
}

func (user *UserModel) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("After update data", config.PrettyPrint(user))
	return
}

func (user *UserModel) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("Before delete data", config.PrettyPrint(user))
	return
}

func (user *UserModel) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("After delete data", config.PrettyPrint(user))
	return
}
