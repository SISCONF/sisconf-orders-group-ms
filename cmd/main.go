package main

import (
	"time"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/files"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
)

func main() {
	ordersGroup := sisconf.OrdersGroup{
		TotalPrice:   64.5,
		OrderDate:    time.Now(),
		ItemQuantity: 2,
		Orders: []sisconf.Order{
			{
				CustomerName: "Alyson",
				Details: []sisconf.OrderFood{
					{
						FoodName:     "Teste",
						Quantity:     12,
						QuantityType: "KG",
					},
				},
			},
		},
	}

	files.CreateOrdersGroupXlsxFile(ordersGroup)
}
