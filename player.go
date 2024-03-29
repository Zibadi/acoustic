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
	isFinished      chan bool
	metadata        tag.Metadata
	autoPauseTicker *time.Ticker
	musics          []Music
	context         *audio.Context
	player          *audio.Player
}

func newPlayer(s *Settings) Player {
	player := Player{
		index:           0,
		volume:          1.0,
		isPaused:        false,
		isGoingForward:  true,
		isFinished:      make(chan bool, 1),
		autoPauseTicker: time.NewTicker(time.Duration(1) * time.Second),
		musics:          loadMusics(s),
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
	music, err := p.preparePlayer()
	if err != nil {
		return err
	}
	defer music.Close()
	printMetadata(p, s)
	quit := printDuration(p, s)
	defer close(quit)
	p.player.Play()
	defer p.player.Close()
	p.listen()
	return nil
}

func (p *Player) listen() {
	for {
		timeout := getMusicTimeout(p)
		select {
		case <-p.isFinished:
			return
		case <-p.autoPauseTicker.C:
			p.autoPause()
		case <-time.After(timeout):
			p.nextMusic()
		}
	}
}

func (p *Player) preparePlayer() (*os.File, error) {
	music := p.getCurrentMusic()
	file, err := os.Open(music.path)
	if err != nil {
		fmt.Printf("[ERROR]: Could not open the %v\n%v\n", music.path, err)
		return nil, err
	}
	stream, err := decode(file)
	if err != nil {
		fmt.Printf("[ERROR]: Could not decode %v\n%v\n", music.path, err)
		return file, err
	}
	p.player, err = p.context.NewPlayer(stream)
	p.player.SetVolume(p.volume)
	if err != nil {
		fmt.Printf("[ERROR]: Could not play %v\n%v\n", music.path, err)
		return file, err
	}
	p.duration = getMusicDuration(stream)
	return file, nil
}

func (p *Player) getCurrentMusic() Music {
	p.index %= len(p.musics)
	return p.musics[p.index]
}

func (p *Player) nextMusic() {
	p.index++
	p.isGoingForward = true
	p.isFinished <- true
}

func (p *Player) previousMusic() {
	p.index--
	if p.index < 0 {
		p.index = len(p.musics) - 1
	}
	p.isGoingForward = false
	p.isFinished <- true
}

func (p *Player) skipMusic() {
	p.musics[p.index] = p.musics[len(p.musics)-1]
	if !p.isGoingForward {
		p.previousMusic()
	}
	p.musics = p.musics[:len(p.musics)-1]
}

func (p *Player) togglePauseOrPlay() {
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
	p.musics = shuffle(p.musics)
}

func (p *Player) autoPause() {
	count, err := getRunningAudioCount()
	if err != nil {
		p.autoPauseTicker.Stop()
		return
	}
	if (count > 2 && !p.isPaused) || (count <= 2 && p.isPaused) {
		p.togglePauseOrPlay()
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
