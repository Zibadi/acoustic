package main

import (
	"math"
	"os"
	"time"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Player struct {
	index          int
	volume         float64
	duration       int
	isPaused       bool
	isGoingForward bool
	metadata       tag.Metadata
	songs          []*Song
	context        *audio.Context
	player         *audio.Player
}

func newPlayer(s *Settings) *Player {
	player := &Player{
		index:          0,
		volume:         1.0,
		isPaused:       false,
		isGoingForward: true,
		songs:          loadSongs(s),
		context:        newContext(),
	}
	return player
}

func (p *Player) getCurrentSong() *Song {
	p.index %= len(p.songs)
	return p.songs[p.index]
}

func (p *Player) nextSong() {
	p.index++
	p.isGoingForward = true
}

func (p *Player) previousSong() {
	p.index--
	if p.index < 0 {
		p.index = len(p.songs) - 1
	}
	p.isGoingForward = false
}

func newContext() *audio.Context {
	const sampleRate = 44100
	context := audio.NewContext(sampleRate)
	return context
}

func decode(f *os.File) (*mp3.Stream, error) {
	stream, err := mp3.DecodeWithSampleRate(44100, f)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func skipSong(p *Player) {
	p.songs[p.index] = p.songs[len(p.songs)-1]
	if !p.isGoingForward {
		p.previousSong()
	}
	p.songs = p.songs[:len(p.songs)-1]
}

func (p *Player) togglePuaseOrPlay() {
	if !p.isPaused {
		p.player.Pause()
	} else {
		p.player.Play()
	}
	p.isPaused = !p.isPaused
}

func (p *Player) increaseVolume() {
	p.volume = math.Min(2, p.volume+0.2)
	p.player.SetVolume(p.volume)
}

func (p *Player) decreaseVolume() {
	p.volume = math.Max(0, p.volume-0.2)
	p.player.SetVolume(p.volume)
}

func (p *Player) seekForward() {
	newPosition := p.player.Position() + (time.Second * 5)
	if int(newPosition.Seconds()) < p.duration {
		p.player.SetPosition(newPosition)
	}
}

func (p *Player) seekBackward() {
	newPosition := p.player.Position() - (time.Second * 5)
	if int(newPosition.Seconds()) > 0 {
		p.player.SetPosition(newPosition)
	} else {
		p.player.SetPosition(0)
	}
}

func (p *Player) shuffle() {
	p.songs = shuffle(p.songs)
}
