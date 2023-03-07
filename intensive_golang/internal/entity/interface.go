package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetTotal() (float64, error)
}
