package repositories

import (
	"errors"
	"smoeji/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB `di.inject:"util::database"`
}

func (ur *UserRepository) GetUsers() ([]domain.User, error) {
	var users []domain.User

	result := ur.database.Find(&users)
	if err := result.Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) CreateUser(user domain.User) (domain.User, error) {
	result := ur.database.Create(&user)

	if err := result.Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	result := ur.database.First(&user, "email = ?", email)
	if err := result.Error; err != nil {
		return domain.User{}, err
	}
	if result.RowsAffected == 0 {
		return domain.User{}, errors.New("no user!")
	}
	return user, nil
}
