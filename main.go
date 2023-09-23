package main

import (
	"fmt"
	"log"
	"os"
)

func checkArgs(args []string) {
	if len(args) < 2 {
		log.Fatal("specify the path of the file.")
	}
}

func main() {
	checkArgs(os.Args)
	dir := os.Args[1]
	entries := openDir(dir)
	paths := createPaths(dir, entries)
	paths = shuffle(paths)
	index := 0
	context, readyChan := newContext()
	for {
		file := openFile(paths[index])
		defer file.Close()
		metadata := readMetadata(file)
		fmt.Printf("[ARTIST]: %v, [SONG]: %v\n", metadata.Artist(), metadata.Title())
		mp3 := decode(file)
		player := newPlayer(mp3, context, readyChan)
		index += listen(player)
	}
}
