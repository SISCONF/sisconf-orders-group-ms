package main

import (
	"fmt"
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
						FoodName:     "ABACATE",
						Quantity:     12,
						QuantityType: "KG",
					},
					{
						FoodName:     "ABACAXI GRANDE",
						Quantity:     10,
						QuantityType: "CX",
					},
				},
			},
		},
	}

	err := files.CreateOrdersGroupXlsxFile(ordersGroup)
	if err != nil {
		fmt.Println(err.Error())
	}
}
