package dto

func GetIntValue(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func GetFloatValue(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}
