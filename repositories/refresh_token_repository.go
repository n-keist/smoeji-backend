package repositories

import (
	"smoeji/domain"
	"smoeji/util"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	database *gorm.DB `di.inject:"database"`
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
