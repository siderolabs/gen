// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package containers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/containers"
)

func TestConcurrentMap(t *testing.T) {
	t.Run("should return nothing if key doesnt exist", func(t *testing.T) {
		m := containers.ConcurrentMap[int, int]{}
		_, ok := m.Get(0)
		require.False(t, ok)
	})

	t.Run("should remove nothing if map is empty", func(t *testing.T) {
		m := containers.ConcurrentMap[int, int]{}
		m.Remove(0)
	})

	t.Run("should return setted value", func(t *testing.T) {
		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)
		val, ok := m.Get(1)
		require.True(t, ok)
		require.Equal(t, 1, val)
	})

	t.Run("should remove value", func(t *testing.T) {
		m := containers.ConcurrentMap[int, int]{}
		m.Set(1, 1)
		m.Remove(1)
		_, ok := m.Get(1)
		require.False(t, ok)
	})

	t.Run("should call fn for every key", func(t *testing.T) {
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
}
