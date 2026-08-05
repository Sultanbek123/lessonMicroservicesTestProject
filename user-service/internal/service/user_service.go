package service

import (
	"context"
	"user-service/internal/dto"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, err
	}

	id, err := s.repo.CreateUser(ctx, req.Username, string(hashedPassword), req.Role)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       id,
		Username: req.Username,
		Role:     req.Role,
	}, nil
}
