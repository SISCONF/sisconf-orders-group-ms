package sisconf

import "math"

type FoodsList []Food

func (foodsList FoodsList) GetSheetColsGroupCount(sheetFoodRowsPerFoodGroupCount float64) float64 {
	return math.Ceil(
		float64(len(foodsList)) /
			sheetFoodRowsPerFoodGroupCount)
}

func (foodsList FoodsList) GetFoodsNames() []string {
	var names []string = make([]string, len(foodsList))
	for index, food := range foodsList {
		names[index] = food.Name
	}
	return names
}
