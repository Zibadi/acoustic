package main

import "log"

func checkArgs(args []string) {
	if len(args) < 2 {
		log.Fatal("specify the path of the directory.")
	}
}
