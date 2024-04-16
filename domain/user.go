package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID        uuid.UUID `ksql:"id" json:"id"`
		Email     string    `ksql:"email" json:"email"`
		Password  string    `ksql:"password" json:"-"`
		UpdatedAt time.Time `ksql:"updated_at,timeNowUTC" json:"-"`
		CreatedAt time.Time `ksql:"created_at,timeNowUTC/skipUpdates" json:"-"`
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
