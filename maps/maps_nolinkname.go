// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build go1.24 || nolinkname

package maps

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}

	r := make([]K, 0, len(m))

	for k := range m {
		r = append(r, k)
	}

	return r
}

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}

	r := make([]V, 0, len(m))

	for _, v := range m {
		r = append(r, v)
	}

	return r
}
