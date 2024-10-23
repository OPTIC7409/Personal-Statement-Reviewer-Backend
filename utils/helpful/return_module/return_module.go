package ReturnModule

import (
	"encoding/json"
	"fmt"
	"net/http"
	miscTypes "psr/types/misc"
	"psr/utils/helpful/discord"

	"github.com/pterm/pterm"
)

func SendResponse(w http.ResponseWriter, response any, statusCode int) {
	responseMarshal, err := json.Marshal(response)
	if err != nil {
		pterm.Error.Printf("Error converting %+v to json: %s\n", response, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch statusCode {
	case http.StatusOK:
		discord.SendMessage(discord.MessageLog, fmt.Sprintf("Request to %s was successful", w.Header().Get("X-Request-ID"))+" "+string(responseMarshal))
	default:
		discord.SendMessage(discord.ErrorLog, fmt.Sprintf("Request to %s was unsuccessful", w.Header().Get("X-Request-ID"))+" "+string(responseMarshal))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(responseMarshal)
	if err != nil {
		pterm.Error.Println("Error writing to HTTP: ", err)
	}
}

func CustomError(w http.ResponseWriter, errorMessage string, errorCode int) {
	errorResponse := miscTypes.ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	SendResponse(w, errorResponse, errorCode)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	errorResponse := miscTypes.ErrorResponse{
		ErrorCode:    http.StatusMethodNotAllowed,
		ErrorMessage: "That method is not accepted at this endpoint.",
	}
	SendResponse(w, errorResponse, http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	ErrorResponse := miscTypes.ErrorResponse{
		ErrorCode:    http.StatusInternalServerError,
		ErrorMessage: "There was an internal server error while trying to handle your request.",
	}
	SendResponse(w, ErrorResponse, http.StatusInternalServerError)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request) {
	SuccessResponse := miscTypes.SuccessResponse{
		Success: true,
	}
	SendResponse(w, SuccessResponse, http.StatusOK)
}
