// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package must provide a set of functions that fail test if the error is not nil.
package must

import (
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xtesting"
)

// Value returns a function that accepts parameter that implements [xtesting.T] and uses [require.NoError] to ensure there was no error.
func Value[V any](v V, err error) func(t xtesting.T) V {
	return func(t xtesting.T) V {
		require.NoError(t, err)

		return v
	}
}

// Values returns a function that accepts parameter that implements [xtesting.T] and uses [require.NoError] to ensure there was no error.
func Values[V, V2 any](v V, v2 V2, err error) func(t xtesting.T) (V, V2) {
	return func(t xtesting.T) (V, V2) {
		require.NoError(t, err)

		return v, v2
	}
}
