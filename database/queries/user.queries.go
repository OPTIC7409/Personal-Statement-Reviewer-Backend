package queries

import (
	"psr/database"
	"psr/types/user"
)

func CreateUser(name, email, password string) error {
	_, err := database.GetConnection().Exec(`
		INSERT INTO users (name, email, password) VALUES ($1, $2, $3)
	`, name, email, password)
	return err
}

func GetUserByEmail(email string) (user.User, error) {
	var user user.User
	err := database.GetConnection().QueryRow(`
		SELECT * FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	return user, err
}

func GetUserProfile(userID int) (user.UserProfile, error) {
	var userProfile user.UserProfile
	err := database.GetConnection().QueryRow(`
		SELECT * FROM user_profiles WHERE user_id = $1
	`, userID).Scan(&userProfile.ID, &userProfile.UserID, &userProfile.Bio, &userProfile.ProfilePictureURL, &userProfile.Preferences)
	return userProfile, err
}

func UpdateUserProfile(userID int, bio string, profilePictureURL string, preferences string) error {
	_, err := database.GetConnection().Exec(`
		UPDATE user_profiles SET bio = $2, profile_picture_url = $3, preferences = $4 WHERE user_id = $1
	`, userID, bio, profilePictureURL, preferences)
	return err
}
