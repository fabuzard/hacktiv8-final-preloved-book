package utils

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserForbidden = errors.New("user not eligible")
	ErrBadReq        = errors.New("request not valid")
	ErrUnauthorized  = errors.New("no credentials or wrong credentials")
)
