package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/dhowden/tag"
	"github.com/mattn/go-runewidth"
	"github.com/nfnt/resize"
)

func run(p *Player) {
	go listenToKeyboard(p)
	for len(p.musics) > 0 {
		err := p.play()
		if err != nil {
			p.skipMusic()
		}
	}
}

func printMetadata(p *Player) {
	file, _ := os.Open(p.getCurrentMusic().path)
	defer file.Close()
	metadata, err := readMusicMetadata(file)
	if err != nil {
		fmt.Printf("[WARNING]: Could not load the music metadata of %v\n%v\n", p.getCurrentMusic().path, err)
	} else {
		printMusicImage(metadata, p.settings.imageChar)
		printMusicMetadata(metadata)
	}
}

func printMusicImage(m tag.Metadata, c string) {
	defer checkImage()
	data := m.Picture().Data
	reader := bytes.NewReader(data)
	image, _, err := image.Decode(reader)
	if err != nil {
		fmt.Println("[WARNING]: Could not decode the music image.")
		return
	}
	printImage(image, c)
}

func checkImage() {
	if r := recover(); r != nil {
		printlnCenter("[NO IMAGE]")
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
	printlnCenter(m.Title())
	printlnCenter(m.Artist())
	printlnCenter(m.Genre())
	printlnCenter(strconv.Itoa(m.Year()))
}

func printStatusSpace() {
	fmt.Println()
}

func printStatus(p *Player) {
	moveCursorToTagsLine()
	status := ""
	status += getCoolTag(p)
	status += getShuffleTag(p.status.isShuffled)
	status += getPuaseTag(p.status.isPaused)
	printCenter(status)
	moveCursorToProgressbarLine()
	updateProgressBar(p)
}

func printVolume(p *Player) {
	moveCursorToTagsLine()
	content := "["
	maxVolume := 2.0
	for i := 0.0; i < maxVolume; i += 0.1 {
		if i < p.volume {
			content += "❚"
		} else {
			content += " "
		}
	}
	content += "]"
	printCenter(content)
	moveCursorToProgressbarLine()
	updateProgressBar(p)
}

func moveCursorToTagsLine() {
	cursor.Hide()
	cursor.Up(1)
	cursor.ClearLine()
	cursor.StartOfLine()
}

func printDuration(p *Player) {
	minutes := int(p.duration.Seconds()) / 60
	seconds := int(p.duration.Seconds()) % 60
	duratoin := fmt.Sprintf("[%d:%02d]", minutes, seconds)
	printlnCenter(duratoin)
}

func getCoolTag(p *Player) string {
	music := p.getCurrentMusic()
	if music.isCool {
		return "[COOL]"
	}
	return ""
}

func getPuaseTag(isPaused bool) string {
	if isPaused {
		return "[PUASE]"
	}
	return ""
}

func getShuffleTag(isShuffled bool) string {
	if isShuffled {
		return "[SHUFFLE]"
	}
	return ""
}

func moveCursorToProgressbarLine() {
	cursor.Down(1)
}

func printProgressbar(p *Player) {
	if !p.status.isPaused {
		fmt.Print(p.settings.progressbarChar)
	}
}

func updateProgressBar(p *Player) {
	interval := getProgressbarInterval(p).Milliseconds()
	cursor.Hide()
	area := cursor.NewArea()
	area.ClearLinesDown(0)
	bound := p.player.Position().Milliseconds()
	content := ""
	for i := 0; i < int(bound/interval); i++ {
		content += p.settings.progressbarChar
	}
	area.Update(content)
	cursor.Show()
}

func printCenter(text string) {
	err := centerizeCursor(text)
	if err != nil {
		return
	}
	fmt.Print(text)
}

func printlnCenter(text string) {
	err := centerizeCursor(text)
	if err != nil {
		return
	}
	fmt.Println(text)
}

func centerizeCursor(text string) error {
	text = strings.TrimSpace(text)
	if text == "" || text == "0" {
		return errors.New("input text is empty")
	}
	width, _, _ := getTerminalSize()
	length := runewidth.StringWidth(text)
	for i := 0; i < (width-length)/2; i++ {
		fmt.Print(" ")
	}
	return nil
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
