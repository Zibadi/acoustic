package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-tty"
)

func listen(p *Player, key <-chan rune) {
	for p.player.IsPlaying() || p.isPaused {
		timeout := getSongTimeout(p)
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
			case "s":
				p.shuffle()
			case "q":
				os.Exit(0)
			}
		case <-time.After(timeout):
			p.nextSong()
			return
		}
	}
	p.nextSong()
}

func readKey(key chan<- rune) {
	tty, err := tty.Open()
	if err != nil {
		fmt.Println("[WARNING]: Could not open the tty, therefore cannot listen to the key events.", err)
	}
	defer tty.Close()
	for {
		r, err := tty.ReadRune()
		if err != nil {
			fmt.Println("[WARNING]: Could not read the keyboard event.", err)
		}
		key <- r
	}
}
