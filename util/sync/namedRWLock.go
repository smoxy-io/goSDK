package sync

import "sync"

type NamedRWLock struct {
	locks sync.Map
}

func (l *NamedRWLock) Lock(name string) {
	m, _ := l.locks.LoadOrStore(name, &sync.RWMutex{})
	m.(*sync.RWMutex).Lock()
}

func (l *NamedRWLock) Unlock(name string) {
	if m, ok := l.locks.Load(name); ok {
		m.(*sync.RWMutex).Unlock()
	}
}

func (l *NamedRWLock) RLock(name string) {
	m, _ := l.locks.LoadOrStore(name, &sync.RWMutex{})
	m.(*sync.RWMutex).RLock()
}

func (l *NamedRWLock) RUnlock(name string) {
	if m, ok := l.locks.Load(name); ok {
		m.(*sync.RWMutex).RUnlock()
	}
}

func (l *NamedRWLock) Delete(name string) {
	l.locks.Delete(name)
}

func NewNamedRWLock() *NamedRWLock {
	return &NamedRWLock{
		locks: sync.Map{},
	}
}
