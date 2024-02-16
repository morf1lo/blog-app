package services

import "errors"

var (
	errInvalidCredenials	error	= errors.New("invalid credentials")
	errInternalServer 		error = errors.New("internal server error")
	errInvalidPassword		error	= errors.New("invalid password")
	errUserNotFound				error = errors.New("user not found")
	errPostNotFound				error = errors.New("post not found")
	errNoAccess						error = errors.New("you have no access")
)
