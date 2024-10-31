// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xiter provides a set of iterator helpers.
package xiter

import (
	"iter"
)

// Concat returns an iterator over the concatenation of the sequences.
func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Concat2 returns an iterator over the concatenation of the sequences.
func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Equal reports whether the two sequences are equal.
func Equal[V comparable](x, y iter.Seq[V]) bool {
	next, stop := iter.Pull(y)
	defer stop()

	v2, ok := next()
	for v1 := range x {
		if !ok || v1 != v2 {
			return false
		}

		v2, ok = next()
	}

	return !ok
}

// Equal2 reports whether the two sequences are equal.
func Equal2[K, V comparable](x, y iter.Seq2[K, V]) bool {
	next, stop := iter.Pull2(y)
	defer stop()

	k2, v2, ok := next()
	for k1, v1 := range x {
		if !ok || k1 != k2 || v1 != v2 {
			return false
		}

		k2, v2, ok = next()
	}

	return !ok
}

// EqualFunc reports whether the two sequences are equal according to the function f.
func EqualFunc[V1, V2 any](f func(V1, V2) bool, x iter.Seq[V1], y iter.Seq[V2]) bool {
	next, stop := iter.Pull(y)
	defer stop()

	v2, ok := next()
	for v1 := range x {
		if !ok || !f(v1, v2) {
			return false
		}

		v2, ok = next()
	}

	return !ok
}

// EqualFunc2 reports whether the two sequences are equal according to the function f.
func EqualFunc2[K1, V1, K2, V2 any](f func(K1, V1, K2, V2) bool, x iter.Seq2[K1, V1], y iter.Seq2[K2, V2]) bool {
	next, stop := iter.Pull2(y)
	defer stop()

	k2, v2, ok := next()
	for k1, v1 := range x {
		if !ok || !f(k1, v1, k2, v2) {
			return false
		}

		k2, v2, ok = next()
	}

	return !ok
}

// Map returns an iterator over f applied to seq.
func Map[In, Out any](f func(In) Out, seq iter.Seq[In]) iter.Seq[Out] {
	return func(yield func(Out) bool) {
		for in := range seq {
			if !yield(f(in)) {
				return
			}
		}
	}
}

// Map2 returns an iterator over f applied to seq.
func Map2[KIn, VIn, KOut, VOut any](f func(KIn, VIn) (KOut, VOut), seq iter.Seq2[KIn, VIn]) iter.Seq2[KOut, VOut] {
	return func(yield func(KOut, VOut) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Filter returns an iterator over the elements in seq for which f returns true.
func Filter[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for e := range seq {
			if !f(e) {
				continue
			}

			if !yield(e) {
				return
			}
		}
	}
}

// Filter2 returns an iterator over the elements in seq for which f returns true.
func Filter2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !f(k, v) {
				continue
			}

			if !yield(k, v) {
				return
			}
		}
	}
}

// Keys returns an iterator over the keys in seq.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over the values in seq.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// ToSeq returns an iterator where each element is the result of applying fn to the elements in seq.
func ToSeq[K, V, R any](fn func(K, V) R, seq iter.Seq2[K, V]) iter.Seq[R] {
	return func(yield func(R) bool) {
		for k, v := range seq {
			if !yield(fn(k, v)) {
				return
			}
		}
	}
}

// ToSeq2 returns an iterator where each element is the result of applying fn to the elements in seq.
func ToSeq2[V1, R1, R2 any](fn func(V1) (R1, R2), seq iter.Seq[V1]) iter.Seq2[R1, R2] {
	return func(yield func(R1, R2) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Reduce applies f to the elements in seq, starting with the initial value.
func Reduce[V, R any](f func(R, V) R, sum R, seq iter.Seq[V]) R {
	for e := range seq {
		sum = f(sum, e)
	}

	return sum
}

// Reduce2 applies f to the elements in seq, starting with the initial value.
func Reduce2[K, V, R any](f func(R, K, V) R, sum R, seq iter.Seq2[K, V]) R {
	for k, v := range seq {
		sum = f(sum, k, v)
	}

	return sum
}

// Empty returns an empty iterator.
func Empty[V any](func(V) bool) {}

// Empty2 returns an empty iterator.
func Empty2[V, V2 any](func(V, V2) bool) {}

// Single returns an iterator over a single element.
func Single[V any](v V) iter.Seq[V] { return func(yield func(V) bool) { yield(v) } }

// Single2 returns an iterator over a single element.
func Single2[K, V any](k K, v V) iter.Seq2[K, V] { return func(yield func(K, V) bool) { yield(k, v) } }

// Find returns the first element in seq for which f returns true.
func Find[V any](f func(V) bool, seq iter.Seq[V]) (V, bool) {
	for e := range seq {
		if f(e) {
			return e, true
		}
	}

	return *new(V), false
}

// Find2 returns the first element in seq for which f returns true.
func Find2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) (K, V, bool) {
	for k, v := range seq {
		if f(k, v) {
			return k, v, true
		}
	}

	return *new(K), *new(V), false
}
