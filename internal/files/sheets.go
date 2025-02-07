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

const headerFillAndRowsColor string = "00B050"

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

	var foodsList []sisconf.Food
	err = json.Unmarshal(foodsByes, &foodsList)
	if err != nil {
		return err
	}

	if len(foodsList) == 0 {
		return errors.New("foods list is empty")
	}

	foodsListGroupsCount := math.Ceil(float64(len(foodsList)) / float64(56))

	var currentCellLetter rune = 'A'
	var currentGroupIndex int = 1
	groupsHeaderValue := map[int]string{
		1: "KG",
		2: "UND",
		3: "CX",
		4: "PRODUTO",
	}
	for currentCellsGroup := 0.0; currentCellsGroup < foodsListGroupsCount*4.0; currentCellsGroup++ {
		cell := fmt.Sprintf("%c1", currentCellLetter)
		err = file.SetCellValue(customerName, cell, groupsHeaderValue[currentGroupIndex])
		if err != nil {
			return err
		}

		if currentGroupIndex == 4 {
			currentGroupIndex = 1
		} else {
			currentGroupIndex++
		}
		currentCellLetter += 1
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
