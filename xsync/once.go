// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xsync contains the additions to std sync package.
package xsync

import "sync"

// Once is small wrapper around [sync.Once]. It stores the result inside.
type Once[T any] struct {
	val  T
	once sync.Once
}

// Do runs the function only once.
func (o *Once[T]) Do(fn func() T) T {
	o.once.Do(func() {
		o.val = fn()
	})

	return o.val
}
