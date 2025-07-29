package utils

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserForbidden = errors.New("user not eligible")
	ErrBadReq        = errors.New("request tidak valid")
	ErrUnauthorized  = errors.New("tidak ada atau salah kredensial")
)
