package main

import (
	"os"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func listen(p *Player) {
	listenToKeyboard(p)
	listenToSong(p)
}

func listenToKeyboard(p *Player) {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.RuneKey {
			handleRuneKeys(&key, p)
		} else {
			handelOtherKeys(&key, p)
		}
		return false, nil // Return false to continue listening
	})
}

func handleRuneKeys(key *keys.Key, p *Player) {
	switch key.String() {
	case "n":
		p.nextSong()
		return
	case "p":
		p.previousSong()
		return
	case "s":
		p.shuffle()
	case "q":
		os.Exit(0)
	}
}

func handelOtherKeys(key *keys.Key, p *Player) {
	switch key.Code {
	case keys.Space:
		p.autoPuaseTicker.Stop()
		p.togglePuaseOrPlay()
	case keys.Up:
		p.increaseVolume()
	case keys.Down:
		p.decreaseVolume()
	case keys.Right:
		p.seekForward()
	case keys.Left:
		p.seekBackward()
	case keys.CtrlC:
		os.Exit(0)
	}
}

func listenToSong(p *Player) {
	for p.player.IsPlaying() || p.isPaused {
		timeout := getSongTimeout(p)
		select {
		case <-p.autoPuaseTicker.C:
			p.autoPuase()
		case <-time.After(timeout):
			p.nextSong()
			return
		}
	}
	p.nextSong()
}
