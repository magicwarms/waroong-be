package model

import (
	"fmt"
	"time"
	"waroong-be/config"

	"gorm.io/gorm"
)

// UserTypeModel Constructs your UserTypeModel under entities.
type UserTypeModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;unique" json:"name"`
	IsActive  *bool          `gorm:"type:boolean;default:true;not null" json:"is_active"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Set table name (GORM)
func (UserTypeModel) TableName() string {
	return "user_types"
}

// DEFINE HOOKS
func (user *UserTypeModel) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before create data", config.PrettyPrint(user))
	return
}

func (user *UserTypeModel) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("After create data", config.PrettyPrint(user))
	return
}

func (user *UserTypeModel) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Before update data", config.PrettyPrint(user))
	return
}

func (user *UserTypeModel) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("After update data", config.PrettyPrint(user))
	return
}

func (user *UserTypeModel) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("Before delete data", config.PrettyPrint(user))
	return
}

func (user *UserTypeModel) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("After delete data", config.PrettyPrint(user))
	return
}
