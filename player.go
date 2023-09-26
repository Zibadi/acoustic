package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

func newContext() *audio.Context {
	const sampleRate = 44100
	context := audio.NewContext(sampleRate)
	return context
}

// func newContextOption() *oto.NewContextOptions {
// 	options := &oto.NewContextOptions{}
// 	options.SampleRate = 44100
// 	options.ChannelCount = 2
// 	options.Format = oto.FormatSignedInt16LE
// 	return options
// }

// func waitFor(readyChan chan struct{}) {
// 	// It might take a bit for the hardware audio devices to be ready,
// 	// so we wait on the channel.
// 	<-readyChan
// }

func newPlayer(s *mp3.Stream, c *audio.Context) (*audio.Player, error) {
	// player := context.NewPlayer(mp3)
	p, err := c.NewPlayer(s)
	return p, err
}

func decode(f *os.File) (*mp3.Stream, error) {
	stream, err := mp3.DecodeWithSampleRate(44100, f)
	if err != nil {
		return nil, err
	}
	return stream, nil
}
