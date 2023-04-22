package auth

import (
	"context"

	"google.golang.org/grpc"
)

type Header string

const AuthHeader Header = "authorization"

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		panic("not implemented")
	}
}
