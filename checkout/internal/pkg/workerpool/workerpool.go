package workerpool

import (
	"context"
	"route256/checkout/internal/pkg/ratelimit"

	"github.com/aitsvet/debugcharts"
	"google.golang.org/grpc"
)

type workerpool[I, O any] chan struct{}

func New[I, O any](limit int) workerpool[I, O] {
	return make(chan struct{}, limit)
}

type Either[T any] struct {
	Value *T
	Err   error
}

func (wp workerpool[I, O]) Exec(
	ctx context.Context,
	in *I,
	work func(context.Context, *I, ...grpc.CallOption) (*O, error),
	r *ratelimit.Ratelimit,
) <-chan Either[O] {
	result := make(chan Either[O])
	select {
	// если протухщий контекст, то идём в этот кейс без рандома
	case <-ctx.Done():
		// нужна горутина, т.к. result небуферизованный или сделать буфер
		go func() {
			result <- Either[O]{Value: nil, Err: ctx.Err()}
		}()
	default:
		select {
		case <-ctx.Done():
			go func() {
				result <- Either[O]{Value: nil, Err: ctx.Err()}
			}()
		case wp <- struct{}{}:
			go func() {
				r.Ratelimiter <- struct{}{}
				debugcharts.RPS.Add(1)
				val, err := work(ctx, in)
				result <- Either[O]{Value: val, Err: err}
				<-wp
			}()
		}
	}
	return result
}
