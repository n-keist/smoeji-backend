package main

import (
	"log"
	"smoeji/controllers"
	"smoeji/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/goioc/di"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())

	app.Get("/_monitor", monitor.New())

	app.Post("/login", di.GetInstance("authController").(*controllers.AuthController).Login)
	app.Get("/users",
		di.GetInstance("authMiddleware").(*middleware.AuthMiddleware).GetMiddleware(),
		di.GetInstance("userController").(*controllers.UserController).GetUsers,
	)
	app.Post("/users", di.GetInstance("userController").(*controllers.UserController).CreateUser)

	log.Fatal(app.Listen(":42069"))
}
