package controllers

import (
	"smoeji/domain"
	"smoeji/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userService *services.UserService `di.inject:"service::user"`
	validate    *validator.Validate   `di.inject:"util::validator"`
}

func (ac *AuthController) Register(ctx *fiber.Ctx) error {
	var createUser domain.UserCreateRequest
	err := ctx.BodyParser(&createUser)
	if err != nil {
		return err
	}
	err = ac.validate.Struct(&createUser)
	if err != nil {
		return err
	}
	user, err := ac.userService.CreateUser(createUser)
	if err != nil {
		return err
	}
	return ctx.JSON(&user)
}

func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest domain.UserLoginRequest
	err := ctx.BodyParser(&loginRequest)
	if err != nil {
		return err
	}
	err = ac.validate.Struct(&loginRequest)
	if err != nil {
		return err
	}
	loginResult, err := ac.userService.LoginUser(loginRequest)
	if err != nil {
		return err
	}
	return ctx.JSON(&loginResult)
}

func (ac *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var refreshRequest domain.RefreshTokenRefreshRequest
	err := ctx.BodyParser(&refreshRequest)
	if err != nil {
		return err
	}
	loginResult, err := ac.userService.RefreshToken(refreshRequest.Token)
	if err != nil {
		return err
	}
	return ctx.JSON(&loginResult)
}
