package workerpool

import (
	"context"

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

func (wp workerpool[I, O]) Exec(ctx context.Context, in *I, work func(context.Context, *I, ...grpc.CallOption) (*O, error)) <-chan Either[O] {
	result := make(chan Either[O])
	select {
	case <-ctx.Done():
		result <- Either[O]{Err: ctx.Err()}
	case wp <- struct{}{}:
		go func() {
			val, err := work(ctx, in)
			result <- Either[O]{Value: val, Err: err}
			close(result)
			<-wp
		}()
	}
	return result
}
