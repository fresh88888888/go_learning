package maps

import (
	"github.com/easierway/concurrent_map"
)

type ConCurrentMapBenchMarkAdapter struct {
	cm *concurrent_map.ConcurrentMap
}

func CreateConCurrentMapBenchMarkAdapter(numOfPartitions int) *ConCurrentMapBenchMarkAdapter {
	conMap := concurrent_map.CreateConcurrentMap(numOfPartitions)
	return &ConCurrentMapBenchMarkAdapter{conMap}
}

func (c *ConCurrentMapBenchMarkAdapter) Set(key interface{}, val interface{}) {
	c.cm.Set(concurrent_map.StrKey(key.(string)), val)
}

func (c *ConCurrentMapBenchMarkAdapter) Get(key interface{}) (interface{}, bool) {
	return c.cm.Get(concurrent_map.StrKey(key.(string)))
}

func (c *ConCurrentMapBenchMarkAdapter) Del(key interface{}) {
	c.cm.Del(concurrent_map.StrKey(key.(string)))
}
