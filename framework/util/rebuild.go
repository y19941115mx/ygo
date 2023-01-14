package util

import (
	"fmt"
	"log"
	"os/exec"
)

func RebuildApp() error {
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("hade go: 请在Path路径中先安装go")
	}

	// modCmd := exec.Command(path, "mod", "tidy", "-compat=1.17")
	// if out, err := modCmd.CombinedOutput(); err != nil {
	// 	fmt.Println("go mod tidy error:")
	// 	fmt.Println(string(out))
	// 	fmt.Println("--------------")
	// 	return err
	// }

	buildCmd := exec.Command(path, "build", "-o", "ygo")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Println("go build error:")
		fmt.Println(string(out))
		fmt.Println("--------------")
		return err
	}
	fmt.Println("重新编译成功")
	return nil
}
