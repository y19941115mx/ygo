package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/util"
)

var refreshTime int
var devPort int
var backendPid = 0

// 初始化Dev命令
func initDevCommand() *cobra.Command {
	devCommand.Flags().IntVarP(&refreshTime, "time", "t", 1, "如果发生文件变更，等待设置的时间再进行更新, 默认1s")
	devCommand.Flags().IntVarP(&devPort, "port", "p", 8888, "启动调试模式的服务端口号，默认为8888")
	return devCommand
}

// devCommand 为调试模式的一级命令
var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "启动调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		go monitorBackend(c.GetContainer())
		// 启动服务
		if err := restartService(); err != nil {
			return err
		}
		select {}
	},
}

// monitorBackend 监听应用文件
func monitorBackend(container framework.Container) error {
	// 监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 开启监听app文件夹
	appService := container.MustMake(contract.AppKey).(contract.App)
	appFolder := appService.AppFolder()
	fmt.Println("监控源码文件夹：", appFolder)
	// 监听所有子目录，需要使用filepath.walk
	filepath.Walk(appFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		// 如果是隐藏的目录比如 . 或者 .. 则不用进行监控
		if util.IsHiddenDirectory(path) {
			return nil
		}
		return watcher.Add(path)
	})

	// 开启计时时间机制
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	// 先停止计时器
	t.Stop()
	for {
		select {
		case <-t.C:
			// 计时器时间到了，代表之前有文件更新事件重置过计时器
			// 即有文件更新
			fmt.Println("...检测到文件更新，重启服务开始...")
			if err := util.RebuildApp(); err != nil {
				fmt.Println("重新编译失败：", err.Error())
			} else {
				if err := restartService(); err != nil {
					fmt.Println("重新启动失败：", err.Error())
				}
			}
			fmt.Println("OK 重启服务成功...")
			// 停止计时器
			t.Stop()
		case _, ok := <-watcher.Events:
			if !ok {
				continue
			}
			// 有文件更新事件，重置计时器
			t.Reset(time.Duration(refreshTime) * time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			// 如果有文件监听错误，则停止计时器
			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
}

// restartService 启动后端服务
func restartService() error {
	// 杀死之前的进程
	if backendPid != 0 {
		util.KillPid(backendPid)
		backendPid = 0
	}
	// 使用命令行启动后端进程 指定端口
	cmd := exec.Command("./ygo", "app", "start", "--port=", fmt.Sprint(devPort))
	// cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("启动调试服务: ", "http://127.0.0.1:", devPort)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	backendPid = cmd.Process.Pid
	fmt.Println("服务pid:", backendPid)
	return nil
}
