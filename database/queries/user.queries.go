package queries

import (
	"context"
	"psr/database"
	"psr/types/user"
)

func CreateUser(name, email, password string) error {
	_, err := database.GetConnection().Exec(context.Background(), `
		INSERT INTO users (name, email, password) VALUES ($1, $2, $3)
	`, name, email, password)
	return err
}

func GetUserByEmail(email string) (user.User, error) {
	var user user.User
	err := database.GetConnection().QueryRow(context.Background(), `
		SELECT * FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	return user, err
}
