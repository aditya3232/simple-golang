package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"simple-golang/internal/adapter/outbound/postgres/model"
	"simple-golang/internal/domain/entity"
	outboundport "simple-golang/internal/port/outbound"

	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) outboundport.UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	modelUser := model.User{}
	if err := r.db.WithContext(ctx).Where("id =?", id).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] DeleteCustomer: User not found")
			return err
		}
		log.Errorf("[UserRepository-2] DeleteCustomer: %v", err)
		return err
	}

	if err := r.db.Delete(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-3] DeleteCustomer: %v", err)
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.UserEntity) error {
	var (
		modelUser = model.User{}
		updates   = map[string]interface{}{}
	)

	if err := r.db.WithContext(ctx).Where("id =?", user.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] UpdateUser: User not found")
			return err
		}
		log.Errorf("[UserRepository-2] UpdateUser: %v", err)
		return err
	}

	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Phone != "" {
		updates["phone"] = user.Phone
	}
	if user.Address != "" {
		updates["address"] = user.Address
	}

	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&modelUser).Updates(updates).Error; err != nil {
			log.Errorf("[UserRepository-3] UpdateUser: %v", err)
			return err
		}
	}

	return nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.UserEntity) error {
	modeluser := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
		Address:  user.Address,
	}

	if err := r.db.WithContext(ctx).Create(&modeluser).Error; err != nil {
		log.Errorf("[UserRepository-1] CreateUser: %v", err)
		return err
	}

	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := r.db.WithContext(ctx).Where("id =?", id).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] GetUserByID: User not found")
			return nil, err
		}
		log.Errorf("[UserRepository-2] GetUserByID: %v", err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Name:     modelUser.Name,
		Email:    modelUser.Email,
		Password: modelUser.Password,
		Phone:    modelUser.Phone,
		Address:  modelUser.Address,
	}, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := r.db.WithContext(ctx).Where("email =?", email).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] GetUserByEmail: User not found")
			return nil, err
		}
		log.Errorf("[UserRepository-2] GetUserByEmail: %v", err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Name:     modelUser.Name,
		Email:    modelUser.Email,
		Password: modelUser.Password,
		Phone:    modelUser.Phone,
		Address:  modelUser.Address,
	}, nil
}

func (r *userRepository) GetUser(ctx context.Context, query entity.QueryParamEntity) ([]entity.UserEntity, int64, int64, error) {
	var (
		modelUsers   []model.User
		respEntities []entity.UserEntity
		countData    int64
	)

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := r.db.WithContext(ctx).
		Select("id, name, email, password, phone, address", "created_at", "updated_at").
		Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			fmt.Sprintf("%%%s%%", query.Search),
			fmt.Sprintf("%%%s%%", query.Search),
			fmt.Sprintf("%%%s%%", query.Search))

	if err := sqlMain.Model(&modelUsers).Count(&countData).Error; err != nil {
		log.Errorf("[UserRepository-1] GetUser: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.Order(order).Limit(int(query.Limit)).Offset(int(offset)).Find(&modelUsers).Error; err != nil {
		log.Errorf("[UserRepository-2] GetUser: %v", err)
		return nil, 0, 0, err
	}

	if len(modelUsers) < 1 {
		err := errors.New("404")
		log.Infof("[UserRepository-3] GetUser: No users found")
		return nil, 0, 0, err
	}

	for _, v := range modelUsers {
		respEntities = append(respEntities, entity.UserEntity{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			Password: v.Password,
			Phone:    v.Phone,
			Address:  v.Address,
		})
	}

	return respEntities, countData, int64(totalPage), nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, req entity.UserEntity) error {
	var modelUser model.User
	if err := r.db.WithContext(ctx).Where("id = ?", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[UserRepository-1] UpdatePassword: User not found")
			return errors.New("404")
		}
		log.Errorf("[UserRepository-2] UpdatePassword: %v", err)
		return err
	}

	// Update hanya kolom password
	if err := r.db.WithContext(ctx).
		Model(&modelUser).
		Where("id = ?", req.ID).
		Updates(map[string]any{
			"password": req.Password,
		}).Error; err != nil {
		log.Errorf("[UserRepository-3] UpdatePassword: %v", err)
		return err
	}

	return nil
}
