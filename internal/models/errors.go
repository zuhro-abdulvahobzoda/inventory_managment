package models

import "errors"

var (
	// auth / user errors
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrInvalidInput          = errors.New("username and password are required")
	ErrInvalidUsername       = errors.New("username must be at most 100 characters")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong       = errors.New("password must be at most 72 bytes")
	ErrInvalidToken          = errors.New("invalid token")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidRole           = errors.New("role must be either admin or employee")
	ErrAdminAlreadyExists    = errors.New("admin already exists")

	// item errors
	ErrItemNotFound    = errors.New("item not found")
	ErrInvalidItemName = errors.New("name is required and must be at most 200 characters")
	ErrInvalidQuantity = errors.New("quantity must be zero or greater")

	// location errors
	ErrLocationNotFound    = errors.New("location not found")
	ErrInvalidLocationName = errors.New("name is required and must be at most 100 characters")
	ErrInvalidCoordinates  = errors.New("latitude (-90..90) and longitude (-180..180) are required")
)
