// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package panicsafe provides panic-handling function wrappers, helpful when spawning goroutines which should never panic.
package panicsafe

import (
	"errors"
	"fmt"
	"runtime/debug"
)

var errPanic = errors.New("goroutine panicked")

// IsPanic checks if the given error is a panic error.
func IsPanic(err error) bool {
	return errors.Is(err, errPanic)
}

// Run runs the given function, handling panics and converting them to errors.
func Run(f func()) error {
	return RunErrF(func() error {
		f()

		return nil
	})()
}

// RunErr runs the given error-returning function, handling panics and converting them to errors.
func RunErr(f func() error) error {
	return RunErrF(f)()
}

// RunErrF returns a function which wraps the given error-returning function, handling panics and converting them to errors.
func RunErrF(f func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()

				if rError, ok := r.(error); ok { // if the panic is an error, we wrap it as well
					err = fmt.Errorf("%w: %w\n%s", errPanic, rError, string(stack))

					return
				}

				err = fmt.Errorf("%w: %s\n%s", errPanic, r, string(stack))
			}
		}()

		return f()
	}
}
