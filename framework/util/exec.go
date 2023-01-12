package util

import (
	"os"

	"github.com/shirou/gopsutil/process"
)

// GetExecDirectory 获取当前执行程序目录
func GetExecDirectory() string {
	file, err := os.Getwd()
	if err == nil {
		return file + "/"
	}
	return ""
}

// CheckProcessExist Will return true if the process with PID exists.
// func CheckProcessExist(pid int) bool {
// 	process, err := os.FindProcess(pid)
// 	if err != nil {
// 		return false
// 	}

// 	err = process.Signal(syscall.Signal(0))
// 	return err == nil
// }

func CheckProcessExist(pid int) bool {
	_, err := process.NewProcess(int32(pid))
	return err == nil

}

func KillPid(pid int) error {
	p, err := process.NewProcess(int32(pid))
	// 杀死进程
	if err != nil {
		return err
	}

	if err := p.Kill(); err != nil {
		return err
	}
	return nil
}
