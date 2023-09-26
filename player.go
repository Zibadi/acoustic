package main

import (
	"log"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

func newContext() (*oto.Context, chan struct{}) {
	option := newContextOption()
	context, readyChan, err := oto.NewContext(option)
	if err != nil {
		log.Fatal("could not create a new context for playing the song.", err)
	}
	return context, readyChan
}

func newContextOption() *oto.NewContextOptions {
	options := &oto.NewContextOptions{}
	options.SampleRate = 44100
	options.ChannelCount = 2
	options.Format = oto.FormatSignedInt16LE
	return options
}

func waitFor(readyChan chan struct{}) {
	// It might take a bit for the hardware audio devices to be ready,
	// so we wait on the channel.
	<-readyChan
}

func newPlayer(mp3 *mp3.Decoder, context *oto.Context) *oto.Player {
	player := context.NewPlayer(mp3)
	return player
}
