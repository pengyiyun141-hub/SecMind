package scheduler

import (
	"context"
	"sync"
)

type BaseScheduler struct {
	ctx				context.Context
	cancel 			context.CancelFunc
	baseSchedulWg	sync.WaitGroup
}

func NewBaseScheduler() (*BaseScheduler) {
	ctx, cancel := context.WithCancel(context.Background())
	
	baseScheduler := &BaseScheduler{
		ctx: ctx,
		cancel: cancel,
		baseSchedulWg: sync.WaitGroup{},
	}

	return baseScheduler 
}

func (BaseScheduler *BaseScheduler)Stop() () {

}
