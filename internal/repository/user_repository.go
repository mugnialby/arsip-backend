package repository

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/auth"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	FindUserForLoginRequest(loginRequest *request.LoginRequest) (*model.User, error)
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

func (r *userRepository) FindUserForLoginRequest(loginRequest *request.LoginRequest) (*model.User, error) {
	var user model.User
	err := r.db.Where("status = ?", "Y").
		Where("user_id = ?", loginRequest.UserId).
		Where("password_hash = ?", loginRequest.Password).
		First(&user).Error
	return &user, err
}
