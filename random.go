package main

import (
	"math/rand"
)

func shuffle[T any](arr []T, index int) ([]T, int) {
	newIndex := index
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
		if i == newIndex {
			newIndex = j
		} else if j == newIndex {
			newIndex = i
		}
	}
	return arr, newIndex
}
