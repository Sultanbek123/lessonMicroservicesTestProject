package service

import (
	"context"
	"product-service/internal/dto"
	"product-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	// Add business logic/validation here if needed
	
	id, err := s.repo.CreateProduct(ctx, req)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ProductResponse{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
	}, nil
}

func (s *ProductService) GetProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	return s.repo.GetProducts(ctx)
}
