// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//nolint:wsl_v5
package concurrent_test

import (
	"testing"

	cnc "github.com/siderolabs/gen/concurrent"
)

func BenchmarkHashTrieMapLoadSmall(b *testing.B) {
	benchmarkHashTrieMapLoad(b, testDataSmall[:])
}

func BenchmarkHashTrieMapLoad(b *testing.B) {
	benchmarkHashTrieMapLoad(b, testData[:])
}

func BenchmarkHashTrieMapLoadLarge(b *testing.B) {
	benchmarkHashTrieMapLoad(b, testDataLarge[:])
}

func benchmarkHashTrieMapLoad(b *testing.B, data []string) {
	b.ReportAllocs()
	var m cnc.HashTrieMap[string, int]
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

func BenchmarkHashTrieMapLoadOrStore(b *testing.B) {
	benchmarkHashTrieMapLoadOrStore(b, testData[:])
}

func BenchmarkHashTrieMapLoadOrStoreLarge(b *testing.B) {
	benchmarkHashTrieMapLoadOrStore(b, testDataLarge[:])
}

func benchmarkHashTrieMapLoadOrStore(b *testing.B, data []string) {
	b.ReportAllocs()
	var m cnc.HashTrieMap[string, int]

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
