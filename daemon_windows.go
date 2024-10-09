package daemon

import "errors"

func catchSignals() {
	log.Printf("info: [daemon] without forking")
}

func ForkChild(addr string, ln net.Listener) (*os.Process, error) {
	return nil, errors.New("daemon: fork is not supported on Windows")
}
