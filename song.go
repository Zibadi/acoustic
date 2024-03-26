package main

import (
	"fmt"
	"os"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Song struct {
	path string
}

func loadSongs(settings *Settings) []*Song {
	entries := readDir(settings.dir)
	songsPath := getSongsPath(settings.dir, entries)
	songs := newSongs(songsPath)
	return songs
}

func readDir(dir string) []os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("[ERROR]: %v", err)
		os.Exit(0)
	}
	return entries
}

func getSongsPath(dir string, entries []os.DirEntry) []string {
	paths := []string{}
	for _, e := range entries {
		// TODO: Read songs of the sub directories too
		if !e.IsDir() {
			fileName := e.Name()
			paths = append(paths, dir+fileName)
		}
	}
	return paths
}

func newSongs(paths []string) []*Song {
	songs := []*Song{}
	for _, p := range paths {
		songs = append(songs, &Song{path: p})
	}
	return songs
}

func readSongMetadata(file *os.File) (tag.Metadata, error) {
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func getSongDuration(s *mp3.Stream) int {
	const sampleRate = 44100
	const sampleSize = 4
	samples := s.Length() / sampleSize
	duration := int(samples / int64(sampleRate))
	return duration
}
