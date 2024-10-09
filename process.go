package daemon

import (
	"context"
	"fmt"
	"github.com/malumar/daemon/plugins"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	// only we can read, modify and save
	PIDS_STORAGE_PATH_PERMS = 0700
	// only we can read and modify
	PIDS_STORAGE_FILE_PERMS = 0600

	EnvWorkerName      = "DAEMON_WORKER"
	EnvWorkerNameValue = "1"
)

func GetWaitingGroup() *sync.WaitGroup {
	return &waitingGroup
}

var waitingGroup sync.WaitGroup

var myPid int

func init() {
	RegisterListener(plugins.TRIGGER_APPLICATION_EXIT, func(ctx context.Context, eventName string, value interface{}) error {
		removeMyProcessPid()
		return nil
	})
}

func addProcessPid(pid int) {
	myPid = pid
	pidFileName := getPidFileName()
	pidsPath := getPidsPath()
	if len(pidsPath) == 0 {
		log.Printf("err: [daemon] no folder specified for saving pids.")
		return
	}
	log.Printf("info: [daemon] saves the PID %d\n", pid)
	if !IsFileExistsAndIsDir(pidsPath) {
		log.Printf("warn: [daemon] path folder `%s` pid does not exist, I am creating it\n", pidsPath)
		if err := os.Mkdir(pidsPath, PIDS_STORAGE_PATH_PERMS); err != nil {
			log.Printf("err: [daemon] an error occurred while trying to create the pid path folder `%v` %v\n",
				pidsPath, err.Error())
			return
		}
	}

	if err := os.WriteFile(pidFileName, nil, PIDS_STORAGE_FILE_PERMS); err != nil {
		log.Printf("err: [daemon] error %v  occurred while trying to create the pid file folder `%v`",
			err.Error(), pidFileName)
	}
}

func removeMyProcessPid() {
	if myPid == 0 {
		log.Printf("warn: [daemon] PID = 0, it means that we have not initialized the application yet\n")
	}
	fn := getPidFileName()

	if IsFileExistsAndIsFile(fn) {
		if err := os.Remove(fn); err != nil {
			log.Printf("err: [daemon] error %v occurred while trying to delete pid file `%v`", err.Error(), fn)
		}
	} else {
		log.Printf("err: [daemon] while trying to delete pid file `%s` an error occurred: the file does not exist\n", fn)

	}
}

func getPidFileName() string {
	return filepath.Join(getPidsPath(), fmt.Sprintf("%d", myPid))
}

func getPidsPath() string {
	return pidsPathStorage
}

func IsFileExistsAndIsFile(filename string) bool {
	fi, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if fi == nil {
		return false
	}
	if !fi.IsDir() {
		return true
	} else {
		return false
	}
}
func IsFileExistsAndIsDir(filename string) bool {
	fi, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if fi.IsDir() {
		return true
	} else {
		return false
	}
}
