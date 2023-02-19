package mocks

import (
	"context"

	"github.com/krixlion/dev_forum-lib/event"
	"github.com/stretchr/testify/mock"
)

type Cmd[T any] struct {
	*mock.Mock
}

func NewCmd[T any]() Cmd[T] {
	return Cmd[T]{
		Mock: new(mock.Mock),
	}
}

func (m Cmd[T]) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m Cmd[T]) Consume(ctx context.Context, queue string, eventType event.EventType) (<-chan event.Event, error) {
	args := m.Called(ctx, queue, eventType)
	return args.Get(0).(<-chan event.Event), args.Error(1)
}

func (m Cmd[T]) Create(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Cmd[T]) Update(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Cmd[T]) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
