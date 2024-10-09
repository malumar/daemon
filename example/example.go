package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/malumar/daemon"
	"github.com/malumar/daemon/example/stuff"
	"github.com/malumar/daemon/plugins"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {

	if dir, err := os.MkdirTemp("/tmp", "daemon-pids*"); err != nil {
		log.Fatalf(err.Error())
	} else {
		tmpDir = dir
	}
	log.Printf("[main] temp path %s\n", tmpDir)
	// we will inform listeners that this is the second stage of initialization
	daemon.FireTrigger(plugins.TRIGGER_BEFORE_FLAGS_PREPARED, nil)
	flag.Parse()

	// inform listeners that the application is in the process of initialization
	daemon.FireTrigger(plugins.TRIGGER_BEFORE_APP_INIT, nil)

	stuff.DoTheReallyImportantThings()

	daemon.RegisterListener(plugins.TRIGGER_RELOAD_APP, func(ctx context.Context, event string, value interface{}) error {
		fmt.Printf("[main] the application has reloaded")
		return nil
	})

	// inform listeners that the application has been initialized
	daemon.FireTrigger(plugins.TRIGGER_AFTER_APP_INIT, nil)

	daemon.RegisterListener(plugins.TRIGGER_APPLICATION_EXIT, func(ctx context.Context, event string, value interface{}) error {
		fmt.Println("[main] the application has ended")
		return nil
	})

	// listen for this trigger to be called
	daemon.RegisterListener(plugins.TRIGGER_FINISH_APP, stopApp)

	run()

}

var listener net.Listener
var srv *http.Server
var tmpDir string

// stopApp indirect method to stop the application after it starts,
// do not call it directly, only via daemon.Stop();
func stopApp(ctx context.Context, event string, value interface{}) error {
	fmt.Println("[main] someone triggered plugins.TRIGGER_FINISH_APP ")
	srv.Close()
	return nil
}

func run() {

	var addr = ":8080"
	srv = &http.Server{Addr: addr}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "text/plain")
		writer.Write([]byte("Hello World"))
	})

	http.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		daemon.Stop()
	})

	http.HandleFunc("/reload", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "text/plain")
		writer.Write([]byte("Reload app"))
		daemon.Reload()
	})
	go func() {

		if listener == nil {
			if ln, _, err := daemon.NewListener(addr); err != nil {
				log.Fatal(err)
			} else {
				listener = ln
			}

		}

		// always returns error. ErrServerClosed on graceful close
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		} else {
			fmt.Println("[main] the server is shutting down")
		}

	}()

	go func() {
		fmt.Println("[main] if no one stops the application myself, I will do it in a few seconds")
		time.Sleep(time.Second * 5)
		fmt.Println("[main] stopping the application")
		daemon.Stop()

	}()

	daemon.Run(tmpDir, func() {
		fmt.Println("[main] someone interrupted the application with CTRL C or killed our process with kill, " +
			"we have to stop everything manually")

	})

	return
}
