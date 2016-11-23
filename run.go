package main

import (
	"os/exec"
	"os"
	"path/filepath"
	"github.com/howeyc/fsnotify"
	"path"
	"strings"
	"sync"
	"time"
)

func killProcess(name string) error {
	defer func() {
		if e := recover(); e != nil {
			timePrintln("kill process recover:", e)
		}
	}()
	if runningCmd != nil && runningCmd.Process != nil {
		timePrintln("trying to kill process:", name)
		return runningCmd.Process.Kill()
	}
	return nil
}

var runningCmd *exec.Cmd
var buildLocker = &sync.RWMutex{}

func runProcess(name string) (*exec.Cmd, error) {
	cmd := exec.Command("./" + name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	timePrintln("trying to run project:", name)
	go cmd.Run()
	return cmd, nil
}

func buildProject() {
	buildLocker.Lock()
	defer buildLocker.Unlock()
	var execName = getOutputName()
	// build the project as bak file (for backup usage)
	os.Remove(execName + ".bak")
	c := exec.Command("go", "build", "-o", execName + ".bak")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		timePrintln("Errors occurred while building the current project")
		return
	}
	// try to kill the running process and remove the executable file
	if IsFile(execName) {
		err := killProcess(execName)
		if err != nil {
			timePrintln("Error:", err.Error())
		}
		timePrintln("trying to remove the old executable file:", execName)
		err = os.Remove(execName)
		if err != nil {
			timePrintln("Error:", err.Error())
		}
	}
	// replace the executable file
	if IsFile(execName) {
		timePrintln("Failed to execute the src: canot delete the old executable file. Maybe it is using by another process.")
		return
	}
	err = os.Rename(execName + ".bak", execName)
	if err != nil {
		timePrintln(err.Error())
		return
	}
	// execute the new file
    cmd,err := runProcess(execName)
	runningCmd = cmd
	if err != nil {
		timePrintln("Error:", err.Error())
	}
}

type srcDetector struct {
}

func (d *srcDetector) CanHandle(path string) bool {
	return strings.HasSuffix(path, ".go")
}

func (d *srcDetector) Handle(ev *fsnotify.FileEvent) {
	strFile := path.Clean(ev.Name)
	if ev.IsDelete() {
		srcWatcher.RemoveWatch(strFile)
	} else if ev.IsCreate() {
		srcWatcher.AddWatch(strFile)
	}
	buildProject()
}

var srcWatcher *FileWatcher

func runProject() {
	w,err := NewWatcher()
	if err == nil {
		srcWatcher = w
	} else {
		panic(err)
	}
	var workingDir = WorkingDir()
	srcWatcher.AddWatch(workingDir)
	filepath.Walk(workingDir, func(p string, info os.FileInfo, er error) error {
		if info.IsDir() {
			srcWatcher.AddWatch(p)
		}
		return nil
	})
	srcWatcher.AddHandler(&srcDetector{})
	srcWatcher.Start()
	buildProject()
	for {
		time.Sleep(10*time.Second)
	}
}