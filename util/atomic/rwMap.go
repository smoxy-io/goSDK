package atomic

import "sync"

type RWMap[K comparable, V any] struct {
	mu *sync.RWMutex
	m  map[K]V
}

func (m *RWMap[K, V]) Get(k K) (v V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok = m.m[k]
	return v, ok
}

func (m *RWMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[k] = v
}

func (m *RWMap[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.m, k)
}

func (m *RWMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.m)
}

func (m *RWMap[K, V]) Range(f func(K, V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}

func NewRWMap[K comparable, V any]() *RWMap[K, V] {
	return &RWMap[K, V]{
		mu: new(sync.RWMutex),
		m:  make(map[K]V),
	}
}
