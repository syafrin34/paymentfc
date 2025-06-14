package grpc

import (
	"context"
	"log"
	"paymentfc/proto/userpb"
	"time"

	"google.golang.org/grpc"
)

type UserClient interface {
	GetUserInfoByUserID(ctx context.Context, userID int64) (*userpb.GetUserInfoResult, error)
}
type userClient struct {
	Client userpb.UserServiceClient
}

func NewUserClient() UserClient {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed init grpc user client: %v", err)
	}
	client := userpb.NewUserServiceClient(conn)
	return &userClient{
		Client: client,
	}
}

func (uc userClient) GetUserInfoByUserID(ctx context.Context, userID int64) (*userpb.GetUserInfoResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	userInfo, err := uc.Client.GetUserInfoByUserID(ctx, &userpb.GetUserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
