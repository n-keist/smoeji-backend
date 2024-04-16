package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	RefreshToken struct {
		ID        uuid.UUID `ksql:"id" json:"-"`
		UserID    uuid.UUID `ksql:"user_id" json:"-"`
		Value     string    `ksql:"value" json:"token"`
		ExpiresAt time.Time `ksql:"expires_at" json:"expiresAt"`
		CreatedAt time.Time `ksql:"created_at,timeNowUTC/skipUpdates" json:"createdAt"`
	}

	RefreshTokenRefreshRequest struct {
		Token string `json:"refreshToken"`
	}
)
