package main

import (
	"fmt"
	"time"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal"
)

func main() {
	ordersGroup := internal.OrdersGroup{
		TotalPrice:   64.50,
		OrderDate:    time.Now(),
		ItemQuantity: 10,
		Orders: []internal.Order{
			{
				CustomerName: "Alyson",
				Details: []internal.OrderFood{
					{
						FoodName:     "Banana",
						Quantity:     10,
						QuantityType: "KG",
					},
				},
			},
		},
	}
	err := internal.CreateOrdersGroupXlsxFile(ordersGroup)
	if err != nil {
		fmt.Println("Erro ao criar documento!")
		return
	}
	fmt.Println("Documento criado com sucesso!")
}
