package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type UnaryHandler struct {
	*mock.Mock
}

func NewUnaryHandler() UnaryHandler {
	return UnaryHandler{Mock: new(mock.Mock)}
}

func (m UnaryHandler) GetMock() grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		args := m.Called()
		return args.Get(0), args.Error(1)
	}
}
