// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package containers

import "sync"

// ConcurrentMap is a map that can be safely accessed from multiple goroutines.
type ConcurrentMap[K comparable, V any] struct {
	m  map[K]V
	mx sync.Mutex
}

// Get returns the value for the given key.
func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	val, ok := m.m[key]

	return val, ok
}

// GetOrCreate returns the existing value for the key if present. Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *ConcurrentMap[K, V]) GetOrCreate(key K, val V) (V, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if res, ok := m.m[key]; ok {
		return res, true
	}

	if m.m == nil {
		m.m = map[K]V{}
	}

	m.m[key] = val

	return val, false
}

// GetOrCall returns the existing value for the key if present. Otherwise, it calls fn, stores the result and returns it.
// The loaded result is true if the value was loaded, false if it was created using fn.
//
// The main reason for this function is to avoid unnecessary allocations if you use pointer types as values, since
// compiler cannot prove that the value does not escape if it's not stored.
func (m *ConcurrentMap[K, V]) GetOrCall(key K, fn func() V) (V, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if res, ok := m.m[key]; ok {
		return res, true
	}

	if m.m == nil {
		m.m = map[K]V{}
	}

	val := fn()

	m.m[key] = val

	return val, false
}

// Set sets the value for the given key.
func (m *ConcurrentMap[K, V]) Set(key K, val V) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.m == nil {
		m.m = map[K]V{}
	}

	m.m[key] = val
}

// Remove removes the value for the given key.
func (m *ConcurrentMap[K, V]) Remove(key K) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.m == nil {
		return
	}

	delete(m.m, key)
}

// RemoveAndGet removes the value for the given key and returns it if it exists.
func (m *ConcurrentMap[K, V]) RemoveAndGet(key K) (V, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.m == nil {
		return *new(V), false
	}

	val, ok := m.m[key]
	delete(m.m, key)

	return val, ok
}

// ForEach calls the given function for each key-value pair.
func (m *ConcurrentMap[K, V]) ForEach(f func(K, V)) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for k, v := range m.m {
		f(k, v)
	}
}

// Len returns the number of elements in the map.
func (m *ConcurrentMap[K, V]) Len() int {
	m.mx.Lock()
	defer m.mx.Unlock()

	return len(m.m)
}

// Clear removes all key-value pairs.
func (m *ConcurrentMap[K, V]) Clear() {
	m.mx.Lock()
	defer m.mx.Unlock()

	for k := range m.m {
		delete(m.m, k)
	}
}

// Reset resets the underlying map.
func (m *ConcurrentMap[K, V]) Reset() {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.m = nil
}
