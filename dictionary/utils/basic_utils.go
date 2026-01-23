package utils

import "strconv"

func ParseInt(s string) (int, error) {
	newInt, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return int(newInt), err
	}
	return int(newInt), nil
}
