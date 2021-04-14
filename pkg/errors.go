package gobit

import (
	"errors"
)

// i18n error
var (
	ErrFailedBinding  = errors.New("errors.bind")
	ErrFailedValidate = errors.New("errors.validate")
	ErrInternalError  = errors.New("errors.internal")
	ErrDBError        = errors.New("errors.db")
)

// framework error
var (
	ErrReopenTransaction = errors.New("no nested transaction")
)
