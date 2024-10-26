package queries

import (
	"psr/database"
)

func GetUserIDFromToken(token string) (int, error) {
	var userID int
	err := database.GetConnection().QueryRow("SELECT user_id FROM sessions WHERE token = $1", token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
