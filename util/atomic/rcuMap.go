package atomic

import "maps"

import "sync/atomic"

type RcuMap[K comparable, V any] struct {
	m *atomic.Pointer[map[K]V]
}

func (m *RcuMap[K, V]) Get(key K) (v V, ok bool) {
	p := m.m.Load()

	if p == nil {
		return v, false
	}

	v, ok = (*p)[key]

	return v, ok
}

func (m *RcuMap[K, V]) Set(key K, value V) {
	op := m.m.Load()

	nm := make(map[K]V)

	if op == nil {
		nm[key] = value
		m.m.Store(&nm)

		return
	}

	maps.Copy(nm, *op)

	nm[key] = value
	m.m.Store(&nm)
}

func (m *RcuMap[K, V]) Delete(key K) {
	op := m.m.Load()

	nm := make(map[K]V)

	if op == nil {
		return
	}

	maps.Copy(nm, *op)

	delete(nm, key)

	m.m.Store(&nm)
}

func (m *RcuMap[K, V]) Len() int {
	p := m.m.Load()

	if p == nil {
		return 0
	}

	return len(*p)
}

func (m *RcuMap[K, V]) Range(f func(K, V) bool) {
	p := m.m.Load()

	if p == nil {
		return
	}

	if len(*p) == 0 {
		return
	}

	for k, v := range *p {
		if !f(k, v) {
			break
		}
	}
}

func NewRcuMap[K comparable, V any]() *RcuMap[K, V] {
	m := &RcuMap[K, V]{
		m: new(atomic.Pointer[map[K]V]),
	}

	m.m.Store(new(make(map[K]V)))

	return m
}
