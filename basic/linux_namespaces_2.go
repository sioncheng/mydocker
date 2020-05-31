package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path"
	"io/ioutil"
	"strconv"
)

const cgroupMemoryHierachyMount = "/sys/fs/cgroup/memory"

func main() {

	fmt.Println("pid ", syscall.Getpid())
	fmt.Println("args", os.Args)

	if "/proc/self/exe" == os.Args[0] {
		fmt.Println("current pid ", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 20M --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("pid after run", cmd.Process.Pid)
	memLimit := path.Join(cgroupMemoryHierachyMount, "testmemorylimit")
	os.Mkdir(memLimit, 0755)
	ioutil.WriteFile(path.Join(memLimit, "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	ioutil.WriteFile(path.Join(memLimit, "memory.limit_in_bytes"), []byte("100m"), 0644)

	cmd.Process.Wait()
	

}