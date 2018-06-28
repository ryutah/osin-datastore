package datastore

import "errors"

// Error definitions
var (
	ErrEmptyClientID       = errors.New("ID field of Client is empty")
	ErrInvalidUserDataType = errors.New("UserData field must be string")
)
