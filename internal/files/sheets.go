package files

import (
	"errors"
	"fmt"
	"log"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

const colsPerFoodsGroupCount rune = 4
const foodRowsPerFoodGroupCount float64 = 56

var colors = map[string]string{
	"black": "#000000",
	"green": "#00B050",
}

var border = []excelize.Border{
	{
		Style: 1,
		Type:  "left",
		Color: colors["black"],
	},
	{
		Style: 1,
		Type:  "right",
		Color: colors["black"],
	},
	{
		Style: 1,
		Type:  "top",
		Color: colors["black"],
	},
	{
		Style: 1,
		Type:  "bottom",
		Color: colors["black"],
	},
}

func createHeaderRowStyle(file *excelize.File) (int, error) {
	return file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Color: []string{colors["green"]}, Type: "pattern", Pattern: 1},
		Font: &excelize.Font{
			Bold:  true,
			Color: colors["black"],
			Size:  10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: border,
	})
}

func createCommonColStyle(file *excelize.File) (int, error) {
	return file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: border,
	})
}

func createProductsColStyle(file *excelize.File) (int, error) {
	return file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  10,
			Color: colors["green"],
		},
		Border: border,
	})
}

func writeOrdersGroupXlsxHeader(file *excelize.File, customerName string, foodsList *sisconf.FoodsList, commonColsStyleIndex, productColStyleIndex, headerRowStyleIndex int) error {
	var err error
	foodsListGroupsCount := foodsList.GetSheetColsGroupCount(foodRowsPerFoodGroupCount)
	var currentCellLetter rune = 'A'

	for currentCellsGroup := 0.0; currentCellsGroup < foodsListGroupsCount; currentCellsGroup++ {
		cell := fmt.Sprintf("%c1", currentCellLetter)
		err = file.SetSheetRow(customerName, cell, &[]any{"KG", "UND", "CX", "PRODUTO"})
		if err != nil {
			return err
		}

		productsCellLetter := currentCellLetter + colsPerFoodsGroupCount - 1
		foodsGroupColRange := fmt.Sprintf(
			"%c:%c",
			currentCellLetter,
			productsCellLetter,
		)
		err = file.SetColStyle(customerName, foodsGroupColRange, commonColsStyleIndex)
		if err != nil {
			return err
		}

		err = file.SetColStyle(customerName, string(productsCellLetter), productColStyleIndex)
		if err != nil {
			return err
		}

		currentCellLetter += colsPerFoodsGroupCount
		err = file.SetColWidth(
			customerName,
			string(productsCellLetter),
			string(productsCellLetter),
			20,
		)
		if err != nil {
			return err
		}
	}

	err = file.SetRowStyle(customerName, 1, 1, headerRowStyleIndex)
	return err
}

func writeFoodsListToGroupXlsx(file *excelize.File, customerName string, foodsList *sisconf.FoodsList) error {
	var err error
	var startSliceIndex int = 0
	var endSliceIndex int
	var foodsNames []string = foodsList.GetFoodsNames()
	productsCells, err := file.SearchSheet(customerName, "PRODUTO")

	for productCellIndex, productCell := range productsCells {
		endSliceIndex = (int(productCellIndex) + 1) * int(foodRowsPerFoodGroupCount)
		if endSliceIndex > len(*foodsList) {
			endSliceIndex = len(*foodsList)
		}

		productCellStartingDataRow := fmt.Sprintf("%c2", rune(productCell[0]))
		foodsNamesSlice := foodsNames[startSliceIndex:endSliceIndex]
		err = file.SetSheetCol(customerName, productCellStartingDataRow, &foodsNamesSlice)
		if err != nil {
			return err
		}

		startSliceIndex += int(foodRowsPerFoodGroupCount)
	}

	return err
}

func writeFoodsQuantityToXlsx(file *excelize.File, order *sisconf.Order) error {
	var err error
	quantityTypeSheetMap := map[string]rune{
		"KG":  3,
		"UND": 2,
		"CX":  1,
	}
	for _, orderDetail := range order.Details {
		cells, err := file.SearchSheet(order.CustomerName, orderDetail.FoodName)
		if err != nil || len(cells) == 0 {
			log.Printf("couldn't find food %s\n", orderDetail.FoodName)
		} else {
			productCell := cells[0]
			quantityCellLetter := rune(productCell[0]) - quantityTypeSheetMap[orderDetail.QuantityType]
			quantityCellNumber := rune(productCell[1])
			quantityCell := fmt.Sprintf("%c%c", quantityCellLetter, quantityCellNumber)
			file.SetCellValue(order.CustomerName, quantityCell, orderDetail.Quantity)
		}
	}
	return err
}

func CreateOrdersGroupXlsxFile(ordersGroup sisconf.OrdersGroup) error {
	var err error
	file := excelize.NewFile()
	defer func() error {
		return file.Close()
	}()

	err = file.SetDefaultFont("Calibri")
	if err != nil {
		return err
	}

	headerRowStyleIndex, err := createHeaderRowStyle(file)
	if err != nil {
		return err
	}

	productsColsStyleIndex, err := createProductsColStyle(file)
	if err != nil {
		return err
	}

	commonColsStyleIndex, err := createCommonColStyle(file)
	if err != nil {
		return err
	}

	foodsList, err := ReadStructJSON[sisconf.FoodsList]("./internal/data/foods.json")
	if err != nil {
		return err
	}

	if len(*foodsList) == 0 {
		return errors.New("foods list is empty")
	}

	for _, order := range ordersGroup.Orders {
		sheetIndex, err := file.NewSheet(order.CustomerName)
		if err != nil {
			return err
		}
		file.SetActiveSheet(sheetIndex)

		err = writeOrdersGroupXlsxHeader(
			file,
			order.CustomerName,
			foodsList,
			commonColsStyleIndex,
			productsColsStyleIndex,
			headerRowStyleIndex,
		)
		if err != nil {
			return err
		}

		err = writeFoodsListToGroupXlsx(file, order.CustomerName, foodsList)
		if err != nil {
			return err
		}

		err = writeFoodsQuantityToXlsx(file, &order)
		if err != nil {
			return err
		}
	}

	// Default sheet created, not useful
	err = file.DeleteSheet("Sheet1")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("pedido_geral_%s.xlsx", uuid.NewString())
	err = file.SaveAs(filename)

	return err
}
