// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ordered contains ordered Pair and Triple types.
package ordered

// Pair is two element tuple.
type Pair[T1, T2 any] struct {
	F1 T1
	F2 T2
}

// MakePair creates a new Pair.
func MakePair[T1, T2 any](v1 T1, v2 T2) Pair[T1, T2] {
	return Pair[T1, T2]{
		F1: v1,
		F2: v2,
	}
}

// Triple is three element tuple of ordered values.
type Triple[T1, T2, T3 any] struct {
	V1 T1
	V2 T2
	V3 T3
}

// MakeTriple creates a new Triple.
func MakeTriple[T1, T2, T3 any](v1 T1, v2 T2, v3 T3) Triple[T1, T2, T3] {
	return Triple[T1, T2, T3]{
		V1: v1,
		V2: v2,
		V3: v3,
	}
}
