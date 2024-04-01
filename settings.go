package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	dir              string
	imageChar        string
	progressbarChar  string
	isHotColdEnabled bool
}

func newSettings(args []string) Settings {
	checkArgs(args)
	settings := getSettings()
	initSettings(settings)
	return settings
}

func checkArgs(args []string) {
	if len(args) < 2 {
		fmt.Println("[ERROR]: Please provide music directory.")
		os.Exit(0)
	}
}

func getSettings() Settings {
	imageChar := flag.String("image-char", "▄", "Set the character used to display the image.")
	progressbarChar := flag.String("progressbar-char", "─", "Set the character used to display the progress bar.")
	isHotColdEnabled := flag.Bool("hot-cold", false, "Set this option to true to seperate Hot music from Cold ones.")
	flag.Parse()
	return Settings{
		dir:              flag.Arg(0),
		imageChar:        *imageChar,
		progressbarChar:  *progressbarChar,
		isHotColdEnabled: *isHotColdEnabled,
	}
}

func initSettings(s Settings) {
	var err error
	s.dir, err = filepath.Abs(s.dir)
	if err != nil {
		fmt.Printf("[ERROR]: Could not get the absolute path.\n%v\n", err)
		os.Exit(0)
	}
	if s.isHotColdEnabled {
		baseDir := filepath.Dir(s.dir)
		err = os.MkdirAll(filepath.Join(baseDir, "Hot"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(baseDir, "Cold"), os.ModePerm)
		if err != nil {
			fmt.Printf("[ERROR]: Could not create Hot and Cold direcotries.\n%v\n", err)
			os.Exit(0)
		}
	}
}
