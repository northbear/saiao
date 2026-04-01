package errors

import stdErrors "errors"

var ErrUnauthorized = stdErrors.New("unauthorized")
var ErrNotFound = stdErrors.New("not found")
var ErrInvalidInput = stdErrors.New("invalid input")
var ErrExecutionFailed = stdErrors.New("execution failed")
var ErrTimeout = stdErrors.New("timeout")
var ErrInternal = stdErrors.New("internal error")

func IsUnauthorized(err error) bool {
	return stdErrors.Is(err, ErrUnauthorized)
}

func IsNotFound(err error) bool {
	return stdErrors.Is(err, ErrNotFound)
}

func IsInvalidInput(err error) bool {
	return stdErrors.Is(err, ErrInvalidInput)
}

func IsExecutionFailed(err error) bool {
	return stdErrors.Is(err, ErrExecutionFailed)
}

func IsTimeout(err error) bool {
	return stdErrors.Is(err, ErrTimeout)
}

func IsInternal(err error) bool {
	return stdErrors.Is(err, ErrInternal)
}
