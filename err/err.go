package e

import (
	"errors"
	"fmt"
)

var (
	ErrMsgIncorrectAuthType = errors.New("incorrect auth type")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidAuthType      = errors.New("invalid auth type")
	ErrRevokedToken         = errors.New("token revoked")
	ErrIncorrectFileName    = errors.New("incorrect filename")
	ErrNotFound             = errors.New("file not found")
	ErrAlreadyExist         = errors.New("already exist")
)

// Wraps error with new message
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
