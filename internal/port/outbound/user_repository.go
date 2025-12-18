package outbound

import (
	"context"
	"simple-golang/internal/domain/entity"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user entity.UserEntity) error
	GetUser(ctx context.Context, query entity.QueryParamEntity) ([]entity.UserEntity, int64, int64, error)
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	UpdateUser(ctx context.Context, user entity.UserEntity) error
	DeleteUser(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, req entity.UserEntity) error
}
