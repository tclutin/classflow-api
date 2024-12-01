package response

type APIError struct {
	Error string `json:"error"`
}

func NewAPIError(error string) APIError {
	return APIError{Error: error}
}
