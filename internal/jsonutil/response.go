package jsonutil

// ErrorResponse struct.
type ErrorResponse struct {
	Error Error `json:"error"`
}

// Error struct.
type Error struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// SuccessfulResponse struct.
type SuccessfulResponse struct {
	Response interface{} `json:"response"`
}

// NewSuccessfulResponse create new SuccessfulResponse.
func NewSuccessfulResponse(response interface{}) SuccessfulResponse {
	return SuccessfulResponse{
		Response: response,
	}
}

// NewError create new ErrorResponse.
func NewError(code int, msg string) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			ErrorCode: code,
			ErrorMsg:  msg,
		},
	}
}
