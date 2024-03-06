package main

import (
	"log"
	"smoeji/bootstrap"
	"smoeji/controllers"
	"smoeji/deps"
	"smoeji/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/goioc/di"
)

func main() {
	bootstrap.Init()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())

	app.Get("/_monitor", monitor.New())

	app.Post("/login", di.GetInstance(deps.Controller_Auth).(*controllers.AuthController).Login)
	app.Get("/users",
		di.GetInstance(deps.Middleware_Auth).(*middleware.AuthMiddleware).GetMiddleware(),
		di.GetInstance(deps.Controller_User).(*controllers.UserController).GetUsers,
	)
	app.Post("/users", di.GetInstance(deps.Controller_User).(*controllers.UserController).CreateUser)

	log.Fatal(app.Listen(":42069"))
}
