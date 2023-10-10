package entity

import "time"

type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsActive  *bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_ata"`
}

type LoginUserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}
