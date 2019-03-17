package api

import "errors"

var (
	ErrTokenError = errors.New("Error parsing auth token")
)

type Auth interface {
	Valid(string) (bool, UserID, error)
	Sign(UserID) (string, error)
}
