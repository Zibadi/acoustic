package main

type Status struct {
	isPaused          bool
	isAutoPaused      bool
	isAutoPauseEnable bool
	isShuffled        bool
	isGoingForward    bool
	isFinished        chan bool
}

func newStatus() Status {
	return Status{
		isPaused:          false,
		isAutoPaused:      false,
		isAutoPauseEnable: true,
		isShuffled:        false,
		isGoingForward:    true,
		isFinished:        make(chan bool, 1),
	}
}
