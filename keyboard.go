package main

import (
	"fmt"
	"os"

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
				p.togglePuaseOrPlay()
			case "n":
				p.nextSong()
				return
			case "p":
				p.previousSong()
				return
			case "A":
				p.increaseVolume()
			case "B":
				p.decreaseVolume()
			case "C":
				p.seekForward()
			case "D":
				p.seekBackward()
			case "q":
				os.Exit(0)
			}
		default:
			continue
		}
	}
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
