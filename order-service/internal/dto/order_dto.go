package dto

type CreateOrderRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderResponse struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

type ProductResponse struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
