package internal

import (
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
	err = file.SetCellValue(customerName, "A1", "KG")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "E1", "KG")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "B1", "UND")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "F1", "UND")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "C1", "CX")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "G1", "CX")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "D1", "PRODUTO")
	if err != nil {
		return err
	}

	err = file.SetCellValue(customerName, "H1", "PRODUTO")
	return err
}

func CreateOrdersGroupXlsxFile(ordersGroup OrdersGroup) error {
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
