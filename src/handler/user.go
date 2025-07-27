package handler

import (
	"context"
	"errors"
	userpb "github.com/Rickykn/drug-proto/gen/user"
	"github.com/Rickykn/user-service/src/logger"
	"github.com/Rickykn/user-service/src/model"
	"github.com/Rickykn/user-service/src/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userSvc service.IUserService
}

func NewUserHandler(svc service.IUserService) *UserHandler {
	return &UserHandler{
		userSvc: svc,
	}
}

func (uh *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	log := logger.WithContext(ctx)

	log.Info().Str("username", req.Username).Msg("get user")

	if req.Username == "" {
		err := errors.New("username is required")
		log.Error().Err(err).Msg("Validation failed")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	user, err := uh.userSvc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("get user failed")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	log.Info().Msg("get user successfully")
	return &userpb.GetUserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil

}

func (uh *UserHandler) Register(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	reqUser := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err := uh.userSvc.CreateUser(ctx, reqUser)

	if err != nil {
		panic("handler error from service")
	}

	return &userpb.CreateUserResponse{
		Code:    200,
		Message: "Success For Register User",
		User: &userpb.User{
			Username:  req.Username,
			Email:     req.Email,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}, nil
}
