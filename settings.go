package main

import (
	"flag"
	"fmt"
	"os"
)

type Settings struct {
	dir             string
	imageChar       string
	progressbarChar string
}

func newSettings(args []string) *Settings {
	checkArgs(args)
	return getSettings()
}

func checkArgs(args []string) {
	if len(args) < 2 {
		fmt.Println("[ERROR]: Please provide song directory.")
		os.Exit(0)
	}
}

func getSettings() *Settings {
	imageChar := flag.String("imageChar", "▄", "Set the character used to display the image. (defualt: ▄)")
	progressbarChar := flag.String("progressbarChar", "-", "Set the character used to display the progress bar. (defualt: -)")
	flag.Parse()
	return &Settings{
		dir:             flag.Arg(0),
		imageChar:       *imageChar,
		progressbarChar: *progressbarChar,
	}
}
