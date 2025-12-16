package service

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// func (s *AuthService) GetAuthByID(id uint) (*model.Auth, error) {
// 	return s.repo.FindByID(id)
// }
