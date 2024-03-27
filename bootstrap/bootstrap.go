package bootstrap

import (
	"fmt"
	"log"
	"os"
	"reflect"
	db_plugins "smoeji/bootstrap/database/plugins"
	"smoeji/controllers"
	"smoeji/deps"
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

func Init() {
	loadDotEnv()
	connectPostgres()

	validate := validator.New()
	_, err := di.RegisterBeanInstance(deps.Util_Validator, validate)
	if err != nil {
		panic(err)
	}

	registerRepositories()
	registerServices()
	registerControllers()
	registerMiddlewares()

	// at this point, beans can no longer be registered
	if err := di.InitializeContainer(); err != nil {
		panic(err)
	}
}

// Register Services in DI
// Use [deps] as registration keys
func registerServices() {
	services := map[string]reflect.Type{
		deps.Service_User: reflect.TypeOf((*services.UserService)(nil)),
	}

	for key, instance := range services {
		registerReflect(key, instance)
	}
}

// Register Repositories in DI
// Use [deps] as registration keys
func registerRepositories() {
	repositories := map[string]reflect.Type{
		deps.Repository_User:         reflect.TypeOf((*repositories.UserRepository)(nil)),
		deps.Repository_RefreshToken: reflect.TypeOf((*repositories.RefreshTokenRepository)(nil)),
	}

	for key, instance := range repositories {
		registerReflect(key, instance)
	}
}

// Register Controllers in DI
// Use [deps] as registration keys
func registerControllers() {
	controllers := map[string]reflect.Type{
		deps.Controller_Auth: reflect.TypeOf((*controllers.AuthController)(nil)),
		deps.Controller_User: reflect.TypeOf((*controllers.UserController)(nil)),
	}

	for key, instance := range controllers {
		registerReflect(key, instance)
	}
}

// Register Middlewares in DI
// Use [deps] as registration keys
func registerMiddlewares() {
	controllers := map[string]reflect.Type{
		deps.Middleware_Auth: reflect.TypeOf((*middleware.AuthMiddleware)(nil)),
	}

	for key, instance := range controllers {
		registerReflect(key, instance)
	}
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		if os.Getenv("SMOEJI") != "true" {
			panic(err)
		}
	}
}

// Connects to Database configured from .env
// Automatically registered in DI
func connectPostgres() *gorm.DB {
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

	db.Use(&db_plugins.UuidDbPlugin{})

	db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})

	_, err = di.RegisterBeanInstance(deps.Util_Database, db)
	if err != nil {
		panic(err)
	}
	return db
}

// Registers a reflection Type in current DI Framework
// Handles Errors on its own
func registerReflect(key string, value reflect.Type) {
	o, err := di.RegisterBean(key, value)
	if err != nil {
		fmt.Printf("could not register %s -> %s", key, err)
	}
	if o {
		fmt.Printf("%s was overriden", key)
	}
}
