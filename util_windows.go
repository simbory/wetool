package main

import (
	"strings"
	"os"
)

func getWorkDirName() string {
	var dir = WorkingDir()
	return dir[strings.LastIndex(dir, "/") + 1:]
}

func getOutputName() string{
	return getWorkDirName() + ".dev.exe"
}

func getGoPathSrc() string {
	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 || !IsDir(goPath) {
		panic("Could not find the GOPATH environment variable.")
	}
	return strings.Replace(strings.TrimRight(goPath, "\\") + "\\src", "\\", "/", -1)
}