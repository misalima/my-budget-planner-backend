package dto

import (
	"github.com/google/uuid"
)

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

func GetUUIDValue(ptr *string) uuid.UUID {
	if ptr == nil || *ptr == "" {
		return uuid.Nil
	}
	u, err := uuid.Parse(*ptr)
	if err != nil {
		return uuid.Nil
	}
	return u
}
