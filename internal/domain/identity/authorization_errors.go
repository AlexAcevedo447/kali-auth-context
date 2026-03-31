package identity

import "errors"

var (
	ErrAuthorizationRequestRequired  = errors.New("authorization request is required")
	ErrAuthorizationResourceRequired = errors.New("authorization resource is required")
	ErrAuthorizationActionRequired   = errors.New("authorization action is required")
)
