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
