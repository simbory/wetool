package main

import (
	"os"
	"time"
	"fmt"
	"strings"
)

// IsDir check if the path is directory
func IsDir(path string) bool {
	state, err := os.Stat(path)
	if err != nil {
		return false
	}
	return state.IsDir()
}

// IsFile check if the path is file
func IsFile(path string) bool {
	state, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !state.IsDir()
}

// WorkingDir get the current working directory
func WorkingDir() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Replace(p, "\\", "/", -1)
}

func timePrintln(args ...interface{}) {
	var t = time.Now()
	args = append([]interface{}{t.Format(time.RFC850), ":"}, args...)
	fmt.Println(args...)
}