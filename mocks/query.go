package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Query[T any] struct {
	*mock.Mock
}

func (m Query[T]) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m Query[T]) Get(ctx context.Context, id string) (T, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(T), args.Error(1)
}

func (m Query[T]) GetMultiple(ctx context.Context, offset string, limit string) ([]T, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]T), args.Error(1)
}

func (m Query[T]) GetBelongingIDs(ctx context.Context, userId string) ([]string, error) {
	panic("not implemented") // TODO: Implement
}
func (m Query[T]) Create(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Query[T]) Update(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Query[T]) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
