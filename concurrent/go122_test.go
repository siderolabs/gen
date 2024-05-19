// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.22 && !go1.24 && !nomaptypehash

//nolint:nlreturn,wsl,testpackage
package concurrent

import (
	"testing"
	"unsafe"
)

func TestHashTrieMap(t *testing.T) {
	testHashTrieMap(t, func() *HashTrieMap[string, int] {
		return NewHashTrieMap[string, int]()
	})
}

func TestHashTrieMapBadHash(t *testing.T) {
	testHashTrieMap(t, func() *HashTrieMap[string, int] {
		// Stub out the good hash function with a terrible one.
		// Everything should still work as expected.
		m := NewHashTrieMap[string, int]()
		m.keyHash = func(_ unsafe.Pointer, _ uintptr) uintptr {
			return 0
		}
		return m
	})
}
