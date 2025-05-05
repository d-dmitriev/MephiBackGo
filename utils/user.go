package utils

import (
	"strconv"
)

// ParseUserID преобразует строковое значение userID в uint
// Используется во всех сервисах для приведения userID из контекста к нужному типу
func ParseUserID(userID string) (uint, error) {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
