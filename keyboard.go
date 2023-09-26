package main

import (
	"log"
	"math"
	"os"

	"github.com/ebitengine/oto/v3"
	"github.com/mattn/go-tty"
)

func listen(p *oto.Player, quitTicker chan struct{}) int {
	tty, err := tty.Open()
	if err != nil {
		log.Println("could not open the tty, therefore cannot listen to the key events.", err)
	}
	defer tty.Close()
	p.Play()
	isPaused := false
	volume := 1.0
	key := make(chan rune)
	go func() {
		for p.IsPlaying() || isPaused {
			r, err := tty.ReadRune()
			if err != nil {
				log.Println("could not read the key.", err)
			}
			key <- r
		}
	}()
	for p.IsPlaying() || isPaused {
		select {
		case r := <-key:
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
				close(quitTicker)
				return 1
			case "p":
				p.Close()
				close(quitTicker)
				return -1
			case "A":
				volume = math.Min(1, volume+0.2)
				p.SetVolume(volume)
			case "B":
				volume = math.Max(0, volume-0.2)
				p.SetVolume(volume)
			case "q":
				p.Close()
				close(quitTicker)
				os.Exit(0)
			}
		default:
			continue
		}
	}
	return 1
}
