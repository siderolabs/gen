// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build go1.22 && !go1.24 && !nolinkname

//nolint:revive
package maps

import "unsafe"

//go:linkname runtime_keys maps.keys
//go:noescape
func runtime_keys(m any, p unsafe.Pointer)

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}

	result := make([]K, 0, len(m))

	runtime_keys(m, unsafe.Pointer(&result))

	return result
}

//go:linkname runtime_values maps.values
//go:noescape
func runtime_values(m any, p unsafe.Pointer)

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}

	result := make([]V, 0, len(m))

	runtime_values(m, unsafe.Pointer(&result))

	return result
}
