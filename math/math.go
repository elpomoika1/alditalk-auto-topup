package math

import "strconv"

func CalculateRemaining(strAllocated, strUsed string) (float64, error) {
	allocated, err := strconv.ParseInt(strAllocated, 10, 64)
	if err != nil {
		return 0.0, err
	}

	used, err := strconv.ParseInt(strUsed, 10, 64)
	if err != nil {
		return 0.0, err
	}

	remaining := float64(allocated-used) / 1024 / 1024

	return remaining, nil
}
