package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Music struct {
	path string
}

func loadMusics(settings *Settings) []Music {
	musics := make([]Music, 0)
	err := filepath.WalkDir(settings.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("[ERROR]: %v", err)
			return err
		}
		if !d.IsDir() {
			musics = append(musics, Music{path: path})
		}
		return nil
	})
	if err != nil {
		fmt.Printf("[ERROR]: Could not load the musics from %v\n%v\n", settings.dir, err)
		os.Exit(0)
	}
	return musics
}

func decode(f *os.File) (*mp3.Stream, error) {
	stream, err := mp3.DecodeWithSampleRate(44100, f)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func readMusicMetadata(file *os.File) (tag.Metadata, error) {
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func getMusicDuration(s *mp3.Stream) time.Duration {
	const sampleRate = 44100
	const sampleSize = 4
	samples := s.Length() / sampleSize
	duration := int(samples / int64(sampleRate))
	return time.Duration(duration * int(time.Second))
}

func getMusicTimeout(p *Player) time.Duration {
	if p.isPaused {
		return time.Duration(24*365) * time.Hour
	}
	return p.duration - p.player.Position()
}
