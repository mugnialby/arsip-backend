package service

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	authRequest "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/auth"
	usersRequest "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/users"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(deleteUserRequest *usersRequest.DeleteUserRequest) error {
	return s.repo.Delete(deleteUserRequest)
}

func (s *UserService) CheckUserLoginRequest(loginRequest *authRequest.LoginRequest) (*model.User, error) {
	return s.repo.FindUserForLoginRequest(loginRequest)
}
