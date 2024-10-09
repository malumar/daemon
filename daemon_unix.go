//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package daemon

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func catchSignals() {
	signalChannel := make(chan os.Signal, 1024)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	msg := <-signalChannel
	if msg == syscall.SIGINT {
		// once we get this we should probably use killing all processes
		log.Printf("err: [daemon] kill all processes")
	}
	MessageWithArgs(FINISH_APP_GRACEFULLY, map[string]string{UseSignalParam: UseSignalParamTrue})
}

// There is one problem with forking - when you do this, the program will detach from the console
// and the other process will, in fact, the program will run in the background and spit on the console,
// so in such a situation, when you use it, you should:
//
// 1. Fork the first program you run as well
// 2. Save logs in the logs folder
func ForkChild(addr string, ln net.Listener) (*os.Process, error) {
	// Get the file descriptor for the listener and marshal the metadata to pass
	// to the child in the environment.
	lnFile, err := getListenerFile(ln)
	if err != nil {
		return nil, err
	}
	defer lnFile.Close()
	l := listener{
		Addr:     addr,
		FD:       3,
		Filename: lnFile.Name(),
	}
	listenerEnv, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	// Pass stdin, stdout, and stderr along with the listener to the child.
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
		lnFile,
	}

	// Get current environment and add in the listener to it.
	environment := append(os.Environ(), "LISTENER="+string(listenerEnv), fmt.Sprintf("%s=%s", EnvWorkerName, EnvWorkerNameValue))

	// Get current process name and directory.
	execName, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execName)
	// Spawn child process.
	p, err := os.StartProcess(execName, append([]string{execName}, flag.Args()...), &os.ProcAttr{
		Dir:   execDir,
		Env:   environment,
		Files: files,
		Sys: &syscall.SysProcAttr{
			Setpgid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}
