package main

import (
	_ "image/gif"  // initialize decoder
	_ "image/jpeg" // initialize decoder
	_ "image/png"  // initialize decoder
	"log"
	"os"
)

func main() {
	checkArgs(os.Args)
	dir := os.Args[1]
	entries := openDir(dir)
	paths := createPaths(dir, entries)
	paths = shuffle(paths)
	index := 0
	context, readyChan := newContext()
	for {
		index %= len(paths)
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
		mp3, err := decode(file)
		if err != nil {
			log.Printf("could not decode the %v. %v\n", paths[index], err)
			index++
			continue
		}
		waitFor(readyChan)
		player := newPlayer(mp3, context)
		quitTicker, err := showLength(mp3, player)
		if err != nil {
			log.Printf("could not get the width of terminal, therefore cannot show the song length. %v\n", err)
		}
		index += listen(player, quitTicker)
	}
}
