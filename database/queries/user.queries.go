package queries

import (
	"errors"
	"fmt"
	"psr/database"
	"psr/types/user"
)

func CreateUser(name, email, password string) (int, error) {
	var userID int
	err := database.GetConnection().QueryRow(`
		INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id
	`, name, email, password).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return userID, nil
}
func GetUserByEmail(email string) (user.User, error) {
	var user user.User
	err := database.GetConnection().QueryRow(`
		SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, fmt.Errorf("error retrieving user by email: %w", err)
	}
	return user, nil
}
func GetUserProfile(userID int) (user.UserProfile, error) {
	fmt.Println(userID)
	if userID <= 0 {
		return user.UserProfile{}, errors.New("user ID must be greater than zero")
	}
	var userProfile user.UserProfile
	query := `
		SELECT id, user_id, bio, profile_picture_url, preferences FROM user_profiles WHERE user_id = $1`
	err := database.GetConnection().QueryRow(query, userID).Scan(&userProfile.ID, &userProfile.UserID, &userProfile.Bio, &userProfile.ProfilePictureURL, &userProfile.Preferences)
	return userProfile, err
}

func UpdateUserProfile(userID int, bio string, profilePictureURL string, preferences string) error {
	_, err := database.GetConnection().Exec(`
		UPDATE user_profiles SET bio = $2, profile_picture_url = $3, preferences = $4 WHERE user_id = $1
	`, userID, bio, profilePictureURL, preferences)
	return err
}

func CreateUserProfile(userID int, bio string, profilePictureURL string, preferences string) (int, error) {
	var newUserID int
	err := database.GetConnection().QueryRow(`
		INSERT INTO user_profiles (user_id, bio, profile_picture_url, preferences) VALUES ($1, $2, $3, $4) RETURNING id
	`, userID, bio, profilePictureURL, preferences).Scan(&newUserID)
	return newUserID, err
}

func GetUserDetails(userID int) (user.User, error) {
	var user user.User
	err := database.GetConnection().QueryRow(`
		SELECT id, name, email FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Name, &user.Email)
	return user, err
}

func CreateSession(userID int, token string) error {
	_, err := database.GetConnection().Exec(`
		INSERT INTO sessions (user_id, token) VALUES ($1, $2)
	`, userID, token)
	return err
}
