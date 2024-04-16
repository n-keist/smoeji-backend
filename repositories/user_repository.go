package repositories

import (
	"context"
	"smoeji/domain"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

type UserRepository struct {
	database *ksql.DB `di.inject:"util::database"`
}

var usersTable = ksql.NewTable("users", "id")
var ctx = ksql.InjectLogger(context.Background(), ksql.Logger)

func (ur *UserRepository) GetUsers() ([]domain.User, error) {
	var users []domain.User

	err := ur.database.Query(ctx, &users, "FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) CreateUser(user domain.User) (domain.User, error) {
	err := ur.database.Insert(ctx, usersTable, &user)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	err := ur.database.QueryOne(ctx, &user, "FROM users WHERE email = $1;", email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (domain.User, error) {
	var user domain.User

	err := ur.database.QueryOne(ctx, &user, "FROM users WHERE id = $1", id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
