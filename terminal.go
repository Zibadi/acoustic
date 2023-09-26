package main

import "fmt"

func hideCursor() {
	fmt.Printf("\033[?25l")
}

func showCursor() {
	fmt.Printf("\033[?25h")
}
