package service

import (
	"context"
	"errors"
	"time"

	"user-service/internal/entities"
	"user-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(u *repository.UserRepository, jwtSecret []byte) *AuthService {
	return &AuthService{userRepo: u, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("password not match")
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
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
