package mocks

import (
	"context"

	"github.com/krixlion/dev_forum-lib/event"
	"github.com/stretchr/testify/mock"
)

type CQRStorage[T any] struct {
	*mock.Mock
}

func NewCQRStorage[T any]() CQRStorage[T] {
	return CQRStorage[T]{
		Mock: new(mock.Mock),
	}
}

func (m CQRStorage[T]) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m CQRStorage[T]) Get(ctx context.Context, id string) (T, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(T), args.Error(1)
}

func (m CQRStorage[T]) GetMultiple(ctx context.Context, offset string, limit string) ([]T, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]T), args.Error(1)
}

func (m CQRStorage[T]) EventHandlers() map[event.EventType][]event.Handler {
	panic("not implemented") // TODO: Implement
}

func (m CQRStorage[T]) GetBelongingIDs(ctx context.Context, userId string) ([]string, error) {
	panic("not implemented") // TODO: Implement
}

func (m CQRStorage[T]) Create(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m CQRStorage[T]) Update(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m CQRStorage[T]) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m CQRStorage[T]) CatchUp(e event.Event) {
	m.Called(e)
}
