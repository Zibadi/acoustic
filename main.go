package main

import (
	"os"
)

func main() {
	settings := newSettings(os.Args)
	player := newPlayer(&settings)
	run(&player)
}
