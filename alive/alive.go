package alive

import (
	"context"
	"time"
)

func KeepAlive() (context.Context, chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan struct{}, 10)
	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				cancel()
				return
			case <-sig:
			}
		}
	}()

	return ctx, sig
}
