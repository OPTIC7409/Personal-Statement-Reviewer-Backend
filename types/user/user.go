package user

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	Bio               string    `json:"bio,omitempty"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty"`
	Preferences       string    `json:"preferences,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
