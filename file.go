package main

import (
	"log"
	"os"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
)

func openDir(dir string) []os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("could not find %v. %v", dir, err)
	}
	return entries
}

func createPaths(dir string, entries []os.DirEntry) []string {
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

func decode(f *os.File) (*mp3.Decoder, error) {
	mp3, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, err
	}
	return mp3, nil
}
