package controllers

import (
	"smoeji/services/interfaces"

	"github.com/gofiber/fiber/v2"
)

type (
	UserController struct {
		userService interfaces.IUserService `di.inject:"service::user"`
	}
)

func (uc *UserController) GetUsers(ctx *fiber.Ctx) error {
	users, err := uc.userService.GetUsers()
	if err != nil {
		return err
	}
	return ctx.JSON(&users)
}
