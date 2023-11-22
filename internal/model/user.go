package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

func NewUser(req *UserRequest) *User {
	return &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
	}
}

type UserResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

func (u *User) UserRes() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginToken struct {
	AccessToken string        `json:"access_token"`
	User        *UserResponse `json:"user"`
}

type Payload struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *Payload) Valid() error {
	// Check if the token has expired
	if time.Now().After(p.ExpiredAt) {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}

	// You can add more validation logic here if needed

	return nil
}

type ResetPasswordReq struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	gorm.Model
	Email string `json:"email"`
	Token string `json:"token"`
}

type ForgetPassword struct {
	NewPassword    string `json:"new_password"`
	RetypePassword string `json:"retype_password"`
	Token          string `json:"token"`
}
