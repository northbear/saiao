package errors

import stdErrors "errors"

var ErrUnauthorized = stdErrors.New("unauthorized")
var ErrNotFound = stdErrors.New("not found")
var ErrInvalidInput = stdErrors.New("invalid input")
var ErrExecutionFailed = stdErrors.New("execution failed")
var ErrTimeout = stdErrors.New("timeout")
var ErrInternal = stdErrors.New("internal error")
