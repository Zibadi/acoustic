package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	dir               string
	imageChar         string
	progressbarChar   string
	isCoolColdEnabled bool
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
	imageChar := flag.String("imageChar", "▄", "Set the character used to display the image.")
	progressbarChar := flag.String("progressbarChar", "─", "Set the character used to display the progress bar.")
	isCoolColdEnabled := flag.Bool("coolCold", false, "Set this option to true to seperate cool musics from cold ones.")
	flag.Parse()
	return Settings{
		dir:               flag.Arg(0),
		imageChar:         *imageChar,
		progressbarChar:   *progressbarChar,
		isCoolColdEnabled: *isCoolColdEnabled,
	}
}

func initSettings(s Settings) {
	var err error
	s.dir, err = filepath.Abs(s.dir)
	if err != nil {
		fmt.Printf("[ERROR]: Could not get the absolute path.\n%v\n", err)
		os.Exit(0)
	}
	if s.isCoolColdEnabled {
		baseDir := filepath.Dir(s.dir)
		err = os.MkdirAll(filepath.Join(baseDir, "COOL"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(baseDir, "cold"), os.ModePerm)
		if err != nil {
			fmt.Printf("[ERROR]: Could not create COOL and cold direcotries.\n%v\n", err)
			os.Exit(0)
		}
	}
}
