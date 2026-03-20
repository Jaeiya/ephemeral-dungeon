package utils

import (
	"math/rand/v2"
	"slices"
	"strconv"
)

// Shuffle uses the fisher-yates established algorithm
func Shuffle[T any](items []T) {
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}

// FilterSlice returns a new slice filtered by the filterFunc.
//
// 🔵 filterFunc() filters OUT an item when it returns true
func FilterSlice[T any](items []T, filterFunc func(arg T) bool) []T {
	items = slices.Clone(items)
	return slices.DeleteFunc(items, filterFunc)
}

func ParseInt(s string) (int, error) {
	newInt, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return int(newInt), err
	}
	return int(newInt), nil
}

func HasLowercaseOnly(s string) bool {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}
