package lock

import (
	"fmt"
	"sync"
	"testing"
)

var cache map[string]string

const num_of_redaer int = 40
const read_times int = 1000

func init() {
	cache = make(map[string]string)

	cache["a"] = "aa"
	cache["b"] = "bb"
}

func lockFreeAccess() {
	var wg sync.WaitGroup
	wg.Add(num_of_redaer)
	for i := 0; i < num_of_redaer; i++ {
		go func() {
			for j := 0; j < read_times; j++ {
				_, err := cache["a"]
				if !err {
					fmt.Println("nothing...")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func lockAccess() {
	var wg sync.WaitGroup
	wg.Add(num_of_redaer)
	m := new(sync.RWMutex)
	for i := 0; i < num_of_redaer; i++ {
		go func() {
			for j := 0; j < read_times; j++ {
				m.RLock()
				_, err := cache["a"]
				if !err {
					fmt.Println("nothing...")
				}
				m.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkLockFreeAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lockFreeAccess()
	}
	b.StartTimer()
}

func BenchmarkLockAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lockAccess()
	}
	b.StartTimer()
}
