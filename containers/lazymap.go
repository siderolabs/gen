// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package containers provides generic containers.
package containers

// LazyBiMap is like BiMap but creates values on demand.
type LazyBiMap[K comparable, V comparable] struct {
	Creator func(K) (V, error)
	biMap   BiMap[K, V]
}

// GetOrCreate returns the value for the given key.
func (m *LazyBiMap[K, V]) GetOrCreate(key K) (V, error) {
	val, ok := m.biMap.Get(key)
	if ok {
		return val, nil
	}

	val, err := m.Creator(key)
	if err != nil {
		return *new(V), err //nolint:gocritic
	}

	m.biMap.Set(key, val)

	return val, nil
}

// Get returns the value for the given key.
func (m *LazyBiMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.biMap.Get(key)

	return val, ok
}

// BiMap (or “bidirectional map”) is a special kind of map that maintains
// an inverse view of the map while ensuring that no duplicate values are present
// and a value can always be used safely to get the key back.
type BiMap[K comparable, V comparable] struct {
	k2v map[K]V
	v2k map[V]K
}

// Get returns the value for the given key.
func (m *BiMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.k2v[key]

	return val, ok
}

// Set sets the value for the given key.
func (m *BiMap[K, V]) Set(key K, val V) {
	if m.k2v == nil {
		m.k2v = map[K]V{}
		m.v2k = map[V]K{}
	}

	// Ensure that there is only one key per value
	if prevKey, ok := m.v2k[val]; ok && prevKey != key {
		delete(m.k2v, prevKey)
	}

	m.k2v[key] = val
	m.v2k[val] = key
}

// LazyMap is like usual map but creates values on demand.
type LazyMap[K comparable, V comparable] struct {
	Creator func(K) (V, error)
	biMap   map[K]V
}

// GetOrCreate returns the value for the given key. It creates it using Creator if it doesn't exist.
func (m *LazyMap[K, V]) GetOrCreate(key K) (V, error) {
	if m.biMap == nil {
		m.biMap = map[K]V{}
	}

	val, ok := m.biMap[key]
	if ok {
		return val, nil
	}

	val, err := m.Creator(key)
	if err != nil {
		return *new(V), err //nolint:gocritic
	}

	m.biMap[key] = val

	return val, nil
}

// Get returns the value for the given key.
func (m *LazyMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.biMap[key]

	return val, ok
}

// Remove deletes the value for the given key.
func (m *LazyMap[K, V]) Remove(key K) {
	if m.biMap == nil {
		return
	}

	delete(m.biMap, key)
}
