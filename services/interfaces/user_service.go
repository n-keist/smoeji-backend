package interfaces

import "smoeji/domain"

type IUserService interface {
	GetUsers() ([]domain.User, error)
	CreateUser(domain.UserCreateRequest) (domain.User, error)
	LoginUser(domain.UserLoginRequest) (*domain.UserLoginResponse, error)
	RefreshToken(token string) (*domain.UserLoginResponse, error)
}
