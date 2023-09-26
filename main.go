package main

import (
	"log"
	"os"
)

func main() {
	checkArgs(os.Args)
	dir := os.Args[1]
	entries := openDir(dir)
	paths := getPaths(dir, entries)
	paths = shuffle(paths)
	index := 0
	context := newContext()
	for {
		index %= len(paths)
		if index < 0 {
			index = len(paths) - 1
		}
		file, err := openFile(paths[index])
		if err != nil {
			log.Printf("could not open the %v. %v\n", paths[index], err)
			index++
			continue
		}
		defer file.Close()
		metadata, err := readMetadata(file)
		if err != nil {
			log.Println("there is no any info tag for this song.", err)
		} else {
			showImage(metadata, "♥")
			showMetadata(metadata)
		}
		stream, err := decode(file)
		if err != nil {
			log.Printf("could not decode the %v. %v\n", paths[index], err)
			index++
			continue
		}
		// waitFor(readyChan)
		player, err := newPlayer(stream, context)
		if err != nil {
			log.Printf("could not create a new player for the %v. %v\n", paths[index], err)
			index++
			continue
		}
		defer player.Close()
		quit := showProgressBar(stream, player)
		index += listen(player)
		close(quit)
	}
}
