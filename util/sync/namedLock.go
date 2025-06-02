package sync

import "sync"

type NamedLock struct {
	locks sync.Map
}

func (l *NamedLock) Lock(name string) {
	m, _ := l.locks.LoadOrStore(name, &sync.Mutex{})
	m.(*sync.Mutex).Lock()
}

func (l *NamedLock) Unlock(name string) {
	if m, ok := l.locks.Load(name); ok {
		m.(*sync.Mutex).Unlock()
	}
}

func (l *NamedLock) Delete(name string) {
	l.locks.Delete(name)
}

func NewNamedLock() *NamedLock {
	return &NamedLock{
		locks: sync.Map{},
	}
}
