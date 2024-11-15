package mocks

import (
	"github.com/krixlion/dev_forum-lib/event"
	"github.com/stretchr/testify/mock"
)

type Handler struct {
	*mock.Mock
}

func NewHandler() Handler {
	return Handler{
		Mock: new(mock.Mock),
	}
}

func (h Handler) Handle(e event.Event) {
	h.Called(e)
}
