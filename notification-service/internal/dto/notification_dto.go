package dto

type NotificationResponse struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	OrderID int    `json:"order_id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
