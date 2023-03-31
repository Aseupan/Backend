package utils

var FoodPoints = map[string]int{
	"fruit and vegetables": 40,
	"starchy food":         50,
	"fast food":            50,
	"dairy":                50,
	"protein":              100,
	"fat":                  150,
	"snack":                30,
	"hydration":            40,
}

func GetFoodPoints(foodType string) int {
	price, ok := FoodPoints[foodType]
	if !ok {
		return 0
	}
	return price
}
