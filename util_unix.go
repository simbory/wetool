// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package main

import "strings"

func getWorkDirName() string {
	var dir = WorkingDir()
	return dir[strings.LastIndex(dir, "/") + 1:]
}

func getOutputName() string{
	return getWorkDirName() + ".debug.out"
}

func getGoPathSrc() string {
	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 || !IsDir(goPath) {
		panic("Could not find the GOPATH environment variable.")
	}
	return strings.TrimRight(goPath, "/") + "/src"
}

func killProcess(name string) error {
	cmd := exec.Command("taskkill.exe", "/f", "/im", name)
	cmd.Stdout = os.Stdout
	fmt.Println("try to kill process:", name)
	return cmd.Run()
}