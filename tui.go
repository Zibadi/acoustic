package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/mattn/go-runewidth"
	"github.com/nfnt/resize"
)

func run(p *Player, s *Settings) {
	go listenToKeyboard(p)
	for len(p.musics) > 0 {
		err := p.play(s)
		if err != nil {
			p.skipMusic()
		}
	}
}

func printMetadata(p *Player, s *Settings) {
	file, _ := os.Open(p.getCurrentMusic().path)
	defer file.Close()
	var err error
	metadata, err := readMusicMetadata(file)
	if err != nil {
		fmt.Printf("[WARNING]: Could not load the music meta tag of %v\n%v\n", p.getCurrentMusic().path, err)
	} else {
		printMusicImage(metadata, s.imageChar)
		printMusicMetadata(metadata)
	}
}

func printMusicImage(m tag.Metadata, c string) {
	defer checkImage()
	data := m.Picture().Data
	reader := bytes.NewReader(data)
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Println("[WARNING]: Could not decode the music image.")
		return
	}
	printImage(image, c)
}

func checkImage() {
	if r := recover(); r != nil {
		printCenter("[NO IMAGE]")
	}
}

func printImage(img image.Image, char string) {
	width, height, _ := getTerminalSize()
	min := math.Min(float64(width), float64(height))
	size := uint(min)
	image := resize.Resize(size, 0, img, resize.Lanczos3)
	maxY := image.Bounds().Max.Y - 1
	maxX := image.Bounds().Max.X

	for y := 0; y < maxY; y += 2 {
		for i := 0; i < (width-maxX)/2; i++ {
			fmt.Printf(" ")
		}
		for x := 0; x < maxX; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			if char == "▄" {
				fmt.Printf("\033[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)
			} else {
				fmt.Printf("\033[48;2;%d;%d;%d", r>>8, g>>8, b>>8)
			}

			r, g, b, _ = image.At(x, y+1).RGBA()
			fmt.Printf("\033[38;2;%d;%d;%dm%v", r>>8, g>>8, b>>8, char)
		}
		fmt.Printf("\033[0m")
		fmt.Printf("\n")
	}
}

func printMusicMetadata(m tag.Metadata) {
	printCenter(m.Title())
	printCenter(m.Artist())
	printCenter(strconv.Itoa(m.Year()))
	printCenter(m.Genre())
}

func printDuration(p *Player) {
	minutes := int(p.duration.Seconds()) / 60
	seconds := int(p.duration.Seconds()) % 60
	duratoin := fmt.Sprintf("[%d:%02d]", minutes, seconds)
	printCenter(duratoin)
}

func printProgressbar(p *Player, s *Settings) {
	if !p.isPaused {
		fmt.Print(s.progressbarChar)
	}
}

func printCenter(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	width, _, _ := getTerminalSize()
	length := runewidth.StringWidth(text)
	for i := 0; i < (width-length)/2; i++ {
		fmt.Print(" ")
	}
	fmt.Println(text)
}

func getTerminalSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("[ERROR]: Could not get the terminal size.\n%v\n", err)
		return 0, 0, err
	}
	trim := strings.TrimRight(string(output), "\n")
	size := strings.Split(trim, " ")
	height, err := strconv.Atoi(size[0])
	width, _ := strconv.Atoi(size[1])
	if err != nil {
		fmt.Printf("[ERROR]: Could not parse the terminal size.\n%v\n", err)
		return 0, 0, err
	}
	return width, height, nil
}
