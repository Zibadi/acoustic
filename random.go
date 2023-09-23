package main

import (
	"math/rand"
	"time"
)

func randInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max)
	return randNum
}

func shuffle(arr []string) []string {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
