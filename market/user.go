package market

import (
	"fmt"
)

//go:generate mockgen -destination=./mock/user_service.go  -package=mock . UserService

type ErrUserNotFound struct {
	UserID string
}

func (err *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with id %v not found", err.UserID)
}

func (err *ErrUserNotFound) Is(target error) bool {
	t, ok := target.(*ErrUserNotFound)
	if !ok {
		return false
	}
	return t.UserID == err.UserID || t.UserID == ""
}

// var ErrUserNotFound = errors.New("user not found")

// UserService represents a user data backend.
type UserService interface {
	User(id string) (*User, error)
}

type User struct {
	ID   string
	Name string
}
