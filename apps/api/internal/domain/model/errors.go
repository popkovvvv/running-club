package model

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrConflict           = errors.New("conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
	ErrWeakPassword       = errors.New("password too weak")
	ErrInvalidInviteCode  = errors.New("invalid invite code")
	ErrAlreadyMember      = errors.New("already a club member")
	ErrNotMember          = errors.New("not a club member")
	ErrInvalidRole        = errors.New("invalid role")
	ErrAlreadySignedUp    = errors.New("already signed up")
	ErrNotSignedUp        = errors.New("not signed up")
)
