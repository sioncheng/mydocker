package main

import (
	"fmt"
	"syscall"
	"os"
)

func main() {
	fmt.Println("current pid", syscall.Getpid())
	self, _ := os.Readlink("/proc/self/exe")
	fmt.Println("readlink", self)
}