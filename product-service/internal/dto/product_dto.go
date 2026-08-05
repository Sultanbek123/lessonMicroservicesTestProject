package dto

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductResponse struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
