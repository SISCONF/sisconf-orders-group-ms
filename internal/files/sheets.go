package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
	"github.com/xuri/excelize/v2"
)

type FoodsList []sisconf.Food

const sheetColsPerFoodsGroupCount rune = 4
const sheetFoodRowsPerFoodGroupCount float64 = 56
const headerFillAndRowsColor string = "00B050"

func (foodsList FoodsList) GetSheetColsGroupCount() float64 {
	return math.Ceil(
		float64(len(foodsList)) /
			float64(sheetFoodRowsPerFoodGroupCount))
}

func createHeaderRowStyle(file *excelize.File) (int, error) {
	const headerRowFontColor string = "000000" // Black

	return file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Color: []string{headerFillAndRowsColor}, Type: "pattern", Pattern: 1},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Calibri",
			Color:  headerRowFontColor,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
}

func createColsStyle(file *excelize.File) (int, error) {
	return file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Family: "Calibri",
			Color:  headerFillAndRowsColor,
		},
	})
}

func writeOrdersGroupXlsxHeader(file *excelize.File, customerName string) error {
	var err error
	foodsByes, err := os.ReadFile("foods.json")
	if err != nil {
		return err
	}

	var foodsList FoodsList
	err = json.Unmarshal(foodsByes, &foodsList)
	if err != nil {
		return err
	}

	if len(foodsList) == 0 {
		return errors.New("foods list is empty")
	}

	foodsListGroupsCount := foodsList.GetSheetColsGroupCount()

	var currentCellLetter rune = 'A'
	for currentCellsGroup := 0.0; currentCellsGroup < foodsListGroupsCount; currentCellsGroup++ {
		cell := fmt.Sprintf("%c1", currentCellLetter)
		err = file.SetSheetRow(customerName, cell, &[]any{"KG", "UND", "CX", "PRODUTO"})
		if err != nil {
			return err
		}

		currentCellLetter += sheetColsPerFoodsGroupCount
	}

	err = file.SetColWidth(customerName, "A", fmt.Sprintf("%c", currentCellLetter), 18)

	return err
}

func CreateOrdersGroupXlsxFile(ordersGroup sisconf.OrdersGroup) error {
	var err error
	file := excelize.NewFile()
	defer func() error {
		return file.Close()
	}()

	for _, order := range ordersGroup.Orders {
		_, err := file.NewSheet(order.CustomerName)
		if err != nil {
			return err
		}

		headerStyleIndex, err := createHeaderRowStyle(file)
		if err != nil {
			return err
		}

		colsStyleIndex, err := createColsStyle(file)
		if err != nil {
			return err
		}

		err = file.SetColStyle(order.CustomerName, "D", colsStyleIndex)
		if err != nil {
			return err
		}

		err = file.SetColStyle(order.CustomerName, "H", colsStyleIndex)
		if err != nil {
			return err
		}

		err = file.SetRowStyle(order.CustomerName, 1, 1, headerStyleIndex)
		if err != nil {
			return err
		}

		err = writeOrdersGroupXlsxHeader(file, order.CustomerName)
		if err != nil {
			return err
		}
	}

	err = file.SaveAs("Teste.xlsx")

	return err
}
