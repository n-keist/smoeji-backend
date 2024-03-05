package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"smoeji/controllers"
	"smoeji/domain"
	"smoeji/middleware"
	"smoeji/repositories"
	"smoeji/services"

	"github.com/go-playground/validator/v10"
	"github.com/goioc/di"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("could not connect to db ", err)
	}

	db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})

	validate := validator.New()

	_, _ = di.RegisterBeanInstance("database", db)
	_, _ = di.RegisterBeanInstance("validator", validate)

	_, _ = di.RegisterBean("userRepository", reflect.TypeOf((*repositories.UserRepository)(nil)))
	_, _ = di.RegisterBean("refreshTokenRepository", reflect.TypeOf((*repositories.RefreshTokenRepository)(nil)))

	_, _ = di.RegisterBean("userService", reflect.TypeOf((*services.UserService)(nil)))

	_, _ = di.RegisterBean("userController", reflect.TypeOf((*controllers.UserController)(nil)))
	_, _ = di.RegisterBean("authController", reflect.TypeOf((*controllers.AuthController)(nil)))

	_, _ = di.RegisterBean("authMiddleware", reflect.TypeOf((*middleware.AuthMiddleware)(nil)))
	_ = di.InitializeContainer()
}
