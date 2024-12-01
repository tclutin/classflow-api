package response

type APIError struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

func NewAPIError(error string) APIError {
	return APIError{Error: error}
}
