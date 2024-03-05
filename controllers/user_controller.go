package controllers

import (
	"smoeji/domain"
	"smoeji/services/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	UserController struct {
		userService interfaces.IUserService `di.inject:"userService"`
		validator   *validator.Validate     `di.inject:"validator"`
	}
)

func (uc *UserController) GetUsers(ctx *fiber.Ctx) error {
	users, err := uc.userService.GetUsers()
	if err != nil {
		return err
	}
	return ctx.JSON(&users)
}

func (uc *UserController) CreateUser(ctx *fiber.Ctx) error {
	var createUser domain.UserCreateRequest
	err := ctx.BodyParser(&createUser)
	if err != nil {
		return err
	}
	err = uc.validator.Struct(&createUser)
	if err != nil {
		return err
	}
	user, err := uc.userService.CreateUser(createUser)
	if err != nil {
		return err
	}
	return ctx.JSON(&user)
}
