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

// ForEach calls the given function for each key-value pair.
func (m *ConcurrentMap[K, V]) ForEach(f func(K, V)) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for k, v := range m.m {
		f(k, v)
	}
}

// Clear removes all key-value pairs.
func (m *ConcurrentMap[K, V]) Clear() {
	m.mx.Lock()
	defer m.mx.Unlock()

	for k := range m.m {
		delete(m.m, k)
	}
}
