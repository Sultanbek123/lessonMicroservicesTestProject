package service

import (
	"errors"

	"product-service/internal/entities"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(jwtSecret []byte) *AuthService {
	return &AuthService{jwtSecret: jwtSecret}
}

func (s *AuthService) ValidateToken(tokenString string) (*entities.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok != true {
			return nil, errors.New("неправильный метод подпись")
		}
		return s.jwtSecret, nil
	})
	if err != nil || token.Valid == false {
		return nil, errors.New("токен недействительный or expired")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok != true {
		return nil, errors.New("can't read claims")
	}
	userID, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)

	var tokenClaims entities.TokenClaims
	tokenClaims.UserID = userID
	tokenClaims.Role = role
	return &tokenClaims, nil
}
