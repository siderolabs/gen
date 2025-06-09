// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package panicsafe_test

import (
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/panicsafe"
)

func TestRunNoPanic(t *testing.T) {
	ran := false

	err := panicsafe.Run(func() { ran = true })
	require.NoError(t, err)

	require.True(t, ran, "function should have run without panic")
}

func TestRunPanic(t *testing.T) {
	err := panicsafe.Run(func() {
		panic("test panic")
	})
	assert.True(t, panicsafe.IsPanic(err))

	assert.ErrorContains(t, err, "test panic")
	assertStackTrace(t, err)
}

func TestRunErr(t *testing.T) {
	testErr := errors.New("test err")

	err := panicsafe.RunErr(func() error {
		return testErr
	})
	assert.False(t, panicsafe.IsPanic(err))

	assert.Equal(t, testErr, err)
}

func TestRunErrPanic(t *testing.T) {
	err := panicsafe.RunErr(func() error {
		panic("test panic")
	})
	assert.True(t, panicsafe.IsPanic(err))

	assertStackTrace(t, err)
}

func TestPanicErrType(t *testing.T) {
	testErr := errors.New("test err")

	err := panicsafe.Run(func() {
		panic(testErr)
	})

	assert.True(t, panicsafe.IsPanic(err))
	assert.ErrorIs(t, err, testErr, "both panicsafe.ErrPanic and the original error should be wrapped")

	assertStackTrace(t, err)
}

func assertStackTrace(t *testing.T, err error) {
	file, callerFunc := trace(1) // skip this very function

	assert.ErrorContains(t, err, "runtime/debug.Stack()")
	assert.ErrorContains(t, err, callerFunc)
	assert.ErrorContains(t, err, file)
}

// trace returns the file and function name of the caller.
func trace(skip int) (file, function string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2+skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.File, frame.Function
}
