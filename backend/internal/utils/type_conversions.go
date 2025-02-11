package utils

import (
	"fmt"
	"strconv"
)

func ConvertToType[T ~int64](s string) (T, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid input: %s", s)
	}
	return T(value), nil
}
