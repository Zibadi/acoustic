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
		log.Fatal(err)
	}
	return entries
}

func createPaths(dir string, entries []os.DirEntry) []string {
	names := make([]string, len(entries))
	for i, e := range entries {
		fileName := e.Name()
		names[i] = dir + fileName
	}
	return names
}

func openFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func readMetadata(file *os.File) tag.Metadata {
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	return metadata
}

func decode(file *os.File) *mp3.Decoder {
	mp3, err := mp3.NewDecoder(file)
	if err != nil {
		log.Fatal(err)
	}
	return mp3
}
