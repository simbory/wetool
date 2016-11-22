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
	if currentRunning != nil && currentRunning.Process != nil {
		timePrintln("try to kill process:", name)
		return currentRunning.Process.Kill()
	}
	return nil
}

var currentRunning *exec.Cmd
var lock = &sync.RWMutex{}

func runProcess(name string) (*exec.Cmd, error) {
	cmd := exec.Command("./" + name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	timePrintln("trying to run project:", name)
	go cmd.Run()
	return cmd, nil
}

func buildProject() {
	lock.Lock()
	var execName = getOutputName()
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
	c := exec.Command("go", "build", "-o", getOutputName())
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		timePrintln("Errors occurred while building current project.")
	}
    cmd,err := runProcess(execName)
	currentRunning = cmd
	if err != nil {
		timePrintln("Error:", err.Error())
	}
	lock.Unlock()
}

type srcDetector struct {
}

func (d *srcDetector) CanHandle(path string) bool {
	return strings.HasSuffix(path, ".go")
}

func (d *srcDetector) Handle(ev *fsnotify.FileEvent) {
	strFile := path.Clean(ev.Name)
	if IsDir(strFile) {
		if ev.IsDelete() {
			srcWatcher.RemoveWatch(strFile)
		} else if ev.IsCreate() {
			srcWatcher.AddWatch(strFile)
		}
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