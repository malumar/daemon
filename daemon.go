package daemon

import (
	"fmt"
	"github.com/malumar/daemon/plugins"
	"log"
	"sync/atomic"
	"syscall"
)

type Command int

const (
	FINISH_APP_GRACEFULLY Command = iota
	RELOAD_APP
	UseSignalParam      = "signal"
	UseSignalParamTrue  = "1"
	UseSignalParamFalse = ""
)

var channelMessage = make(chan ChannelMessage)
var finishHim = make(chan bool)

type ChannelMessage struct {
	Command Command
	Args    map[string]string
}

func (self ChannelMessage) Arg(name string, dfault string) string {
	if self.Args != nil {
		if val, ok := self.Args[name]; ok {
			return val
		}
	}

	return dfault
}

var hookAfterFinish func()
var pidsPathStorage string

// Run
// @pathForPids - where to save pids - create if it does not exist
// @handlerAfterFinish called at the end
func Run(pathForPids string, handlerAfterFinish func()) {

	addProcessPid(syscall.Getpid())

	hookAfterFinish = handlerAfterFinish
	pidsPathStorage = pathForPids
	//	GetWaitingGroup().Add(1)
	go func() {
		//		defer GetWaitingGroup().Done()
		catchSignals()

	}()
	waitForMessages()
}

// Reload Anyone can do this to make the server reload functions for end applications
func Reload() {
	// uwaga nie wspieramy
	Message(RELOAD_APP)
}

// DoReloadAppGraceful block accepting new connections and wait until all current connections are completed,
// then reload the application
// Deprecated: use Reload instead
func DoReloadAppGraceful() {
	channelMessage <- ChannelMessage{
		Command: RELOAD_APP,
	}
}

// Stop Anyone can do this to stop the server as a function for end applications
func Stop() {
	Message(FINISH_APP_GRACEFULLY)
}

func Message(cmd Command) {
	MessageWithArgs(cmd, nil)
}

func MessageWithArgs(cmd Command, args map[string]string) {
	channelMessage <- ChannelMessage{
		Command: cmd,
		Args:    args,
	}

}

func waitForMessages() {
	for {
		msg, ok := <-channelMessage
		if ok {
			switch msg.Command {
			case FINISH_APP_GRACEFULLY:
				// Force all processes to end
				FireTrigger(plugins.TRIGGER_FINISH_APP, nil)
				close(channelMessage)

				if msg.Arg(UseSignalParam, UseSignalParamFalse) == UseSignalParamTrue {
					stopAppGracefully(true)
				} else {
					stopAppGracefully(false)
				}

				return
			case RELOAD_APP:
				reloadAppGracefully()
				break

			default:
				panic(fmt.Sprintf("daemon unknown command: %v", msg.Command))
				break
			}
		} else {
			stopAppGracefully(false)
			return

		}
	}
}

// IsFallingDown is the application reloading? If so, you should not accept
// new connections to your server in this situation
func IsFallingDown() bool {
	return inReloadState.Load()
}

func reloadAppGracefully() {
	if inReloadState.Load() {
		log.Println("warn: [daemon] we are currently restarting the application")
		return
	}

	inReloadState.Store(true)

	FireTrigger(plugins.TRIGGER_RELOAD_APP, nil)
}

func stopAppGracefully(useHook bool) {
	inReloadState.Store(true)

	// we only execute it once, without it we will execute it twice, once when we send the FINISH_APP_NOW command
	// and the next time when we close the channel

	if useHook && hookAfterFinish != nil {
		hookAfterFinish()
	}
	log.Printf("info: [daemon] i'm waiting for the processes to be completed")
	waitingGroup.Wait()
	FireTrigger(plugins.TRIGGER_APPLICATION_EXIT, nil)
}

func doFinishHim() {
	finishHim <- true
}

var inReloadState atomic.Bool

func init() {
	inReloadState.Store(false)
}
