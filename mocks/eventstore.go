package mocks

import (
	"context"

	"github.com/krixlion/dev_forum-lib/event"
	"github.com/stretchr/testify/mock"
)

type Eventstore[T any] struct {
	*mock.Mock
}

func NewEventstore[T any]() Eventstore[T] {
	return Eventstore[T]{
		Mock: new(mock.Mock),
	}
}

func (m Eventstore[T]) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m Eventstore[T]) Create(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Eventstore[T]) Update(ctx context.Context, a T) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m Eventstore[T]) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m Eventstore[T]) CatchUp(e event.Event) {
	m.Called(e)
}

func (m Eventstore[T]) Consume(ctx context.Context, queue string, eventType event.EventType) (<-chan event.Event, error) {
	args := m.Called(ctx, queue, eventType)
	return args.Get(0).(<-chan event.Event), args.Error(1)
}
