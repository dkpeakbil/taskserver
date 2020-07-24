package domain

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user is exists")
	ErrUsernameHasTaken  = errors.New("username has taken before")
	ErrUserNotFound      = errors.New("user could not found")
)
