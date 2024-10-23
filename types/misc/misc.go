package types

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
