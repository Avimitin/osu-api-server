package server

import (
	"fmt"
	"strconv"
)

func atoi(s string) (int, error) {
	i, e := strconv.Atoi(s)
	if e != nil {
		return -1, fmt.Errorf("cast %s to int: %v", s, e)
	}
	return i, nil
}

func parseInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("cast %s to int64: %v", s, err)
	}
	return i, nil
}

func parseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1.0, fmt.Errorf("cast %s to float32: %v", s, err)
	}
	return f, nil
}
