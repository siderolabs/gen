// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package containers_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/containers"
	"github.com/siderolabs/gen/xsync"
)

func TestConcurrentMap(t *testing.T) {
	t.Parallel()

	t.Run("should return nothing if key doesnt exist", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		_, ok := m.Get(0)
		require.False(t, ok)
	})

	t.Run("should remove nothing if map is empty", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Remove(0)
	})

	t.Run("should return setted value", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)
		val, ok := m.Get(1)
		require.True(t, ok)
		require.Equal(t, 1, val)
	})

	t.Run("should remove value", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)
		m.Remove(1)
		_, ok := m.Get(1)
		require.False(t, ok)

		m.Set(2, 2)
		got, ok := m.RemoveAndGet(2)
		require.True(t, ok)
		require.Equal(t, 2, got)

		got, ok = m.RemoveAndGet(2)
		require.False(t, ok)
		require.Zero(t, got)

		m.Reset()
		got, ok = m.RemoveAndGet(2)
		require.False(t, ok)
		require.Zero(t, got)

		require.False(t, ok)
	})

	t.Run("should call fn for every key", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		var count int
		m.ForEach(func(key int, value int) {
			count++
		})
		require.Equal(t, 3, count)
	})

	t.Run("should clear the map", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)

		require.Equal(t, 1, m.Len())

		m.Clear()

		require.Equal(t, 0, m.Len())
	})

	t.Run("should trunc the map", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)

		require.Equal(t, 1, m.Len())

		m.Reset()

		require.Equal(t, 0, m.Len())
	})
}

func TestConcurrentMap_GetOrCall(t *testing.T) {
	var m containers.ConcurrentMap[int, int]

	t.Run("group", func(t *testing.T) {
		t.Run("try to insert value", func(t *testing.T) {
			parallelGetOrCall(t, &m, 100, 1000)
		})

		t.Run("try to insert value #2", func(t *testing.T) {
			parallelGetOrCall(t, &m, 1000, 100)
		})
	})
}

func parallelGetOrCall(t *testing.T, m *containers.ConcurrentMap[int, int], our, another int) {
	t.Parallel()

	oneAnotherGet := false

	for i := 0; i < 10000; i++ {
		key := int(rand.Int63n(10000))

		res, ok := m.GetOrCall(key, func() int { return key * our })
		if ok {
			switch res {
			case key * our:
			case key * another:
				oneAnotherGet = true
			default:
				t.Fatalf("unexpected value %d", res)
			}
		}
	}

	require.True(t, oneAnotherGet)
}

func TestConcurrentMap_GetOrCreate(t *testing.T) {
	var m containers.ConcurrentMap[int, int]

	t.Run("group", func(t *testing.T) {
		t.Run("try to insert value", func(t *testing.T) {
			parallelGetOrCreate(t, &m, 100, 1000)
		})

		t.Run("try to insert value #2", func(t *testing.T) {
			parallelGetOrCreate(t, &m, 1000, 100)
		})
	})
}

func parallelGetOrCreate(t *testing.T, m *containers.ConcurrentMap[int, int], our, another int) {
	t.Parallel()

	oneAnotherGet := false

	for i := 0; i < 10000; i++ {
		key := int(rand.Int63n(10000))

		res, ok := m.GetOrCreate(key, key*our)
		if ok {
			switch res {
			case key * our:
			case key * another:
				oneAnotherGet = true
			default:
				t.Fatalf("unexpected value %d", res)
			}
		}
	}

	require.True(t, oneAnotherGet)
}

func Example_benchConcurrentMap() {
	var sink int

	benchResult := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()

		var m containers.ConcurrentMap[int, *xsync.Once[int]]

		for i := 0; i < b.N; i++ {
			variable := 0

			res, _ := m.GetOrCall(10, func() *xsync.Once[int] {
				return &xsync.Once[int]{}
			})

			sink = res.Do(func() int {
				variable++

				return variable
			})
		}
	})

	if benchResult.AllocsPerOp() > 0 {
		fmt.Println("this benchmark should not allocate memory")
	}

	fmt.Println("ok")

	// Output:
	// ok

	_ = sink
}
