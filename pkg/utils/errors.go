package utils

import "errors"

// Controller errors
var (
	// ErrBadRequest is used when the request body is not valid format
	ErrBadRequestBody = errors.New("bad request body")

	// ErrInvalidCredentials is used when the user's credentials are invalid
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrInvalidTemplateID is used when the template id is invalid or not found
	ErrInvalidTemplateID = errors.New("invalid template id")

	// ErrDidntHavePermission is used when the user doesn't have permission to access or modify the resource
	ErrDidntHavePermission = errors.New("you didn't have permission to do this action")
)

// Service errors
var (
	// ErrFieldNotMatch is used when the field in the request body is not match with the field in the template that saved in the database
	ErrFieldNotMatch = errors.New("document fields doesn't match with template fields")

	// ErrAlreadyVerified is used when the document is already verified
	ErrAlreadyVerified = errors.New("already verified")

	// ErrNotVerifiedYet is used when the document is not verified yet
	ErrNotVerifiedYet = errors.New("not verified yet")

	//ErrAlreadySigned is used when the document is already signed
	ErrAlreadySigned = errors.New("already signed")
)

// Repository errors
var (
	// ErrUsernameAlreadyExist is used when the username is already exist in the database
	ErrUsernameAlreadyExist = errors.New("user with provided username already exist")

	// ErrNIKAlreadyExist is used when the NIK is already exist in the database
	ErrNIKAlreadyExist = errors.New("user with provided nik already exist")

	// ErrNIPAlreadyExist is used when the NIP is already exist in the database
	ErrNIPAlreadyExist = errors.New("user with provided nip already exist")

	// ErrUserNotFound is used when the user is not found in the database
	ErrUserNotFound = errors.New("user not found")

	// ErrTemplateNotFound is used when the template is not found in the database
	ErrTemplateNotFound = errors.New("template not found")

	// ErrTemplateFieldNotFound is used when the template field is not found in the database
	ErrTemplateFieldNotFound = errors.New("template field not found")

	// ErrDuplicateRegister is used when the document register is already exist in the database
	ErrDuplicateRegister = errors.New("document with provided register already exist")

	// ErrDocumentNotFound is used when the document is not found in the database
	ErrDocumentNotFound = errors.New("document not found")
)
