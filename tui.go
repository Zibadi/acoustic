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
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/ssh/terminal"
)

func run(p *Player, s *Settings) {
	key := make(chan rune)
	go readKey(key)
	for len(p.songs) > 0 {
		err := p.play(s, key)
		if err != nil {
			p.skipSong()
		}
	}
}

func printMetadata(p *Player, s *Settings) {
	file, _ := os.Open(p.getCurrentSong().path)
	defer file.Close()
	var err error
	p.metadata, err = readSongMetadata(file)
	if err != nil {
		fmt.Printf("[WARNING]: Could not load the song meta tag of %v\n%v\n", p.getCurrentSong().path, err)
	} else {
		printSongImage(p, s.imageChar)
		printSongDetails(p)
	}
}

func printSongImage(p *Player, char string) {
	defer checkImage()
	data := p.metadata.Picture().Data
	reader := bytes.NewReader(data)
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Println("[WARNING]: Could not decode the song image.")
		return
	}
	printImage(image, char)
}

func checkImage() {
	if r := recover(); r != nil {
		printCenter("[NO IMAGE]")
	}
}

func printImage(img image.Image, char string) {
	width, height, _ := terminal.GetSize(int(os.Stdin.Fd()))
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

func printSongDetails(p *Player) {
	printCenter(p.metadata.Title())
	printCenter(p.metadata.Artist())
	printCenter(strconv.Itoa(p.metadata.Year()))
	printCenter(p.metadata.Genre())
	printCenter(fmt.Sprintf("%d:%02d", p.duration/60, p.duration%60))
}

func printDuration(p *Player, s *Settings) chan struct{} {
	width, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("[WARNING]: Could not get the width of terminal, therefore cannot show the song progress bar. %v\n", err)
		return nil
	}
	interval := p.duration * 1000 / width
	quit := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				if p.player.IsPlaying() {
					fmt.Print(s.progressbarChar)
				}
			case <-quit:
				fmt.Print("\n")
				return
			}
		}
	}()
	return quit
}

func printCenter(text string) {
	text = strings.TrimSpace(text)
	if text == "" || text == "0" {
		return
	}
	width, _, _ := terminal.GetSize(int(os.Stdin.Fd()))
	length := runewidth.StringWidth(text)
	for i := 0; i < (width-length)/2; i++ {
		fmt.Printf(" ")
	}
	fmt.Println(text)
}
