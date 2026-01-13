package utils

import (
	"math/rand/v2"
	"slices"
)

// Shuffle uses the fisher-yates established algorithm but returns
// a copy of the items instead of the original slice.
func Shuffle[T any](items []T) {
	for i := len(items) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

// FilterSlice returns a new slice filtered by the filterFunc.
//
// 🔵 filterFunc() filters OUT an item when it returns true
func FilterSlice[T any](items []T, filterFunc func(arg T) bool) []T {
	items = slices.Clone(items)
	return slices.DeleteFunc(items, filterFunc)
}
