package main

import (
	"strings"
	"os"
)

type cmdType uint8

const (
	cmdEmpty cmdType = iota
	cmdInit
	cmdNewNs
	cmdNewCtrl
	cmdRun
	cmdInvalidCmd
)

func getCmdType() (cmdType, []string) {
	args := os.Args
	if len(args) == 1 {
		return cmdEmpty, nil
	}
	if args[1] == "init" {
		return cmdInit, nil
	}
	if args[1] == "ns" {
		return cmdNewNs, args[2:]
	}
	if args[1] == "ctrl" {
		return cmdNewCtrl, args[2:]
	}
	if args[1] == "run" {
		return cmdRun, nil
	}
	return cmdInvalidCmd, []string{args[1]}
}

func showHelp() {
	println("wetool is a tool that can help to create a new wemvc project easily.")
	println("")
	println("Usage:")
	println("    wetool command [arguments]")
	println("")
	println("The commands are:")
	println("    init:           Initialize the project in an empty GO package.")
	println("    ns [name]:      Create a new wemvc namespace with name '[NAME]'.")
	println("    ctrl [name]:    Create a new wemvc controller with name '[NAME]'.")
	println("    run:            Compile and run the wemvc application.")
	return
}

var (
	dir       = WorkingDir()
	goPathSrc = getGoPathSrc()
	pkgPath   string
)

func main() {
	if !strings.HasPrefix(dir, goPathSrc) {
		println("Error: It seems that this is not a GO package folder. This command only can be executed under an empty GO pachage folder.")
		return
	}
	pkgPath = dir[len(goPathSrc)+1:]
	pkgPath = strings.Replace(pkgPath, "\\", "/", -1)
	cmdName, args := getCmdType()
	if cmdName == cmdEmpty {
		showHelp()
		return
	}
	if cmdName == cmdInit {
		initProject()
		return
	}
	if cmdName == cmdNewNs {
		if len(args) != 1 {
			println("wetool: Invalid name of the new namespace.")
			showHelp()
			return
		}
		newNs(args[0])
		return
	}
	if cmdName == cmdNewCtrl {
		if len(args) != 1 {
			println("wetool: Invalid name of the new controller.")
			showHelp()
			return
		}
		newCtrl(args[0])
		return
	}

	if cmdName == cmdRun {
		runProject()
		return
	}
}
