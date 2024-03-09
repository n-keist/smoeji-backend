package services

import (
	"errors"
	"smoeji/domain"
	"smoeji/repositories"
	"smoeji/services/interfaces"
	"smoeji/util"
)

type UserService struct {
	interfaces.IUserService
	userRepository         *repositories.UserRepository         `di.inject:"repository::user"`
	refreshTokenRepository *repositories.RefreshTokenRepository `di.inject:"repository::refreshToken"`
}

func (us *UserService) GetUsers() ([]domain.User, error) {
	return us.userRepository.GetUsers()
}

func (us *UserService) CreateUser(entity domain.UserCreateRequest) (domain.User, error) {
	hash, err := util.HashPassword(entity.Password)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		Email:    entity.Email,
		Password: hash,
	}
	return us.userRepository.CreateUser(user)
}

func (us *UserService) LoginUser(request domain.UserLoginRequest) (*domain.UserLoginResponse, error) {
	user, err := us.userRepository.GetUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if !util.VerifyPassword(user.Password, request.Password) {
		return nil, errors.New("login failed")
	}

	jwt, err := util.CreateJWT(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := us.refreshTokenRepository.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.UserLoginResponse{
		User:         user,
		Token:        jwt,
		RefreshToken: refreshToken.Value,
	}, nil
}

func (us *UserService) RefreshToken(token string) (*domain.UserLoginResponse, error) {
	refreshToken, err := us.refreshTokenRepository.GetTokenByValue(token)
	if err != nil {
		return &domain.UserLoginResponse{}, err
	}
	user, err := us.userRepository.GetUserById(refreshToken.UserID)
	if err != nil {
		return &domain.UserLoginResponse{}, err
	}
	newToken, err := us.refreshTokenRepository.CreateToken(user)
	if err != nil {
		return &domain.UserLoginResponse{}, err
	}
	jwt, err := util.CreateJWT(user)
	if err != nil {
		return nil, err
	}
	return &domain.UserLoginResponse{
		User:         user,
		Token:        jwt,
		RefreshToken: newToken.Value,
	}, nil
}
