package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"simple-golang/internal/domain/entity"
	"simple-golang/internal/port/outbound"
	"simple-golang/util"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
	CreateUserAccount(ctx context.Context, user entity.UserEntity) error
	UpdatePassword(ctx context.Context, req entity.UserEntity) error

	GetUser(ctx context.Context, query entity.QueryParamEntity) ([]entity.UserEntity, int64, int64, error)
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	UpdateUser(ctx context.Context, user entity.UserEntity) error
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repo       outbound.UserRepositoryInterface
	jwtService JwtServiceInterface
	redis      *redis.Client
}

func NewUserService(repo outbound.UserRepositoryInterface, jwtService JwtServiceInterface, redis *redis.Client) UserServiceInterface {
	return &userService{
		repo:       repo,
		jwtService: jwtService,
		redis:      redis,
	}
}

func (s *userService) SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] SignIn: %v", err)
		return nil, "", err
	}

	if checkPass := util.CheckPasswordHash(req.Password, user.Password); !checkPass {
		err = errors.New("password is incorrect")
		log.Errorf("[UserService-2] SignIn: %v", err)
		return nil, "", err
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		log.Errorf("[UserService-3] SignIn: %v", err)
		return nil, "", err
	}

	sessionData := entity.JwtUserData{
		CreatedAt: time.Now().String(),
		Email:     user.Email,
		LoggedIn:  true,
		Name:      user.Name,
		Token:     token,
		UserID:    user.ID,
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, "", err
	}

	err = s.redis.Set(ctx, token, jsonData, time.Hour*23).Err()
	if err != nil {
		log.Errorf("[UserService-4] SignIn: %v", err)
		return nil, "", err
	}

	return user, token, nil
}

func (s *userService) CreateUserAccount(ctx context.Context, user entity.UserEntity) error {
	passwordNoEncrypt := user.Password
	password, err := util.HashPassword(passwordNoEncrypt)
	if err != nil {
		log.Fatalf("[UserService-1] CreateUser: %v", err)
		return err
	}

	user.Password = password
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) UpdatePassword(ctx context.Context, req entity.UserEntity) error {
	password, err := util.HashPassword(req.Password)
	if err != nil {
		log.Errorf("[UserService-1] UpdatePassword: %v", err)
		return err
	}
	req.Password = password

	err = s.repo.UpdatePassword(ctx, req)
	if err != nil {
		log.Errorf("[UserService-2] UpdatePassword: %v", err)
		return err
	}

	err = s.redis.Del(ctx, req.Token).Err()
	if err != nil {
		log.Errorf("[UserService-3] UpdatePassword: %v", err)
		return err
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, query entity.QueryParamEntity) ([]entity.UserEntity, int64, int64, error) {
	return s.repo.GetUser(ctx, query)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, user entity.UserEntity) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if user, err := s.GetUserByID(ctx, id); err != nil {
		log.Errorf("[UserService-1] DeleteUser: %v", err)
		return err
	} else if user.Email == "superadmin@mail.com" {
		err := errors.New("user not allowed to delete")
		return err
	}

	return s.repo.DeleteUser(ctx, id)
}
