package main

import (
	"log"

	"github.com/ebitengine/oto/v3"
	"github.com/mattn/go-tty"
)

func listen(p *oto.Player) int {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	p.Play()
	isPaused := false
	for p.IsPlaying() {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		switch string(r) {
		case " ":
			if !isPaused {
				p.Pause()
				isPaused = true
			} else {
				p.Play()
				isPaused = false
			}
		case "n":
			p.Close()
			return 1
		case "p":
			p.Close()
			return -1
		}
	}
	return 1
}
