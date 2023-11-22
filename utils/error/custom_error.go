package error

import "errors"

var (
	ErrRequirdName   = errors.New("required name")
	ErrRequiredEmail = errors.New("required email")
	ErrTokenExpired  = errors.New("Token already expired")
	ErrEmptyPassword = errors.New("Ensure both fileds are token")
	ErrPasswordLen   = errors.New("Length of password must be atleast 6 characters")
	ErrPasswordMatch = errors.New("Password do not match")
	ErrUpdatePass    = errors.New("Error while updating password")
)
