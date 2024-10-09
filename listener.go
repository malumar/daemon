package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

// TcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type listener struct {
	Addr     string `json:"addr"`
	FD       int    `json:"fd"`
	Filename string `json:"filename"`
}

func NewListener(address string) (*TcpKeepAliveListener, bool, error) {
	oldListner := os.Getenv("LISTENER")

	log.Printf("info: [daemon] LISTENER env=%v\n", oldListner)
	if oldListner != "" {
		log.Printf("info: [daemon] Importowanie listnera %s\n", oldListner)
		if ln, err := ImportListener(address); err != nil {
			return nil, false, err
		} else {

			return &TcpKeepAliveListener{ln.(*net.TCPListener)}, true, nil
		}
	}
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, false, err
	}
	return &TcpKeepAliveListener{l.(*net.TCPListener)}, false, nil
}

func CreateOrImportListener(addr string) (net.Listener, error) {
	// Try and import a listener for addr. If it's found, use it.
	ln, err := ImportListener(addr)
	if err == nil {
		fmt.Printf("Imported listener file descriptor for %v.\n", addr)
		return ln, nil
	}

	// No listener was imported, that means this process has to create one.
	ln, err = CreateListener(addr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created listener file descriptor for %v.\n", addr)
	return ln, nil
}

func ImportListener(addr string) (net.Listener, error) {
	// Extract the encoded listener metadata from the environment.
	listenerEnv := os.Getenv("LISTENER")
	if listenerEnv == "" {
		return nil, fmt.Errorf("unable to find LISTENER environment variable")
	}

	// Unmarshal the listener metadata.
	var l listener
	err := json.Unmarshal([]byte(listenerEnv), &l)
	if err != nil {
		return nil, err
	}
	if l.Addr != addr {
		return nil, fmt.Errorf("unable to find listener for %v l.Addr= %v ", addr, l.Addr)
	}

	// The file has already been passed to this process, extract the file
	// descriptor and name from the metadata to rebuild/find the *os.File for
	// the listener.
	listenerFile := os.NewFile(uintptr(l.FD), l.Filename)
	if listenerFile == nil {
		return nil, fmt.Errorf("unable to create listener file: %v", err)
	}
	defer listenerFile.Close()

	// Tell the parent to stop the server now.
	parent := syscall.Getppid()
	log.Printf("info: [daemon] %d Telling parent process (%d) to stop server\n", syscall.Getpid(), parent)
	syscall.Kill(parent, syscall.SIGTERM)

	// Give the parent some time.
	time.Sleep(100 * time.Millisecond)

	ln, err := net.FileListener(listenerFile)
	if err != nil {
		// INFO.Printf("Auuu", syscall.Getpid())
		return nil, err
	}

	return ln, nil
}

func CreateListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func getListenerFile(ln net.Listener) (*os.File, error) {
	switch t := ln.(type) {
	case *net.TCPListener:
		return t.File()
	case *net.UnixListener:
		return t.File()
	case *TcpKeepAliveListener:
		return t.File()
	}
	return nil, fmt.Errorf("unsupported listener: %T", ln)
}
