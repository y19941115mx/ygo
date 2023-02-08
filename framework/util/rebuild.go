package util

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func RebuildApp() error {
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("ygo go: 请在Path路径中先安装go")
	}
	var binFile string
	sysType := runtime.GOOS

	if sysType == "linux" {
		// LINUX系统
		binFile = "ygo"
	}

	if sysType == "windows" {
		// windows系统
		binFile = "ygo.exe"
	}

	buildCmd := exec.Command(path, "build", "-o", binFile)
	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Println("go build error:")
		fmt.Println(string(out))
		fmt.Println("--------------")
		return err
	}
	fmt.Println("本地编译成功")
	return nil
}
