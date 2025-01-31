package internal

type OrderFood struct {
	foodId   int64
	quantity int
}

type Order struct {
	details []OrderFood
}

type OrdersGroup struct {
	orders []Order
}
