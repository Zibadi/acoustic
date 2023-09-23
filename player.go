package main

import (
	"log"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

func newContext() (*oto.Context, chan struct{}) {
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		log.Fatal(err)
	}
	return otoCtx, readyChan
}

func newPlayer(mp3 *mp3.Decoder, context *oto.Context, readyChan chan struct{}) *oto.Player {
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan
	player := context.NewPlayer(mp3)
	return player
}
