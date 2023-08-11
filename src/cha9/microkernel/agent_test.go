package microkernel

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type DemoCollector struct {
	eventReceiver EventReceiver
	agentContext  context.Context
	stopChan      chan struct{}
	name          string
	content       string
}

func newCollector(name string, content string) *DemoCollector {
	return &DemoCollector{
		stopChan: make(chan struct{}),
		name:     name,
		content:  content,
	}
}

func (dc *DemoCollector) Init(eventReceiver EventReceiver) error {
	fmt.Println("initialize collector...", dc.name)
	dc.eventReceiver = eventReceiver
	return nil
}

func (dc *DemoCollector) Start(ac context.Context) error {
	fmt.Println("start collector...", dc.name)
	for {
		select {
		case <-ac.Done():
			dc.stopChan <- struct{}{}
			break
		default:
			time.Sleep(time.Microsecond * 50)
			dc.eventReceiver.OnEvent(Event{dc.name, dc.content})
		}
	}
}

func (dc *DemoCollector) Stop() error {
	fmt.Println("stop collector...", dc.name)
	select {
	case <-dc.stopChan:
		return nil
	case <-time.After(time.Second * 1):
		return errors.New("Failed to stop for timeout...")
	}
}

func (dc *DemoCollector) Destory() error {
	fmt.Println(dc.name, " release resources...")
	return nil
}

func TestAgent(t *testing.T) {
	agt := newAgent(100)
	c1 := newCollector("c1", "1")
	c2 := newCollector("c2", "2")
	agt.registerCollector("c1", c1)
	agt.registerCollector("c2", c2)
	agt.Start()
	fmt.Println(agt.Start())
	time.Sleep(time.Second * 1)
	agt.Stop()
	agt.Destory()
}
