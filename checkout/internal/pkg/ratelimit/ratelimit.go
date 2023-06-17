package ratelimit

import (
	"context"
	"time"
)

type ratelimit chan struct{}

func New(ctx context.Context, limit int) chan<- struct{} {
	r := make(chan struct{}, limit)
	go ratelimit(r).clean(ctx, limit)
	return r
}

// вычитывает данные из канала с определённым интервалом
func (r ratelimit) clean(ctx context.Context, limit int) {
	interval := time.Second / time.Duration(limit)
	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			select {
			// если канал закрыт, выходим из clean
			case <-r:
				return
			default:
			}
			if len(r) == limit {
				<-r
			}
		}
	}
}
