package common

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type Heartbeat struct {
	cancelc chan struct{}
}

func StartHeartbeat(ctx context.Context, coolingTime time.Duration) (h *Heartbeat) {
	h = &Heartbeat{cancelc: make(chan struct{})}
	go func() {
		for {
			select {
			case <-h.cancelc:
				return
			default:
				{
					time.Sleep(coolingTime)
					activity.RecordHeartbeat(ctx)
				}
			}
		}
	}()
	return
}

func (h *Heartbeat) Stop() {
	close(h.cancelc)
}
