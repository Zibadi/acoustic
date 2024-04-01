package main

import (
	"os"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

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
		p.nextMusic()
	case "p":
		p.previousMusic()
	case "s":
		p.shuffle()
	case "h":
		p.toggleIsHot()
	case "q":
		os.Exit(0)
	}
}

func handelOtherKeys(key *keys.Key, p *Player) {
	switch key.Code {
	case keys.Space:
		p.autoPauseTicker.Stop()
		p.togglePauseOrPlay()
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
