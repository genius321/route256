package ratelimit

import (
	"context"
	"time"
)

type Ratelimit struct {
	Ratelimiter chan struct{}
}

func New(ctx context.Context, limit int) *Ratelimit {
	r := Ratelimit{Ratelimiter: make(chan struct{}, limit)}
	go r.clean(ctx, limit)
	return &r
}

// вычитывает данные из канала с определённым интервалом только по достижению лимита
func (r *Ratelimit) clean(ctx context.Context, limit int) {
	interval := time.Second / time.Duration(limit)
	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if len(r.Ratelimiter) == limit {
				<-r.Ratelimiter
			}
		}
	}
}
