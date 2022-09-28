// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xerrors contains the additions to std errors package.
package xerrors

import (
	"errors"
	"fmt"
)

// TypeIs is wrapper around errors.As which check the error type.
func TypeIs[T error](err error) bool {
	var expected T

	return errors.As(err, &expected)
}

// TagIs is wrapper around errors.As which checks the error tag.
func TagIs[T Tag](err error) bool {
	var expected *Tagged[T]

	return errors.As(err, &expected)
}

// Tag is a type which can be used to tag Tagged errors.
type Tag interface {
	~struct{}
}

// Tagged is an error with a tag attached. Tag can only be an empty struct.
//
//nolint:errname
type Tagged[T Tag] struct {
	err error
}

// NewTagged creates a new typed error.
func NewTagged[T Tag](err error) error {
	return &Tagged[T]{err: err}
}

// NewTaggedf creates a new typed error.
func NewTaggedf[T Tag](format string, a ...any) error {
	return &Tagged[T]{err: fmt.Errorf(format, a...)}
}

// Error implements error interface.
func (e *Tagged[T]) Error() string {
	return e.err.Error()
}

// Unwrap implements errors.Unwrap.
func (e *Tagged[T]) Unwrap() error {
	return e.err
}
