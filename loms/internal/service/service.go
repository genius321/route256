package service

import (
	"context"
	orderModels "route256/loms/internal/models/order"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Repository interface {
	CreateOrder(
		ctx context.Context,
		user orderModels.User,
		items orderModels.Items,
	) (orderModels.OrderId, error)
}

// ничего не знает про транспортный уровень,
// но знает, что бд и тракнзакционный мендеджер реализует необходимое поведение
type Service struct {
	Repository
	TransactionManager
}

func NewService(r Repository, tm TransactionManager) *Service {
	return &Service{Repository: r, TransactionManager: tm}
}

func (s *Service) CreateOrder(
	ctx context.Context,
	user orderModels.User,
	items orderModels.Items,
) (orderModels.OrderId, error) {
	return s.Repository.CreateOrder(ctx, user, items)
}
