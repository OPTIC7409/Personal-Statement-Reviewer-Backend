package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"psr/database/queries"
	"psr/utils/JWT"
	ReturnModule "psr/utils/helpful/return_module"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/signup", http.HandlerFunc(h.Signup)).Methods("POST")
	router.HandleFunc("/auth/login", http.HandlerFunc(h.Login)).Methods("POST")
	router.HandleFunc("/auth/account", http.HandlerFunc(h.GetAccountData)).Methods("GET")
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&userData)
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, err := queries.GetUserByEmail(userData.Email)
	if err == nil && existingUser.Email != "" {
		ReturnModule.SendResponse(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		ReturnModule.SendResponse(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	userID, err := queries.CreateUser(userData.Name, userData.Email, string(hashedPassword))
	if err != nil {
		ReturnModule.SendResponse(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	token, err := JWT.GenerateJWT(userData.Email)
	if err != nil {
		ReturnModule.SendResponse(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	queries.CreateSession(userID, token)

	ReturnModule.SendResponse(w, map[string]string{"token": token}, http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&loginData)
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	user, err := queries.GetUserByEmail(loginData.Email)
	if err != nil || user.Email == "" {
		fmt.Println(err)
		ReturnModule.SendResponse(w, "User not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password))
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := JWT.GenerateJWT(user.Email)
	if err != nil {
		ReturnModule.SendResponse(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	err = queries.CreateSession(user.ID, token)
	if err != nil {
		ReturnModule.SendResponse(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	ReturnModule.SendResponse(w, map[string]string{"token": token}, http.StatusOK)
}
func (h *Handler) GetAccountData(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		ReturnModule.SendResponse(w, "No token found in Authorization header", http.StatusUnauthorized)
		return
	}

	id, err := JWT.ValidateJWT(token)
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		ReturnModule.SendResponse(w, "Invalid user ID type", http.StatusBadRequest)
		return
	}
	user, err := queries.GetUserDetails(userID)
	if err != nil {
		ReturnModule.SendResponse(w, "User not found", http.StatusNotFound)
		return
	}

	ReturnModule.SendResponse(w, user, http.StatusOK)
}
