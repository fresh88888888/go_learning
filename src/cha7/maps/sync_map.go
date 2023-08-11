package maps

import "sync"

type SyncMapBenchMarkAdapter struct {
	m sync.Map
}

func CreateSyncMapBenchMarkAdapter() *SyncMapBenchMarkAdapter {
	return &SyncMapBenchMarkAdapter{}
}

func (s *SyncMapBenchMarkAdapter) Set(key interface{}, val interface{}) {
	s.m.Store(key, val)
}

func (s *SyncMapBenchMarkAdapter) Get(key interface{}) (interface{}, bool) {
	return s.m.Load(key)
}

func (s *SyncMapBenchMarkAdapter) Del(key interface{}) {
	s.m.Delete(key)
}
