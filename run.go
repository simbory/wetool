package main

import (
	"os/exec"
	"os"
	"fmt"
)

func killProcess(name string) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("kill process recover:", e)
		}
	}()
	fmt.Println("try to kill process:", name)
	if currentRunning != nil && currentRunning.Process != nil {
		return currentRunning.Process.Kill()
	}
	return nil
}

var currentRunning *exec.Cmd

func runProcess(name string) (*exec.Cmd, error) {
	cmd := exec.Command("./" + name)
	cmd.Stdout = os.Stdout
	fmt.Println("try to run project:", name)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func buildProject() {
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
		println("Error while building current project.")
	}
    cmd,err := runProcess(execName)
	currentRunning = cmd
	if err != nil {
		println("Error:", err.Error())
	}
}

func runProject() {
	buildProject()
}