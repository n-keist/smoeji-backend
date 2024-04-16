package bootstrap

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"reflect"
	"smoeji/controllers"
	"smoeji/deps"
	"smoeji/middleware"
	"smoeji/repositories"
	"smoeji/services"

	"github.com/go-playground/validator/v10"
	"github.com/goioc/di"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/pressly/goose/v3"
	"github.com/vingarcia/ksql"
	kpgx "github.com/vingarcia/ksql/adapters/kpgx5"
)

func Init() {
	loadDotEnv()
	connectPostgres()
	connectNats()

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
func connectPostgres() *ksql.DB {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	ctx := ksql.InjectLogger(context.Background(), ksql.Logger)

	db, err := kpgx.New(ctx, connectionString, ksql.Config{})

	if err != nil {
		log.Fatalln("could not connect to db ", err)
	}

	migratePostgres()

	_, err = di.RegisterBeanInstance(deps.Util_Database, &db)
	if err != nil {
		panic(err)
	}
	return &db
}

func connectNats() *nats.Conn {
	natsConnection, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		panic(err)
	}
	_, err = di.RegisterBeanInstance(deps.Util_PubSub, natsConnection)
	if err != nil {
		panic(err)
	}
	return natsConnection
}

//go:embed database/migrations/*.sql
var embedMigrations embed.FS

func migratePostgres() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "database/migrations"); err != nil {
		panic(err)
	}
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
