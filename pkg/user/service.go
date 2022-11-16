package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

type UserService interface {
	CreateUser(ctx context.Context, userName string) (*User, error)
	GetUser(ctx context.Context, uid uint64) (*User, error)
	GetOrCreateGoogleUser(ctx context.Context, email, userName string) (*User, error)
}

type UserServiceImpl struct {
	userRepo UserRepo
	sf       common.IDGenerator
}

func NewUserService(userRepo UserRepo, sf common.IDGenerator) UserService {
	return &UserServiceImpl{userRepo, sf}
}
func (svc *UserServiceImpl) CreateUser(ctx context.Context, userName string) (*User, error) {
	userID, err := svc.sf.NextID()
	if err != nil {
		return nil, fmt.Errorf("error create snowflake ID: %w", err)
	}
	user, err := svc.userRepo.CreateUser(ctx, &User{
		ID:   userID,
		Name: userName,
	})
	if err != nil {
		return nil, fmt.Errorf("error create user %d: %w", userID, err)
	}
	return user, nil
}
func (svc *UserServiceImpl) GetUser(ctx context.Context, uid uint64) (*User, error) {
	user, err := svc.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error get user %d: %w", uid, err)
	}
	return user, nil
}

func (svc *UserServiceImpl) GetOrCreateGoogleUser(ctx context.Context, email, userName string) (*User, error) {
	var user *User
	var err error
	user, err = svc.userRepo.GetGoogleUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("error get google user %s: %w", email, err)
		}
		userID, err := svc.sf.NextID()
		if err != nil {
			return nil, fmt.Errorf("error create snowflake ID: %w", err)
		}
		user, err = svc.userRepo.CreateGoogleUser(ctx, email, &User{
			ID:   userID,
			Name: userName,
		})
		if err != nil {
			return nil, fmt.Errorf("error create google user %s: %w", email, err)
		}
		return user, nil
	}
	return user, nil
}
