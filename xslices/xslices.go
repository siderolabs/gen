// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xslices contains a utility functions to work with slices.
package xslices

// Map applies the function fn to each element of the slice and returns a new slice with the results.
func Map[T, R any](slc []T, fn func(T) R) []R {
	if len(slc) == 0 {
		return nil
	}

	r := make([]R, 0, len(slc))

	for _, v := range slc {
		r = append(r, fn(v))
	}

	return r
}

// FlatMap applies the function fn to each element of the slice and returns a new slice with the results.
// It flattens the result of fn into the result slice.
func FlatMap[T, R any](slc []T, fn func(T) []R) []R {
	if len(slc) == 0 {
		return nil
	}

	r := make([]R, 0, len(slc))

	for _, v := range slc {
		r = append(r, fn(v)...)
	}

	return r
}

// Filter returns a slice containing all the elements of s that satisfy fn.
func Filter[S ~[]T, T any](slc S, fn func(T) bool) S {
	// NOTE(DmitriyMV): We use type parameter S here to return exactly the same type as the input slice.
	if len(slc) == 0 {
		return nil
	}

	r := make(S, 0, len(slc))

	for _, v := range slc {
		if fn(v) {
			r = append(r, v)
		}
	}

	// No reason to return empty slice if we filtered everything out.
	if len(r) == 0 {
		return nil
	}

	return r[:len(r):len(r)]
}

// FilterInPlace filters the slice in place.
func FilterInPlace[S ~[]V, V any](slc S, fn func(V) bool) S {
	// NOTE(DmitriyMV): We use type parameter S here to return exactly the same type as the input slice.
	if len(slc) == 0 {
		// We return original empty slice instead of a nil slice unlike Filter function.
		return slc
	}

	r := slc[:0]

	for _, v := range slc {
		if fn(v) {
			r = append(r, v)
		}
	}

	// We return original slice even if we filtered everything out unlike Filter function.
	return r
}

// ToMap converts a slice to a map.
func ToMap[T any, K comparable, V any](slc []T, fn func(T) (K, V)) map[K]V {
	if len(slc) == 0 {
		return nil
	}

	r := make(map[K]V, len(slc))

	for _, v := range slc {
		key, val := fn(v)
		r[key] = val
	}

	return r
}

// ToSet converts a slice to a set.
func ToSet[T comparable](slc []T) map[T]struct{} {
	if len(slc) == 0 {
		return nil
	}

	r := make(map[T]struct{}, len(slc))

	for _, v := range slc {
		r[v] = struct{}{}
	}

	return r
}

// ToSetFunc converts a slice to a set using the function fn.
func ToSetFunc[T any, K comparable](slc []T, fn func(T) K) map[K]struct{} {
	if len(slc) == 0 {
		return nil
	}

	r := make(map[K]struct{}, len(slc))

	for _, v := range slc {
		r[fn(v)] = struct{}{}
	}

	return r
}

// CopyN copies first n elements. If n is greater than the length of the slice, it will copy the whole slice.
func CopyN[S ~[]V, V any](s S, n int) S {
	if s == nil {
		return nil
	}

	if n > len(s) {
		n = len(s)
	}

	result := make([]V, n)
	copy(result, s)

	return result
}

// Deduplicate removes duplicate elements from the slice based on the key function provided.
func Deduplicate[T any, K comparable](items []T, keyFunc func(T) K) []T {
	var (
		result []T
		seen   = make(map[K]struct{})
	)

	for _, item := range items {
		k := keyFunc(item)

		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}

			result = append(result, item)
		}
	}

	return result
}
