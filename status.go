package main

type Status struct {
	isPaused       bool
	isShuffled     bool
	isGoingForward bool
	isFinished     chan bool
}

func newStatus() *Status {
	return &Status{
		isPaused:       false,
		isShuffled:     false,
		isGoingForward: true,
		isFinished:     make(chan bool, 1),
	}
}
