package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID				       int64     `json:"id"`
	Username	       string    `json:"username" validate:"min=3,max=16,required"`
	Email    	       string    `json:"email" validate:"email,max=150,required"`
	Password	       string    `json:"password" validate:"min=8,max=32,required"`
	Avatar		       string    `json:"avatar" validate:"max=100"`
	CreatedAt	       string    `json:"created_at"`
	Activated        bool      `json:"activated"`
	ActivationLink   string    `json:"activation_link"`
	ResetToken       string    `json:"reset_token"`
	ResetTokenExpiry time.Time `json:"reset_token_expiry"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
