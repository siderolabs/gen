// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package containers_test

import (
	"fmt"
	"testing"

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
}
