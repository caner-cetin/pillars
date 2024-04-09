package types

type Error struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// error codes
const (
	InvalidRequestBody = "INVALID_REQUEST_BODY"
)
