package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secret, issuer string) *JWTService {
	return &JWTService{
		secretKey: secret,
		issuer:    issuer,
	}
}

type Claim struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token
func (j *JWTService) GenerateToken(userID uint) (string, error) {
	claims := &Claim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken verifies token validity and returns claims
func (j *JWTService) ValidateToken(tokenStr string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// usage
// jwtSvc := utils.NewJWTService(cfg.JWTSecret, "go-gin-app")
// token, _ := jwtSvc.GenerateToken(user.ID)
// claims, _ := jwtSvc.ValidateToken(token)
