package repository

type Order struct {
	ID		int64
	Name	string
}

type OrderDatabase interface {
	ListOrders() ([]*Order, error)
	AddOrder(order *Order) (int64, error)
}
