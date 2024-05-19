// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//nolint:wsl,testpackage
package concurrent

import (
	"testing"
)

func BenchmarkHashTrieMapHasherLoadSmall(b *testing.B) {
	benchmarkHashTrieMapHasherLoad(b, testDataSmall[:])
}

func BenchmarkHashTrieMapHasherLoad(b *testing.B) {
	benchmarkHashTrieMapHasherLoad(b, testData[:])
}

func BenchmarkHashTrieMapHasherLoadLarge(b *testing.B) {
	benchmarkHashTrieMapHasherLoad(b, testDataLarge[:])
}

func benchmarkHashTrieMapHasherLoad(b *testing.B, data []string) {
	b.ReportAllocs()
	m := NewHashTrieMapHasher[string, int](stringHasher)
	for i := range data {
		m.LoadOrStore(data[i], i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = m.Load(data[i])
			i++
			if i >= len(data) {
				i = 0
			}
		}
	})
}

func BenchmarkHashTrieMapHasherLoadOrStore(b *testing.B) {
	benchmarkHashTrieMapHasherLoadOrStore(b, testData[:])
}

func BenchmarkHashTrieMapHasherLoadOrStoreLarge(b *testing.B) {
	benchmarkHashTrieMapHasherLoadOrStore(b, testDataLarge[:])
}

func benchmarkHashTrieMapHasherLoadOrStore(b *testing.B, data []string) {
	b.ReportAllocs()
	m := NewHashTrieMapHasher[string, int](stringHasher)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = m.LoadOrStore(data[i], i)
			i++
			if i >= len(data) {
				i = 0
			}
		}
	})
}
