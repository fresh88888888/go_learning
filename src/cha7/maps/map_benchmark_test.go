package maps

import (
	"strconv"
	"sync"
	"testing"
)

const (
	num_of_reader = 100
	num_of_writer = 10
)

type Map interface {
	Set(key interface{}, val interface{})
	Get(key interface{}) (interface{}, bool)
	Del(key interface{})
}

func benchmarkMap(b *testing.B, hm Map) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < num_of_writer; j++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 100; j++ {
					hm.Set(strconv.Itoa(j), j*j)
					hm.Set(strconv.Itoa(j), j*j)
					hm.Del(strconv.Itoa(j))
				}
				wg.Done()
			}()
		}
		for i := 0; i < num_of_reader; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 100; j++ {
					hm.Get(strconv.Itoa(j))
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	b.StopTimer()
}

func BenchmarkSyncMap(b *testing.B) {
	b.Run("map with RWLock", func(b *testing.B) {
		hm := CreateRWLockMap()
		benchmarkMap(b, hm)
	})

	b.Run("sync.map", func(b *testing.B) {
		hm := CreateSyncMapBenchMarkAdapter()
		benchmarkMap(b, hm)
	})

	b.Run("concurrent map", func(b *testing.B) {
		hm := CreateConCurrentMapBenchMarkAdapter(199)
		benchmarkMap(b, hm)
	})
}
