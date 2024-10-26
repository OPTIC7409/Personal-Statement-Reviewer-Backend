package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"psr/database/queries"
	"psr/types/user"
	ReturnModule "psr/utils/helpful/return_module"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users/me/{id}", http.HandlerFunc(h.GetUser)).Methods("GET")
	router.HandleFunc("/users/me/{id}", http.HandlerFunc(h.UpdateUser)).Methods("PUT")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ReturnModule.SendResponse(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid user ID type", http.StatusBadRequest)
		return
	}

	userProfile, err := queries.GetUserProfile(userID)
	if err != nil {
		fmt.Printf("Error getting user profile: %v\n", err)
		ReturnModule.SendResponse(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	// Ensure the response is formatted correctly
	ReturnModule.SendResponse(w, userProfile, http.StatusOK)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(int)
	var userProfile user.UserProfile
	err := json.NewDecoder(r.Body).Decode(&userProfile)
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid user profile data", http.StatusBadRequest)
		return
	}

	preferencesJSON, err := json.Marshal(userProfile.Preferences)
	if err != nil {
		ReturnModule.SendResponse(w, "Failed to encode preferences", http.StatusInternalServerError)
		return
	}

	err = queries.UpdateUserProfile(userID, userProfile.Bio, userProfile.ProfilePictureURL, string(preferencesJSON))
	if err != nil {
		ReturnModule.SendResponse(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	ReturnModule.SendResponse(w, "User profile updated successfully", http.StatusOK)
}
