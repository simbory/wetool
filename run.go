package main

import (
	"os/exec"
	"os"
	"fmt"
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
			fmt.Println("kill process recover:", e)
		}
	}()
	if currentRunning != nil && currentRunning.Process != nil {
		fmt.Println("try to kill process:", name)
		return currentRunning.Process.Kill()
	}
	return nil
}

var currentRunning *exec.Cmd
var lock = &sync.RWMutex{}

func runProcess(name string) (*exec.Cmd, error) {
	cmd := exec.Command("./" + name)
	cmd.Stdout = os.Stdout
	fmt.Println("try to run project:", name)
	go cmd.Run()
	return cmd, nil
}

func buildProject() {
	lock.Lock()
	var execName = getOutputName()
	if IsFile(execName) {
		err := killProcess(execName)
		if err != nil {
			println("Error:", err.Error())
		}
		err = os.Remove(execName)
		if err != nil {
			println("Error:", err.Error())
		}
	}
	c := exec.Command("go", "build", "-o", getOutputName())
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		println("Errors occurred while building current project.")
	}
    cmd,err := runProcess(execName)
	currentRunning = cmd
	if err != nil {
		println("Error:", err.Error())
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