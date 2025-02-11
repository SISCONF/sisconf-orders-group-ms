package sheets

import (
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/files"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
	"github.com/xuri/excelize/v2"
)

type SheetMaker struct {
	foodsList sisconf.FoodsList
	fileName  string
	file      *excelize.File
}

func NewSheetMaker() (*SheetMaker, error) {
	fileName := "data/foods.json"
	foodsList, err := files.ReadStructJSON[sisconf.FoodsList](fileName)
	if err != nil {
		return nil, err
	}

	return &SheetMaker{
		foodsList: *foodsList,
		fileName:  fileName,
		file:      excelize.NewFile(),
	}, nil
}
