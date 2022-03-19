package apperror

import "net/http"

type ClientErrorType int

const (
	Common ClientErrorType = iota
	WrongAuthorization
	WrongRequest
	WrongAuthToken
	NoClient
	WrongOneTimeCode
)

// thread safe, cus nobody writes.
var errorsMap = map[ClientErrorType]ClientError{
	Common:             errSomethingWentWrong,
	WrongRequest:       errWrongRequest,
	WrongAuthorization: errWrongAuthorization,
	WrongAuthToken:     errWrongAuthToken,
	NoClient:           errNoClient,
	WrongOneTimeCode:   errWrongOneTimeCode,
}

var errWrongAuthorization = ClientError{
	Code:    http.StatusUnauthorized,
	Message: "wrong password or email",
}

var errSomethingWentWrong = ClientError{
	Code:    http.StatusBadRequest,
	Message: "something went wrong",
}

var errWrongRequest = ClientError{
	Code:    http.StatusBadRequest,
	Message: "wrong request",
}

var errWrongAuthToken = ClientError{
	Code:    http.StatusUnauthorized,
	Message: "wrong auth token",
}

var errNoClient = ClientError{
	Code:    http.StatusBadRequest,
	Message: "client doesn't exist",
}

var errWrongOneTimeCode = ClientError{
	Code:    http.StatusBadRequest,
	Message: "wrong one-time code",
}

type ClientError struct {
	Code       int
	Message    string
	ErrorToLog error
}

func NewClientError(errorType ClientErrorType, err error) ClientError {
	clientError, ok := errorsMap[errorType]
	if !ok {
		clientError = errorsMap[Common]
	}

	clientError.ErrorToLog = err

	return clientError
}

func (err ClientError) Error() string {
	return err.Message
}

func (err ClientError) Unwrap() error {
	return err.ErrorToLog
}
