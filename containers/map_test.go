// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package containers_test

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/containers"
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

		m.ForEach(func(int, int) { count++ })

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

	t.Run("filter map", func(t *testing.T) {
		t.Parallel()

		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		m.FilterInPlace(func(key int, val int) bool {
			return key == 1 || val == 3
		})

		require.Equal(t, 2, m.Len())
	})
}

func TestConcurrentMap_GetOrCall(t *testing.T) {
	var m containers.ConcurrentMap[int, int]

	res, ok := m.GetOrCall(1, func() int { return 1 })
	require.False(t, ok)
	require.Equal(t, 1, res)

	res, ok = m.GetOrCall(1, func() int { return 2 })
	require.True(t, ok)
	require.Equal(t, 1, res)

	m.Clear()

	generatedKeys := generateUniqueKeyVals(10000)

	t.Run("group", func(t *testing.T) {
		t.Run("try to insert value", func(t *testing.T) {
			parallelGetOrCall(t, &m, generatedKeys)
		})

		t.Run("try to insert value #2", func(t *testing.T) {
			parallelGetOrCall(t, &m, generatedKeys)
		})
	})
}

type keyVal struct {
	key int
	val int
}

func generateUniqueKeyVals(num int) []keyVal {
	res := make([]keyVal, 0, num)

	for i := range num {
		res = append(res, keyVal{
			key: i,
			val: rand.IntN(num),
		})
	}

	return res
}

func parallelGetOrCall(t *testing.T, m *containers.ConcurrentMap[int, int], keyVals []keyVal) {
	t.Parallel()

	cloned := slices.Clone(keyVals)

	rand.Shuffle(len(cloned), func(i, j int) { cloned[i], cloned[j] = cloned[j], cloned[i] })

	for _, keyVal := range cloned {
		res, ok := m.GetOrCall(keyVal.key, func() int { return keyVal.val })
		if ok && res != keyVal.val {
			// tests above ensure that ok can only be true if key is already in the map
			t.Fatalf("unexpected value %d", res)
		}
	}
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

	for range 10000 {
		key := int(rand.Int64N(10000))

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

		var m containers.ConcurrentMap[int, func() int]

		for b.Loop() {
			variable := 0

			res, _ := m.GetOrCall(10, func() func() int {
				return sync.OnceValue(func() int {
					variable++

					return variable
				})
			})

			sink = res()
		}
	})

	if allocsPerOp := benchResult.AllocsPerOp(); allocsPerOp > 1 {
		fmt.Printf("this benchmark should not make more than one allocation, but it makes %d allocations per operation", allocsPerOp)
	}

	fmt.Println("ok")
	fmt.Println(sink)

	// Output:
	// ok
	// 1
}
