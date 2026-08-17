package event

type OrderCreated struct {
	OrderID   int    `json:"order_id"`
	UserID    string `json:"user_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
