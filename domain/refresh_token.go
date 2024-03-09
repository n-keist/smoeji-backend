package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	RefreshToken struct {
		ID        uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"-"`
		UserID    uuid.UUID `gorm:"not null;index"`
		Value     string    `gorm:"not null"`
		CreatedAt time.Time `json:"createdAt"`
	}

	RefreshTokenRefreshRequest struct {
		Token string `json:"refreshToken"`
	}
)
