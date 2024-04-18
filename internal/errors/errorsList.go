package errorslist

import "errors"

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSqlFailed         = "sql failed"
)
