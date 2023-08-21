// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ensure provides a set of functions that panic if the error is not nil.
package ensure

// NoError panics if the error is not nil.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

// Value return value or panics if the error is not nil.
func Value[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}

	return v
}

// Values returns values or panics if the error is not nil.
func Values[V, V2 any](v V, v2 V2, err error) (V, V2) {
	if err != nil {
		panic(err)
	}

	return v, v2
}
