package entities

type Order struct {
	ID        int
	UserID    string
	ProductID int
	Quantity  int
	Status    string
}
