package worker

import (
	"context"

	"github.com/gammazero/workerpool"
)

type Worker interface {
	Size() int
	Stop()
	StopWait()
	Stopped() bool
	Submit(func())
	SubmitWait(func())
	WaitingQueueSize() int
	Pause(context.Context)
}

func NewWorker(maxWorkers int, fn func(w Worker)) Worker {
	wp := workerpool.New(maxWorkers)

	go func() {
		fn(wp)
	}()

	return wp
}
