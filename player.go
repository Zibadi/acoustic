package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Player struct {
	index    int
	volume   float64
	isPaused bool
	songs    []*Song
	context  *audio.Context
	player   *audio.Player
}

func newPlayer(s *Settings) *Player {
	player := &Player{
		index:    0,
		volume:   1.0,
		isPaused: false,
		songs:    loadSongs(s),
		context:  newContext(),
	}
	return player
}

func (p *Player) getNextSong() *Song {
	p.index %= len(p.songs)
	if p.index < 0 {
		p.index = len(p.songs) - 1
	}
	return p.songs[p.index]
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
