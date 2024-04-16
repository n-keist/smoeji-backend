package main

import (
	"fmt"
	"log"
	"os"
	"smoeji/bootstrap"
	"smoeji/controllers"
	"smoeji/deps"
	"smoeji/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/goioc/di"
	"github.com/nats-io/nats.go"
)

func main() {
	bootstrap.Init()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())

	authGroup := app.Group("/auth")
	authGroup.Post("/register", di.GetInstance(deps.Controller_Auth).(*controllers.AuthController).Register)
	authGroup.Post("/login", di.GetInstance(deps.Controller_Auth).(*controllers.AuthController).Login)
	authGroup.Post("/refresh",
		di.GetInstance(deps.Middleware_Auth).(*middleware.AuthMiddleware).GetMiddleware(),
		di.GetInstance(deps.Controller_Auth).(*controllers.AuthController).RefreshToken,
	)

	app.Get("/users",
		di.GetInstance(deps.Middleware_Auth).(*middleware.AuthMiddleware).GetMiddleware(),
		di.GetInstance(deps.Controller_User).(*controllers.UserController).GetUsers,
	)

	app.Get("/_healthy", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"healthy": true,
		})
	})
	natsConnection := di.GetInstance(deps.Util_PubSub).(*nats.Conn)
	go handleNatsMessages(natsConnection)

	log.Fatal(app.Listen(":3000"))
}

func handleNatsMessages(natsConnection *nats.Conn) {
	_, err := natsConnection.Subscribe(os.Getenv("NATS_PG_CRUD_SUB"), func(m *nats.Msg) {
		fmt.Println(string(m.Data))
	})
	if err != nil {
		panic(err)
	}
}
