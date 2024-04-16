package repositories

import (
	"errors"
	"smoeji/domain"
	"smoeji/util"
	"time"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

type RefreshTokenRepository struct {
	database *ksql.DB `di.inject:"util::database"`
}

var tokensTable = ksql.NewTable("refresh_tokens", "id")

//var ctx = ksql.InjectLogger(context.Background(), ksql.Logger)

func (repository *RefreshTokenRepository) CreateToken(user domain.User) (*domain.RefreshToken, error) {
	token := domain.RefreshToken{
		UserID:    user.ID,
		Value:     util.RandomString(48),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 14),
	}

	err := repository.database.Insert(ctx, tokensTable, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (repository *RefreshTokenRepository) GetTokenByValue(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := repository.database.QueryOne(ctx, &refreshToken, "FROM refresh_tokens WHERE value = $1", token)
	if err != nil {
		return &domain.RefreshToken{}, err
	}
	if !time.Now().Before(refreshToken.ExpiresAt) {
		return &domain.RefreshToken{}, errors.New("refresh token expired")
	}
	return &refreshToken, nil
}

func (repository *RefreshTokenRepository) DeleteToken(uuid uuid.UUID) error {
	return repository.database.Delete(ctx, tokensTable, uuid)
}
