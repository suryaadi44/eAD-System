package utils

import "errors"

// Controller errors
var (
	ErrBadRequestBody = errors.New("bad request body")

	ErrInvalidCredentials = errors.New("invalid email or password")
)

// Service errors
var ()

// Repository errors
var (
	ErrUsernameAlreadyExist = errors.New("user with provided username already exist")
	ErrNIKAlreadyExist      = errors.New("user with provided nik already exist")
	ErrNIPAlreadyExist      = errors.New("user with provided nip already exist")
	ErrUserNotFound         = errors.New("user not found")
)
