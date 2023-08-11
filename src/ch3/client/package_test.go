package client

import (
	"ch3/service"
	"testing"

	cmap "github.com/easierway/concurrent_map"
)

func TestService(t *testing.T) {
	t.Log(service.GetFibonacci(8))
	t.Log(service.Square(23))
}

func TestCurrentMapInit(t *testing.T) {
	m := cmap.CreateConcurrentMap(10)
	m.Set(cmap.StrKey("key"), 10)
	t.Log(m.Get(cmap.StrKey("key")))
}
