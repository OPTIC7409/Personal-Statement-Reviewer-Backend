package database

import (
	"context"
	"fmt"
	"psr/cmd/api/secrets"
	"psr/utils/helpful/discord"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var databaseConnection *pgxpool.Pool

const (
	PermPostNotification = 1 << iota
)

func InitializeDatabase() bool {
	encodedIP := secrets.GetEnvVariable("DATABASE_IP")
	encodedPassword := secrets.GetEnvVariable("DATABASE_PASSWORD")

	if encodedIP == "" || encodedPassword == "" {
		return false
	}

	connectionString := fmt.Sprintf("postgres://validate:%s@%s/validate", encodedPassword, encodedIP)

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		fmt.Println("Error parsing connection string:", err)
		return false
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Minute
	poolConfig.MaxConnIdleTime = 2 * time.Minute

	databaseConnection, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return false
	}

	fmt.Println("Connected to the database")
	return true
}

func CreateTables() bool {
	_, err := GetConnection().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS personal_statements (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS feedback (
			id SERIAL PRIMARY KEY,
			statement_id INT NOT NULL,
			feedback_text TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (statement_id) REFERENCES personal_statements (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS feedback (
			id SERIAL PRIMARY KEY,
			statement_id INT NOT NULL,
			feedback_text TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (statement_id) REFERENCES personal_statements (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS revisions (
			id SERIAL PRIMARY KEY,
			statement_id INT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (statement_id) REFERENCES personal_statements (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS user_profiles (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			bio TEXT,
			profile_picture_url VARCHAR(255),
			preferences JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			action VARCHAR(100) NOT NULL,
			details JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);
	`)

	if err != nil {
		discord.SendMessage(discord.ErrorLog, "Error creating tables: "+err.Error())
		fmt.Println("Error creating tables:", err)
		return false
	}

	fmt.Println("Created tables")
	return true
}

func GetConnection() *pgxpool.Pool {
	return databaseConnection
}
