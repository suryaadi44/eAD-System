package utils

import "errors"

// Controller errors
var (
	ErrBadRequestBody       = errors.New("bad request body")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrInvalidTemplateID    = errors.New("invalid template id")
	ErrDocumentAccessDenied = errors.New("document access denied")
)

// Service errors
var (
	ErrFieldNotMatch = errors.New("document fields doesn't match with template fields")
)

// Repository errors
var (
	ErrUsernameAlreadyExist  = errors.New("user with provided username already exist")
	ErrNIKAlreadyExist       = errors.New("user with provided nik already exist")
	ErrNIPAlreadyExist       = errors.New("user with provided nip already exist")
	ErrUserNotFound          = errors.New("user not found")
	ErrTemplateNotFound      = errors.New("template not found")
	ErrTemplateFieldNotFound = errors.New("template field not found")
	ErrDuplicateRegister     = errors.New("document with provided register already exist")
	ErrDocumentNotFound      = errors.New("document not found")
)
