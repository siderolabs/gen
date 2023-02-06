// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:dupl
package containers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/containers"
)

func TestLazyBiMap(t *testing.T) {
	prev := 0

	m := containers.LazyBiMap[int, int]{
		Creator: func(i int) (int, error) {
			if i == prev {
				return 0, fmt.Errorf("unxpected call to Creator")
			}

			if i < 0 {
				return 0, fmt.Errorf("key must be >= 0")
			}

			prev = i

			return (i % 10) * 100, nil
		},
	}

	t.Run("should return nothing if key or value doesnt exist", func(t *testing.T) {
		_, ok := m.Get(0)
		require.False(t, ok)

		_, ok = m.GetInverse(0)
		require.False(t, ok)
	})

	t.Run("should create value on demand", func(t *testing.T) {
		create, err := m.GetOrCreate(1)
		require.NoError(t, err)
		require.Equal(t, 100, create)
	})

	t.Run("should return existing value", func(t *testing.T) {
		create, err := m.GetOrCreate(1)
		require.NoError(t, err)
		require.Equal(t, 100, create)
	})

	t.Run("should return existing key by value", func(t *testing.T) {
		key, ok := m.GetInverse(100)
		require.True(t, ok)

		assert.Equal(t, 1, key)
	})

	t.Run("should remove old key if has new key", func(t *testing.T) {
		create, err := m.GetOrCreate(11)
		require.NoError(t, err)
		require.Equal(t, 100, create)
		val, ok := m.Get(1)
		require.False(t, ok)
		require.Equal(t, 0, val)
	})

	t.Run("should return error on Creator error", func(t *testing.T) {
		_, err := m.GetOrCreate(-1)
		require.Error(t, err)
	})

	t.Run("should remove key", func(t *testing.T) {
		_, err := m.GetOrCreate(12)
		require.NoError(t, err)

		m.Remove(12)
		_, ok := m.Get(12)
		require.False(t, ok)
	})

	t.Run("should remove key by value", func(t *testing.T) {
		_, err := m.GetOrCreate(13)
		require.NoError(t, err)

		m.RemoveInverse(300)
		_, ok := m.Get(13)
		require.False(t, ok)
	})

	t.Run("should remove all entries", func(t *testing.T) {
		_, err := m.GetOrCreate(14)
		require.NoError(t, err)

		_, err = m.GetOrCreate(15)
		require.NoError(t, err)

		m.Clear()

		_, ok := m.Get(13)
		require.False(t, ok)

		_, ok = m.Get(14)
		require.False(t, ok)
	})

	t.Run("should iterate over entries", func(t *testing.T) {
		var keys []int
		var values []int

		m.ForEach(func(k int, v int) {
			keys = append(keys, k)
			values = append(values, v)
		})

		assert.Empty(t, keys)
		assert.Empty(t, values)

		_, err := m.GetOrCreate(1)
		require.NoError(t, err)

		_, err = m.GetOrCreate(2)
		require.NoError(t, err)

		m.ForEach(func(k int, v int) {
			keys = append(keys, k)
			values = append(values, v)
		})

		assert.Equal(t, []int{1, 2}, keys)
		assert.Equal(t, []int{100, 200}, values)
	})
}

func TestLazyMap(t *testing.T) {
	prev := 0

	m := containers.LazyMap[int, int]{
		Creator: func(i int) (int, error) {
			if i == prev {
				return 0, fmt.Errorf("unxpected call to Creator")
			}

			if i < 0 {
				return 0, fmt.Errorf("key must be >= 0")
			}

			prev = i

			return i * 100, nil
		},
	}

	t.Run("should delete nothing if there is nothing to delete", func(t *testing.T) {
		m.Remove(0)
	})

	t.Run("should return nothing if key doesnt exist", func(t *testing.T) {
		_, ok := m.Get(0)
		require.False(t, ok)
	})

	t.Run("should create value on demand", func(t *testing.T) {
		create, err := m.GetOrCreate(1)
		require.NoError(t, err)
		require.Equal(t, 100, create)
	})

	t.Run("should return existing value", func(t *testing.T) {
		create, err := m.GetOrCreate(1)
		require.NoError(t, err)
		require.Equal(t, 100, create)
	})

	t.Run("should return error on Creator error", func(t *testing.T) {
		_, err := m.GetOrCreate(-1)
		require.Error(t, err)
	})

	t.Run("should properly remove key", func(t *testing.T) {
		m.Remove(1)
		_, ok := m.Get(1)
		require.False(t, ok)
	})

	t.Run("should remove all entries", func(t *testing.T) {
		_, err := m.GetOrCreate(2)
		require.NoError(t, err)

		_, err = m.GetOrCreate(3)
		require.NoError(t, err)

		m.Clear()

		_, ok := m.Get(1)
		require.False(t, ok)

		_, ok = m.Get(2)
		require.False(t, ok)
	})

	t.Run("should iterate over entries", func(t *testing.T) {
		var keys []int
		var values []int

		m.ForEach(func(k int, v int) {
			keys = append(keys, k)
			values = append(values, v)
		})

		assert.Empty(t, keys)
		assert.Empty(t, values)

		_, err := m.GetOrCreate(4)
		require.NoError(t, err)

		_, err = m.GetOrCreate(5)
		require.NoError(t, err)

		m.ForEach(func(k int, v int) {
			keys = append(keys, k)
			values = append(values, v)
		})

		assert.Equal(t, []int{4, 5}, keys)
		assert.Equal(t, []int{400, 500}, values)
	})
}
