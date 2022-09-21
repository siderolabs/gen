// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package optional contains optional type.
package optional

// Some creates a new Optional with provided value.
func Some[T any](val T) Optional[T] {
	return Optional[T]{
		value:   val,
		present: true,
	}
}

// None is helper function which returns empty Optional.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Optional represents a type T that may be "null" - that, which value is missing.
type Optional[T any] struct {
	value   T
	present bool
}

// Get returns a value, true on non-empty optional and zero-value and false otherwise.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// ValueOrZero returns a value on non-empty optional and zero-value otherwise.
func (o Optional[T]) ValueOrZero() T {
	return o.value
}

// ValueOr returns a value on non-empty optional or provided value otherwise.
func (o Optional[T]) ValueOr(val T) T {
	if o.present {
		return o.value
	}

	return val
}

// IsPresent tells of there any value inside the optional.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// Ptr returns a pointer to value or nil pointer. It is safe to change the value inside the pointer
// (if the T type supports this).
func (o Optional[T]) Ptr() *T {
	if o.present {
		store := o.value

		return &store
	}

	return nil
}
