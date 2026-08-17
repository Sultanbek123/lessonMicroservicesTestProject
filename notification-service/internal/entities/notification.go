package entities

type Notification struct {
	ID      int
	UserID  string
	OrderID int
	Message string
	Status  string
}
