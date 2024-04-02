package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Player struct {
	index             int
	volume            float64
	duration          time.Duration
	status            Status
	settings          *Settings
	autoPauseTicker   *time.Ticker
	progressbarTicker *time.Ticker
	musics            []Music
	context           *audio.Context
	player            *audio.Player
}

func newPlayer(s *Settings) Player {
	player := Player{
		index:           0,
		volume:          1.0,
		settings:        s,
		status:          newStatus(),
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

func newProgressbarTicker(p *Player) *time.Ticker {
	interval := getProgressbarInterval(p)
	return time.NewTicker(interval)
}

func getProgressbarInterval(p *Player) time.Duration {
	width, _, _ := getTerminalSize()
	interval := p.duration.Milliseconds() / int64(width)
	return time.Duration(interval * int64(time.Millisecond))
}

func (p *Player) play() error {
	music, err := p.preparePlayer()
	if err != nil {
		return err
	}
	defer music.Close()
	printMetadata(p)
	printDuration(p)
	printStatusSpace()
	printStatus(p)
	p.player.Play()
	defer p.dispose()
	p.listen()
	return nil
}

func (p *Player) listen() {
	for {
		select {
		case <-p.status.isFinished:
			p.finished()
			return
		case <-p.progressbarTicker.C:
			printProgressbar(p)
			printStatus(p)
		case <-p.autoPauseTicker.C:
			p.autoPause()
		}
	}
}

func (p *Player) dispose() {
	p.player.Pause()
	p.player.Close()
	p.progressbarTicker.Stop()
}

func (p *Player) finished() {
	fmt.Println()
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
	if err != nil {
		fmt.Printf("[ERROR]: Could not play %v\n%v\n", music.path, err)
		return file, err
	}
	p.setupPlayerConfigs(stream)
	return file, nil
}

func (p *Player) setupPlayerConfigs(s *mp3.Stream) {
	p.player.SetVolume(p.volume)
	p.duration = getMusicDuration(s)
	p.progressbarTicker = newProgressbarTicker(p)
	p.status = newStatus()
}

func (p *Player) getCurrentMusic() *Music {
	p.index %= len(p.musics)
	return &p.musics[p.index]
}

func (p *Player) nextMusic() {
	if p.settings.isHotColdEnabled {
		p.moveToColdDir()
	}
	p.index++
	p.status.isGoingForward = true
	p.status.isFinished <- true
}

func (p *Player) previousMusic() {
	if p.settings.isHotColdEnabled {
		p.moveToColdDir()
	}
	p.index--
	if p.index < 0 {
		p.index = len(p.musics) - 1
	}
	p.status.isGoingForward = false
	p.status.isFinished <- true
}

func (p *Player) skipMusic() {
	p.musics[p.index] = p.musics[len(p.musics)-1]
	if !p.status.isGoingForward {
		p.previousMusic()
	}
	p.musics = p.musics[:len(p.musics)-1]
}

func (p *Player) disableAutoPause() {
	if p.status.isAutoPaused {
		p.autoPauseTicker.Stop()
		p.status.isAutoPaused = false
		p.status.isAutoPauseEnable = false
	}
	p.status.isAutoPauseEnable = !p.status.isAutoPauseEnable
}

func (p *Player) togglePauseOrPlay() {
	if !p.status.isPaused {
		p.player.Pause()
	} else {
		p.player.Play()
	}
	p.status.isPaused = !p.status.isPaused
	printStatus(p)
}

func (p *Player) increaseVolume() {
	p.volume = math.Min(2.0, p.volume+0.1)
	p.player.SetVolume(p.volume)
	printVolume(p)
}

func (p *Player) decreaseVolume() {
	p.volume = math.Max(0, p.volume-0.1)
	p.player.SetVolume(p.volume)
	printVolume(p)
}

func (p *Player) seekForward() {
	newPosition := p.player.Position() + (time.Second * 5)
	if newPosition < p.duration {
		p.player.SetPosition(newPosition)
		updateProgressBar(p)
	}
}

func (p *Player) seekBackward() {
	newPosition := p.player.Position() - (time.Second * 5)
	if int(newPosition.Seconds()) > 0 {
		p.player.SetPosition(newPosition)
	} else {
		p.player.SetPosition(0)
	}
	updateProgressBar(p)
}

func (p *Player) shuffle() {
	p.musics, p.index = shuffle(p.musics, p.index)
	p.status.isShuffled = true
	printStatus(p)
}

func (p *Player) autoPause() {
	if !p.status.isAutoPauseEnable {
		return
	}
	count, err := getRunningAudioCount()
	if err != nil {
		p.autoPauseTicker.Stop()
		return
	}
	if count > 2 && !p.status.isAutoPaused {
		p.status.isAutoPaused = true
		p.togglePauseOrPlay()
	} else if count <= 2 && p.status.isAutoPaused {
		p.status.isAutoPaused = false
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

func (p *Player) moveToColdDir() {
	music := p.getCurrentMusic()
	if music.isHot {
		return
	}
	fileName := filepath.Base(music.path)
	absolutePath, _ := filepath.Abs(p.settings.dir)
	baseDir := filepath.Dir(absolutePath)
	newPath := filepath.Join(baseDir, "Cold", fileName)
	if music.path == newPath {
		return
	}
	err := moveFile(music.path, newPath)
	if err != nil {
		return
	}
	music.path = newPath
}

func (p *Player) toggleIsHot() {
	if !p.settings.isHotColdEnabled {
		return
	}
	music := p.getCurrentMusic()
	fileName := filepath.Base(music.path)
	absolutePath, _ := filepath.Abs(p.settings.dir)
	baseDir := filepath.Dir(absolutePath)
	var newPath string
	if music.isHot {
		newPath = filepath.Join(baseDir, "Cold", fileName)
		music.isHot = false
	} else {
		newPath = filepath.Join(baseDir, "Hot", fileName)
		music.isHot = true
	}
	err := moveFile(music.path, newPath)
	if err != nil {
		return
	}
	music.path = newPath
	printStatus(p)
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		fmt.Printf("[ERROR]: Could not open source file: %v", err)
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("[ERROR]: Could not open dest file: %v", err)
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		fmt.Printf("[ERROR]: Could not copy to dest from source: %v", err)
		return err
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		fmt.Printf("[ERROR]: Could not remove source file: %v", err)
		return err
	}
	return nil
}
