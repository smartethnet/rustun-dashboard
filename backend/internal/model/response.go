package model

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// ErrorResponseWithCode creates an error response with custom code
func ErrorResponseWithCode(code int, message string, err string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Error:   err,
	}
}

