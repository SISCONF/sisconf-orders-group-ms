package sisconf

import "time"

type OrderFood struct {
	FoodName     string
	Quantity     int
	QuantityType string
}

type Order struct {
	CustomerName string
	Details      []OrderFood
}

type OrdersGroup struct {
	Id           int
	TotalPrice   float64
	OrderDate    time.Time
	ItemQuantity int
	Orders       []Order
}
