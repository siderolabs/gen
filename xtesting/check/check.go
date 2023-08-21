// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package check provides a set of functions that can be used for error checking in table-driven test.
package check

import (
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xerrors"
	"github.com/siderolabs/gen/xtesting"
)

// Check is a function that checks for an error or absence of it.
type Check func(t xtesting.T, err error)

// NoError returns a function that checks that no error happened.
func NoError() Check {
	return func(t xtesting.T, err error) {
		require.NoError(t, err)
	}
}

// EqualError returns a function that checks for a specific error text.
func EqualError(msg string) Check {
	return func(t xtesting.T, err error) {
		require.EqualError(t, err, msg)
	}
}

// ErrorContains returns a function that checks that error message contains a specific text.
func ErrorContains(msg string) Check {
	return func(t xtesting.T, err error) {
		require.ErrorContains(t, err, msg)
	}
}

// ErrorRegexp returns a function that checks that error message matches a specific regexp.
func ErrorRegexp(rx any) Check {
	return func(t xtesting.T, err error) {
		require.Error(t, err)
		require.NotZero(t, rx)
		require.Regexp(t, rx, err.Error())
	}
}

// ErrorAs returns a function that checks that error can be converted to a specific type.
func ErrorAs[T error]() Check {
	var target T

	return func(t xtesting.T, err error) {
		require.Error(t, err)
		require.ErrorAs(t, err, &target)
	}
}

// ErrorIs returns a function that checks that error is equal to a specific error.
func ErrorIs(target error) Check {
	return func(t xtesting.T, err error) {
		require.Error(t, err)
		require.ErrorIs(t, err, target)
	}
}

// ErrorTagIs returns a function that checks that error has a specific tag.
func ErrorTagIs[T xerrors.Tag]() Check {
	return func(t xtesting.T, err error) {
		require.Error(t, err)
		require.True(t, xerrors.TagIs[T](err))
	}
}
