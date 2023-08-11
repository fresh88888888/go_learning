package microkernel

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var WrongStateError = errors.New("can not take the operation in the current state.")

type CollectorsError struct {
	CollectorErrors []error
}

type Event struct {
	name    string
	content string
}

type EventReceiver interface {
	OnEvent(evt Event)
}

type Collector interface {
	Init(evtReceiver EventReceiver) error
	Start(ctx context.Context) error
	Stop() error
	Destory() error
}

type Agent struct {
	collectors map[string]Collector
	evtBuf     chan Event
	cancel     context.CancelFunc
	ctx        context.Context
	state      int
}

const (
	Waiting int = 0
	Running int = 1
)

func (ce CollectorsError) Error() string {
	var errStr strings.Builder
	for _, err := range ce.CollectorErrors {
		errStr.WriteString(err.Error())
	}

	return errStr.String()
}

func (agent *Agent) OnEvent(evt Event) {
	agent.evtBuf <- evt
}

func (agent *Agent) registerCollector(name string, collector Collector) error {
	if agent.state != Waiting {
		return WrongStateError
	}

	agent.collectors[name] = collector
	return collector.Init(agent)
}

func (agent *Agent) startCollector() error {
	var err error
	var errs CollectorsError
	var mutex sync.Mutex
	for name, collector := range agent.collectors {
		go func(name string, collector Collector, ctx context.Context) {
			defer func() {
				mutex.Unlock()
			}()
			err = collector.Start(ctx)
			mutex.Lock()
			if err != nil {
				errs.CollectorErrors = append(errs.CollectorErrors, errors.New(name+":"+err.Error()))
			}
		}(name, collector, agent.ctx)
	}

	return errs
}

func (agent *Agent) stopCollector() error {
	var err error
	var errs CollectorsError
	for name, collector := range agent.collectors {
		if err = collector.Stop(); err != nil {
			errs.CollectorErrors = append(errs.CollectorErrors, errors.New(name+":"+err.Error()))
		}
	}

	return errs
}

func (agent *Agent) destoryCollector() error {
	var err error
	var errs CollectorsError
	for name, collector := range agent.collectors {
		if err = collector.Destory(); err != nil {
			errs.CollectorErrors = append(errs.CollectorErrors, errors.New(name+":"+err.Error()))
		}
	}

	return errs
}

func (agent *Agent) eventprocessGroutine() {
	var evtSeg [10]Event
	for {
		for i := 0; i < 10; i++ {
			select {
			case evtSeg[i] = <-agent.evtBuf:
			case <-agent.ctx.Done():
				return
			}
		}
		fmt.Println(evtSeg)
	}
}

func newAgent(sizeBuf int) *Agent {
	agt := Agent{
		collectors: map[string]Collector{},
		evtBuf:     make(chan Event, sizeBuf),
		state:      Waiting,
	}

	return &agt
}

func (agent *Agent) Start() error {
	if agent.state != Waiting {
		return WrongStateError
	}
	agent.state = Running
	agent.ctx, agent.cancel = context.WithCancel(context.Background())
	go agent.eventprocessGroutine()
	return agent.startCollector()
}

func (agent *Agent) Stop() error {
	if agent.state != Running {
		return WrongStateError
	}
	agent.state = Waiting
	agent.cancel()

	return agent.stopCollector()
}

func (agent *Agent) Destory() error {
	if agent.state != Waiting {
		return WrongStateError
	}

	return agent.destoryCollector()
}
