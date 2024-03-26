package main

import (
	"math/rand"
)

func shuffle[T any](arr []T) []T {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
