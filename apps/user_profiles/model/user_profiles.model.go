package model

import (
	"fmt"
	"time"
	"waroong-be/config"

	"gorm.io/gorm"
)

// UserProfileModel Constructs your UserProfileModel under entities.
type UserProfileModel struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	UserID    uint64 `gorm:"not null" json:"user_id"`
	FirstName string `gorm:"index:,not null" json:"first_name"`
	LastName  string `gorm:"index:,not null" json:"last_name"`
	// if you want to make it column not null, do below!
	Phone     string         `gorm:"index:,unique;not null;default:null" json:"phone"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Set tablename (GORM)
func (UserProfileModel) TableName() string {
	return "user_profiles"
}

// DEFINE HOOKS
func (userProfile *UserProfileModel) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before create data", config.PrettyPrint(userProfile))
	return
}

func (userProfile *UserProfileModel) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("After create data", config.PrettyPrint(userProfile))
	return
}

func (userProfile *UserProfileModel) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Before update data", config.PrettyPrint(userProfile))
	return
}

func (userProfile *UserProfileModel) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("After update data", config.PrettyPrint(userProfile))
	return
}

func (userProfile *UserProfileModel) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("Before delete data", config.PrettyPrint(userProfile))
	return
}

func (userProfile *UserProfileModel) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("After delete data", config.PrettyPrint(userProfile))
	return
}
