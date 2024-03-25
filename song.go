package main

import (
	"log"
	"os"

	"github.com/dhowden/tag"
)

type Song struct {
	path     string
	length   int
	metadata *tag.Metadata
}

func loadSongs(settings *Settings) []*Song {
	entries := readDir(settings.dir)
	paths := getPaths(settings.dir, entries)
	songs := createSongs(paths)
	return songs
}

func readDir(dir string) []os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("could not find %v. %v", dir, err)
	}
	return entries
}

func getPaths(dir string, entries []os.DirEntry) []string {
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

func createSongs(paths []string) []*Song {
	songs := []*Song{}
	for _, p := range paths {
		songs = append(songs, &Song{path: p})
	}
	return songs
}

func openFile(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readMetadata(file *os.File) (tag.Metadata, error) {
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
