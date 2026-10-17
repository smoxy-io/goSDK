package atomic

import (
	"sync"
	"sync/atomic"
)

type Map[K comparable, V any] struct {
	m   *sync.Map
	len *atomic.Int64
}

func (m *Map[K, V]) Get(k K) (v V, ok bool) {
	val, ok := m.m.Load(k)

	if !ok {
		return v, ok
	}

	return val.(V), ok
}

func (m *Map[K, V]) Set(k K, v V) {
	m.m.Store(k, v)
	m.len.Add(1)
}

func (m *Map[K, V]) Delete(k K) {
	m.m.Delete(k)
	m.len.Add(-1)
}

func (m *Map[K, V]) Len() int {
	return int(m.len.Load())
}

func (m *Map[K, V]) Range(f func(K, V) bool) {
	m.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m:   new(sync.Map),
		len: new(atomic.Int64),
	}
}
