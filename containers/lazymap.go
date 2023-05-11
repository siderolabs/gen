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
		return *new(V), err
	}

	m.biMap.Set(key, val)

	return val, nil
}

// Get returns the value for the given key.
func (m *LazyBiMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.biMap.Get(key)

	return val, ok
}

// GetInverse returns the key for the given value.
func (m *LazyBiMap[K, V]) GetInverse(value V) (K, bool) {
	key, ok := m.biMap.GetInverse(value)

	return key, ok
}

// Remove removes the value for the given key.
func (m *LazyBiMap[K, V]) Remove(key K) {
	m.biMap.Remove(key)
}

// RemoveInverse removes the key for the given value.
func (m *LazyBiMap[K, V]) RemoveInverse(value V) {
	m.biMap.RemoveInverse(value)
}

// Clear removes all values.
func (m *LazyBiMap[K, V]) Clear() {
	m.biMap.Clear()
}

// ForEach calls the given function for each key-value pair.
func (m *LazyBiMap[K, V]) ForEach(f func(K, V)) {
	m.biMap.ForEach(f)
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

// GetInverse returns the key for the given value.
func (m *BiMap[K, V]) GetInverse(val V) (K, bool) {
	key, ok := m.v2k[val]

	return key, ok
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

// Remove deletes the value for the given key.
func (m *BiMap[K, V]) Remove(key K) {
	if m.k2v == nil {
		return
	}

	val, ok := m.k2v[key]
	if !ok {
		return
	}

	delete(m.k2v, key)
	delete(m.v2k, val)
}

// RemoveInverse deletes the key for the given value.
func (m *BiMap[K, V]) RemoveInverse(val V) {
	if m.v2k == nil {
		return
	}

	key, ok := m.v2k[val]
	if !ok {
		return
	}

	delete(m.v2k, val)
	delete(m.k2v, key)
}

// Clear removes all key-value pairs.
func (m *BiMap[K, V]) Clear() {
	m.k2v = nil
	m.v2k = nil
}

// ForEach calls the given function for each key-value pair.
func (m *BiMap[K, V]) ForEach(f func(K, V)) {
	for k, v := range m.k2v {
		f(k, v)
	}
}

// LazyMap is like usual map but creates values on demand.
type LazyMap[K comparable, V comparable] struct {
	Creator func(K) (V, error)
	dataMap map[K]V
}

// GetOrCreate returns the value for the given key. It creates it using Creator if it doesn't exist.
func (m *LazyMap[K, V]) GetOrCreate(key K) (V, error) {
	if m.dataMap == nil {
		m.dataMap = map[K]V{}
	}

	val, ok := m.dataMap[key]
	if ok {
		return val, nil
	}

	val, err := m.Creator(key)
	if err != nil {
		return *new(V), err
	}

	m.dataMap[key] = val

	return val, nil
}

// Get returns the value for the given key.
func (m *LazyMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.dataMap[key]

	return val, ok
}

// Remove deletes the value for the given key.
func (m *LazyMap[K, V]) Remove(key K) {
	if m.dataMap == nil {
		return
	}

	delete(m.dataMap, key)
}

// Clear removes all key-value pairs.
func (m *LazyMap[K, V]) Clear() {
	m.dataMap = nil
}

// ForEach calls the given function for each key-value pair.
func (m *LazyMap[K, V]) ForEach(f func(K, V)) {
	for k, v := range m.dataMap {
		f(k, v)
	}
}
