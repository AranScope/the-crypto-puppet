package errors

import "fmt"

type ErrorPrefix string

const (
	ErrBadRequest          ErrorPrefix = "bad_request"
	ErrUnauthorized        ErrorPrefix = "unauthorized"
	ErrForbidden           ErrorPrefix = "forbidden"
	ErrNotFound            ErrorPrefix = "not_found"
	ErrInternalServerError ErrorPrefix = "internal_server_error"
	ErrUnknown             ErrorPrefix = "unknown"
)

type Error struct {
	msg    string
	prefix ErrorPrefix
}

func New(prefix ErrorPrefix, msg string) *Error {
	return &Error{
		msg:    msg,
		prefix: prefix,
	}
}

func Augment(e Error, msg string) *Error {
	e.msg = msg
	return &e
}

func (e *Error) Error() string {
	if e.prefix != "" {
		return fmt.Sprintf("%s: %s", e.prefix, e.msg)
	}
	return e.msg
}

func PrefixMatches(e error, prefix ErrorPrefix) bool {
	ce, ok := e.(*Error)
	if !ok {
		return false
	}

	return ce.prefix == prefix
}
