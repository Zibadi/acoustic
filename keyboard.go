package main

import (
	"fmt"
	"math"
	"os"

	"github.com/mattn/go-tty"
)

func listen(p *Player) {
	tty, err := tty.Open()
	if err != nil {
		fmt.Println("[WARNING]: Could not open the tty, therefore cannot listen to the key events.", err)
	}
	defer tty.Close()
	p.player.Play()
	key := make(chan rune)
	go func() {
		for p.player.IsPlaying() || p.isPaused {
			r, err := tty.ReadRune()
			if err != nil {
				fmt.Println("[WARNING]: Could not read the keyboard event.", err)
			}
			key <- r
		}
	}()
	for p.player.IsPlaying() || p.isPaused {
		select {
		case r := <-key:
			switch string(r) {
			case " ":
				if !p.isPaused {
					p.player.Pause()
					p.isPaused = true
				} else {
					p.player.Play()
					p.isPaused = false
				}
			case "n":
				p.player.Close()
				p.index++
				return
			case "p":
				p.player.Close()
				p.index--
				return
			case "A":
				p.volume = math.Min(1, p.volume+0.2)
				p.player.SetVolume(p.volume)
			case "B":
				p.volume = math.Max(0, p.volume-0.2)
				p.player.SetVolume(p.volume)
			case "q":
				p.player.Close()
				os.Exit(0)
			}
		default:
			continue
		}
	}
	p.index++
}
