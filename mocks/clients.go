package mocks

import (
	"context"

	articlePb "github.com/krixlion/dev_forum-proto/article_service/pb"
	userPb "github.com/krixlion/dev_forum-proto/user_service/pb"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var _ userPb.UserServiceClient = (*UserClient)(nil)

type UserClient struct {
	*mock.Mock
}

func NewUserClient() UserClient {
	return UserClient{
		Mock: new(mock.Mock),
	}
}

func (m UserClient) Create(ctx context.Context, in *userPb.CreateUserRequest, opts ...grpc.CallOption) (*userPb.CreateUserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*userPb.CreateUserResponse), args.Error(1)
}

func (m UserClient) Update(ctx context.Context, in *userPb.UpdateUserRequest, opts ...grpc.CallOption) (*userPb.UpdateUserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*userPb.UpdateUserResponse), args.Error(1)
}

func (m UserClient) Delete(ctx context.Context, in *userPb.DeleteUserRequest, opts ...grpc.CallOption) (*userPb.DeleteUserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*userPb.DeleteUserResponse), args.Error(1)
}

func (m UserClient) Get(ctx context.Context, in *userPb.GetUserRequest, opts ...grpc.CallOption) (*userPb.GetUserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*userPb.GetUserResponse), args.Error(1)
}

func (m UserClient) GetSecret(ctx context.Context, in *userPb.GetUserSecretRequest, opts ...grpc.CallOption) (*userPb.GetUserSecretResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*userPb.GetUserSecretResponse), args.Error(1)
}

func (m UserClient) GetStream(ctx context.Context, in *userPb.GetUsersRequest, opts ...grpc.CallOption) (userPb.UserService_GetStreamClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(userPb.UserService_GetStreamClient), args.Error(1)
}

var _ articlePb.ArticleServiceClient = (*ArticleClient)(nil)

type ArticleClient struct {
	*mock.Mock
}

func NewArticleClient() ArticleClient {
	return ArticleClient{
		Mock: new(mock.Mock),
	}
}
func (m ArticleClient) Create(ctx context.Context, in *articlePb.CreateArticleRequest, opts ...grpc.CallOption) (*articlePb.CreateArticleResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*articlePb.CreateArticleResponse), args.Error(1)
}

func (m ArticleClient) Update(ctx context.Context, in *articlePb.UpdateArticleRequest, opts ...grpc.CallOption) (*articlePb.UpdateArticleResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*articlePb.UpdateArticleResponse), args.Error(1)
}

func (m ArticleClient) Delete(ctx context.Context, in *articlePb.DeleteArticleRequest, opts ...grpc.CallOption) (*articlePb.DeleteArticleResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*articlePb.DeleteArticleResponse), args.Error(1)
}

func (m ArticleClient) Get(ctx context.Context, in *articlePb.GetArticleRequest, opts ...grpc.CallOption) (*articlePb.GetArticleResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*articlePb.GetArticleResponse), args.Error(1)
}

func (m ArticleClient) GetStream(ctx context.Context, in *articlePb.GetArticlesRequest, opts ...grpc.CallOption) (articlePb.ArticleService_GetStreamClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(articlePb.ArticleService_GetStreamClient), args.Error(1)
}
