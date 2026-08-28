package scheduler

import (
	"context"
	"sync"
)

type BaseScheduler struct {
	SignalCtx		context.Context
	baseSchedulWg	sync.WaitGroup
}

func NewBaseScheduler(SignalCtx context.Context) (*BaseScheduler) {
	baseScheduler := &BaseScheduler{
		SignalCtx: SignalCtx,
		baseSchedulWg: sync.WaitGroup{},
	}

	return baseScheduler 
}

func (BaseScheduler *BaseScheduler) Close() () {
	BaseScheduler.baseSchedulWg.Wait()
}
