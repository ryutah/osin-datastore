package datastore

import "errors"

// Error definitions
var (
	ErrInvalidUserDataType = errors.New("UserData field must be string")
)
