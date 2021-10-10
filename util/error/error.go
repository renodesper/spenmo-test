package error

import (
	"fmt"

	"github.com/go-errors/errors"
)

type (
	// Error ...
	Error struct {
		Status      int    `json:"status,omitempty"`
		Code        string `json:"code,omitempty"`
		Message     string `json:"message,omitempty"`
		Trace       string `json:"trace,omitempty"`
		RedirectURL string `json:"redirect_url,omitempty"`
		Meta        Meta   `json:"meta,omitempty"`
	}

	// Meta ...
	Meta map[string]interface{}
)

// Error ...
func (e Error) Error() string {
	return e.Message
}

// Wrap ...
func (e Error) Wrap(err error) Error {
	return Error{
		Status:  e.Status,
		Code:    e.Code,
		Message: err.Error(),
		Trace:   errors.Wrap(err, 1).ErrorStack(),
	}
}

// WrapWithURL ...
func (e Error) WrapWithURL(err error, url string) Error {
	return Error{
		Status:      e.Status,
		Code:        e.Code,
		Message:     err.Error(),
		Trace:       errors.Wrap(err, 1).ErrorStack(),
		RedirectURL: url,
	}
}

// ProduceStackTrace ...
func (e Error) ProduceStackTrace() Error {
	e.Trace = errors.Wrap(e, 1).ErrorStack()
	return e
}

// WithoutStackTrace ...
func (e Error) WithoutStackTrace() Error {
	e.Trace = ""
	return e
}

// AppendError ...
func (e Error) AppendError(err error) Error {
	return Error{
		Status:  e.Status,
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", e.Error(), err.Error()),
		Trace:   errors.Wrap(err, 1).ErrorStack(),
	}
}

// RemoveStackTrace ...
func (e *Error) RemoveStackTrace() {
	e.Trace = ""
}

// NewError ...
func NewError(status int, code string, err error) Error {
	return Error{
		Status:  status,
		Code:    code,
		Message: err.Error(),
		Trace:   errors.Wrap(err, 1).ErrorStack(),
	}
}

// NewErrorWithURL ...
func NewErrorWithURL(status int, code string, err error, url string) Error {
	e := NewError(status, code, err)
	e.RedirectURL = url
	return e
}
