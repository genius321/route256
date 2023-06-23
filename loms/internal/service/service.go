package service

import (
	"context"
	orderModels "route256/loms/internal/models/order"
)

// type TransactionManager interface {
// 	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
// 	Serializable(ctx context.Context, fn func(ctxTx context.Context) error) error
// }

// type Repository interface {
// 	CreateOrder(
// 		ctx context.Context,
// 		user orderModels.User,
// 		items orderModels.Items,
// 	) (orderModels.OrderId, error)
// }

// type Transport interface {
// 	CreateOrder(
// 		ctx context.Context,
// 		user orderModels.User,
// 		items orderModels.Items,
// 	) (orderModels.OrderId, error)
// }

type Service struct {
	// Transport
	// Repository
	// TransactionManager
}

// func NewService(r Repository, tm TransactionManager) *Service {
// 	return &Service{Repository: r, TransactionManager: tm}
// }

// func NewService(t Transport) *Service {
// 	return &Service{Transport: t}
// }

// есть просто какой-то бизнес, который отвечает заглушкой
// он ничего не знает про транспортный уровень
func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateOrder(
	ctx context.Context,
	user orderModels.User,
	items orderModels.Items,
) (orderModels.OrderId, error) {
	return 666, nil
}
