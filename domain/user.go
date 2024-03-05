package domain

import "github.com/google/uuid"

type (
	User struct {
		ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
		Email         string         `gorm:"uniqueIndex;not null" json:"email"`
		Password      string         `gorm:"not null" json:"-"`
		RefreshTokens []RefreshToken `json:"-"`
	}

	UserCreateRequest struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,min=8" json:"password"`
	}

	UserLoginRequest struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,min=8" json:"password"`
	}

	UserLoginResponse struct {
		User         User   `json:"user"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
)
