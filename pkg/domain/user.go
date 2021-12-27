package domain

import (
	"time"

	"github.com/arangodb/go-driver"
)

type User struct {
	Id                       driver.DocumentID `json:"id,omitempty"`
	Status                   bool              `json:"status,omitempty"`
	FirstName                string            `json:"first_name,omitempty"`
	LastName                 string            `json:"last_name,omitempty"`
	Email                    string            `json:"email,omitempty"`
	Password                 string            `json:"password,omitempty"`
	LastLogin                *time.Time        `json:"last_login,omitempty"`
	CreatedAt                time.Time         `json:"created_at,omitempty"`
	PasswordResetRequestedAt *time.Time        `json:"password_reset_requested_at,omitempty"`
	LatestToken              string            `json:"latest_token,omitempty"`
	RefreshToken             string            `json:"refresh_token,omitempty"`
	PasswordResetToken       string            `json:"password_reset_token,omitempty"`
	ConfirmationToken        string            `json:"confirmation_token,omitempty"`
	// Subscriptions            []Subscription    `json:"subscriptions,omitempty"`
}

type UserRegisterDto struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=32"`
	LastName  string `json:"last_name" validate:"required,min=2,max=32"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,weak_password,min=6,max=32"`
	// Subscription Subscription `json:"subscription"`
}

type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserPasswordResetRequestDto struct {
	Email string `json:"email" validate:"required,email"`
}

type UserSetNewPasswordDto struct {
	Password        string `json:"password" validate:"required,weak_password,min=6,max=32"`
	PasswordConfirm string `json:"password_confirm" validate:"required,weak_password,min=6,max=32"`
}

type UserAuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
