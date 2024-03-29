package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Player struct {
	index           int
	volume          float64
	duration        int
	isPaused        bool
	isGoingForward  bool
	metadata        tag.Metadata
	autoPuaseTicker *time.Ticker
	songs           []Song
	context         *audio.Context
	player          *audio.Player
}

func newPlayer(s *Settings) Player {
	player := Player{
		index:           0,
		volume:          1.0,
		isPaused:        false,
		isGoingForward:  true,
		autoPuaseTicker: time.NewTicker(time.Duration(1) * time.Second),
		songs:           loadSongs(s),
		context:         newContext(),
	}
	return player
}

func newContext() *audio.Context {
	const sampleRate = 44100
	context := audio.NewContext(sampleRate)
	return context
}

func (p *Player) play(s *Settings) error {
	song, err := p.preparePlayer()
	if err != nil {
		return err
	}
	defer song.Close()
	printMetadata(p, s)
	quit := printDuration(p, s)
	defer close(quit)
	p.player.Play()
	defer p.player.Close()
	listen(p)
	return nil
}

func (p *Player) preparePlayer() (*os.File, error) {
	song := p.getCurrentSong()
	file, err := os.Open(song.path)
	if err != nil {
		fmt.Printf("[ERROR]: Could not open the %v\n%v\n", song.path, err)
		return nil, err
	}
	stream, err := decode(file)
	if err != nil {
		fmt.Printf("[ERROR]: Could not decode %v\n%v\n", song.path, err)
		return file, err
	}
	p.player, err = p.context.NewPlayer(stream)
	p.player.SetVolume(p.volume)
	if err != nil {
		fmt.Printf("[ERROR]: Could not play %v\n%v\n", song.path, err)
		return file, err
	}
	p.duration = getSongDuration(stream)
	return file, nil
}

func (p *Player) getCurrentSong() Song {
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

func (p *Player) skipSong() {
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

func (p *Player) autoPuase() {
	count, err := getRunningAudioCount()
	if err != nil {
		p.autoPuaseTicker.Stop()
		return
	}
	if (count > 2 && !p.isPaused) || (count <= 2 && p.isPaused) {
		p.togglePuaseOrPlay()
	}
}

func getRunningAudioCount() (int, error) {
	dump := exec.Command("pw-dump")
	output, err := dump.Output()
	if err != nil {
		return 0, err
	}
	return strings.Count(string(output), "running"), nil
}
