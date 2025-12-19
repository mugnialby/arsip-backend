package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	authRequest "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/auth"
	usersRequest "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/users"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(deleteUserRequest *usersRequest.DeleteUserRequest) error
	FindUserForLoginRequest(loginRequest *authRequest.LoginRequest) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Where("status = ?", "Y").
		Where("role_id not in (?)", 1).
		Preload("Department", "status = ?", "Y").
		Preload("Role", "status = ?", "Y").
		Order("user_id asc").
		Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(deleteUserRequest *usersRequest.DeleteUserRequest) error {
	result := r.db.Model(&model.User{}).
		Where("id = ?", deleteUserRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteUserRequest.SubmittedBy,
			"modified_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no data found to delete")
	}

	return nil
}

func (r *userRepository) FindUserForLoginRequest(loginRequest *authRequest.LoginRequest) (*model.User, error) {
	var user model.User
	err := r.db.Where("status = ?", "Y").
		Where("user_id = ?", loginRequest.UserId).
		Where("password_hash = ?", loginRequest.Password).
		Preload("Department", "status = ?", "Y").
		Preload("Role", "status = ?", "Y").
		First(&user).Error
	return &user, err
}
