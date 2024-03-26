package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mattn/go-tty"
)

func listen(p *Player) {
	p.player.Play()
	defer p.player.Close()
	key := make(chan rune)
	go readKey(p, key)

	for p.player.IsPlaying() || p.isPaused {
		select {
		case r := <-key:
			switch string(r) {
			case " ":
				if !p.isPaused {
					p.player.Pause()
				} else {
					p.player.Play()
				}
				p.isPaused = !p.isPaused
			case "n":
				p.index++
				return
			case "p":
				p.index--
				return
			case "A":
				p.volume = math.Min(2, p.volume+0.2)
				p.player.SetVolume(p.volume)
			case "B":
				p.volume = math.Max(0, p.volume-0.2)
				p.player.SetVolume(p.volume)
			case "C":
				p.player.SetPosition(p.player.Position() + (time.Second * 5))
			case "D":
				p.player.SetPosition(p.player.Position() - (time.Second * 5))
			case "q":
				os.Exit(0)
			}
		default:
			continue
		}
	}
	p.index++
}

func readKey(p *Player, key chan<- rune) {
	tty, err := tty.Open()
	if err != nil {
		fmt.Println("[WARNING]: Could not open the tty, therefore cannot listen to the key events.", err)
	}
	defer tty.Close()
	for p.player.IsPlaying() || p.isPaused {
		r, err := tty.ReadRune()
		if err != nil {
			fmt.Println("[WARNING]: Could not read the keyboard event.", err)
		}
		key <- r
	}
}
