package repositories

import (
	"errors"
	"smoeji/domain"
	"smoeji/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	database *gorm.DB `di.inject:"util::database"`
}

func (repository *RefreshTokenRepository) CreateToken(user domain.User) (*domain.RefreshToken, error) {
	token := domain.RefreshToken{
		UserID: user.ID,
		Value:  util.RandomString(48),
	}

	result := repository.database.Create(&token)
	if err := result.Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (repository *RefreshTokenRepository) GetTokenByValue(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	result := repository.database.First(&refreshToken, "value = ?", token)
	if result.Error != nil {
		return &domain.RefreshToken{}, result.Error
	}
	if !time.Now().Before(refreshToken.CreatedAt.Add(time.Hour * 24 * 14)) {
		return &domain.RefreshToken{}, errors.New("refresh token expired")
	}
	return &refreshToken, nil
}

func (repository *RefreshTokenRepository) DeleteToken(uuid uuid.UUID) error {
	result := repository.database.Delete(&domain.RefreshToken{}, "id = ?", uuid)
	return result.Error
}
