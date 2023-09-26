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

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/mattn/go-runewidth"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/ssh/terminal"
)

func showImage(m tag.Metadata, char string) {
	defer func() {
		if r := recover(); r != nil {
			printCenter("[NO IMAGE]")
		}
	}()
	data := m.Picture().Data
	reader := bytes.NewReader(data)
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Println("could not decode the image.")
		return
	}
	print(image, char)
}

func print(img image.Image, char string) {
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
			fmt.Printf("\033[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)

			r, g, b, _ = image.At(x, y+1).RGBA()
			fmt.Printf("\033[38;2;%d;%d;%dm▄", r>>8, g>>8, b>>8)
		}
		fmt.Printf("\033[0m")
		fmt.Printf("\n")
	}
}

func showMetadata(m tag.Metadata) {
	printCenter(m.Title())
	printCenter(m.Artist())
	printCenter(strconv.Itoa(m.Year()))
	printCenter(m.Genre())
}

func showProgressBar(s *mp3.Stream, p *audio.Player) chan struct{} {
	const sampleRate = 44100
	const sampleSize = 4                  // From documentation.
	samples := s.Length() / sampleSize    // Number of samples.
	length := int(samples / int64(44100)) // Audio length in seconds.
	printCenter(fmt.Sprintf("%d:%02d", length/60, length%60))
	width, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("could not get the width of terminal, therefore cannot show the song length. %v\n", err)
	}
	duration := length * 1000 / width
	quit := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
		counter := 0
		for {
			select {
			case <-ticker.C:
				if counter == width {
					continue
				} else if p.IsPlaying() {
					fmt.Printf("─")
					counter++
				}
			case <-quit:
				fmt.Printf("\n")
				ticker.Stop()
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
