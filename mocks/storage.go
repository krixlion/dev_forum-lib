package mocks

import (
	"context"

	"github.com/krixlion/dev_forum-lib/event"
	"github.com/stretchr/testify/mock"
)

type Storage[T any] struct {
	*mock.Mock
}

func NewStorage[T any]() Storage[T] {
	return Storage[T]{
		Mock: new(mock.Mock),
	}
}

func (m Storage[T]) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m Storage[T]) Get(ctx context.Context, id string) (T, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(T), args.Error(1)
}

func (m Storage[T]) GetMultiple(ctx context.Context, offset string, limit string) ([]T, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]T), args.Error(1)
}

func (m Storage[T]) Create(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Storage[T]) Update(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Storage[T]) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m Storage[T]) CatchUp(e event.Event) {
	m.Called(e)
}
