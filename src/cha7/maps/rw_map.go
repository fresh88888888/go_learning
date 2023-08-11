package maps

import "sync"

type RWLockMap struct {
	m    map[interface{}]interface{}
	lock sync.RWMutex
}

func (rw *RWLockMap) Get(key interface{}) (interface{}, bool) {
	rw.lock.RLock()
	v, ok := rw.m[key]
	rw.lock.RUnlock()
	return v, ok
}

func (rw *RWLockMap) Set(key interface{}, val interface{}) {
	rw.lock.Lock()
	rw.m[key] = val
	rw.lock.Unlock()
}

func (rw *RWLockMap) Del(key interface{}) {
	rw.lock.Lock()
	delete(rw.m, key)
	rw.lock.Unlock()
}

func CreateRWLockMap() *RWLockMap {
	m := make(map[interface{}]interface{}, 0)
	return &RWLockMap{m: m}
}
