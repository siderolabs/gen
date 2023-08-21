// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package slices contains a utility functions to work with slices.
package slices

import (
	"slices"

	"github.com/siderolabs/gen/xslices"
)

// NOTE(DmitriyMV): I tried to implement this generic functions to be as performant as possible.
// However, I couldn't find a way to do it, since Go (1.18 at the time of writing) cannot inline closures if (generic)
// function, which accepts the closure, was not inlined itself.
// And inlining budget of 80 is quite small, since most of it is going towards closure call.
// Somewhat relevant issue: https://github.com/golang/go/issues/41988

// Map applies the function fn to each element of the slice and returns a new slice with the results.
//
// Deprecated: Use [xslices.Map] instead.
func Map[T, R any](slc []T, fn func(T) R) []R { return xslices.Map(slc, fn) }

// FlatMap applies the function fn to each element of the slice and returns a new slice with the results.
// It flattens the result of fn into the result slice.
//
// Deprecated: Use [xslices.FlatMap] instead.
func FlatMap[T, R any](slc []T, fn func(T) []R) []R { return xslices.FlatMap(slc, fn) }

// Filter returns a slice containing all the elements of s that satisfy fn.
//
// Deprecated: Use [xslices.Filter] instead.
func Filter[S ~[]T, T any](slc S, fn func(T) bool) S { return xslices.Filter(slc, fn) }

// FilterInPlace filters the slice in place.
//
// Deprecated: Use [xslices.FilterInPlace] instead.
func FilterInPlace[S ~[]V, V any](slc S, fn func(V) bool) S { return xslices.FilterInPlace(slc, fn) }

// ToMap converts a slice to a map.
//
// Deprecated: Use [xslices.ToMap] instead.
func ToMap[T any, K comparable, V any](slc []T, fn func(T) (K, V)) map[K]V {
	return xslices.ToMap(slc, fn)
}

// ToSet converts a slice to a set.
//
// Deprecated: Use [xslices.ToSet] instead.
func ToSet[T comparable](slc []T) map[T]struct{} { return xslices.ToSet(slc) }

// ToSetFunc converts a slice to a set using the function fn.
//
// Deprecated: Use [xslices.ToSetFunc] instead.
func ToSetFunc[T any, K comparable](slc []T, fn func(T) K) map[K]struct{} {
	return xslices.ToSetFunc(slc, fn)
}

// IndexFunc returns the first index satisfying fn(slc[i]),
// or -1 if none do.
//
// Deprecated: Use [slices.IndexFunc] instead.
func IndexFunc[T any](slc []T, fn func(T) bool) int { return slices.IndexFunc(slc, fn) }

// Contains reports whether v is present in s.
//
// Deprecated: Use [slices.ContainsFunc] instead.
func Contains[T any](s []T, fn func(T) bool) bool { return slices.ContainsFunc(s, fn) }

// Copy copies first n elements. If n is greater than the length of the slice, it will copy the whole slice.
//
// Deprecated: Use [xslices.CopyN] instead.
func Copy[S ~[]V, V any](s S, n int) S { return xslices.CopyN(s, n) }

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
//
// Deprecated: Use [slices.Clone] instead.
func Clone[S ~[]E, E any](s S) S { return slices.Clone(s) }

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
//
// Deprecated: Use [slices.Clip] instead.
func Clip[S ~[]E, E any](s S) S { return slices.Clip(s) }

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
//
// Deprecated: Use [slices.Grow] instead.
func Grow[S ~[]E, E any](s S, n int) S { return slices.Grow(s, n) }
