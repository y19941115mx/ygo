package command

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
	"github.com/shirou/gopsutil/process"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/util"
)

// app启动地址
var appPort = ""
var appDaemon = false

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVarP(&appPort, "port", "p", "", "设置app启动的端口号默认为8888")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appStopCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务",
	RunE: func(c *cobra.Command, args []string) error {
		// 从Command中获取服务容器
		container := c.GetContainer()
		// 从服务容器中获取kernel的服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		// 从kernel服务实例中获取引擎
		core := kernelService.HttpEngine()

		// 如果启动命令时没有设置端口号port 首先读取环境变量YGO_ADDRESS 再读取配置文件 app.address 默认端口 8888
		if appPort == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("YGO_ADDRESS") != "" {
				appPort = envService.Get("YGO_ADDRESS")
			} else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.address") {
					appPort = configService.GetString("app.address")
				} else {
					appPort = "8888"
				}
			}
		}

		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    ":" + appPort,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return err
			}
		}
		// 应用日志
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()

		// daemon 模式
		if appDaemon {
			// 创建一个Context
			cntxt := &daemon.Context{
				// 设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				// 设置日志文件
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				// 设置工作路径
				WorkDir: currentFolder,
				// 设置所有设置文件的mask，默认为750
				Umask: 027,
				// 子进程的参数，按照这个参数设置，子进程的命令为 ./hade app start --daemon=true
				Args: []string{"", "app", "start", "--daemon=true"},
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				// 父进程直接打印启动成功信息，不做任何操作
				fmt.Println("app启动成功，pid:", d.Pid)
				fmt.Println("日志文件:", serverLogFile)
				return nil
			}
			defer cntxt.Release()
			// 子进程执行真正的app启动操作
			fmt.Println("deamon started")
			// gspt.SetProcTitle("ygo app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}
		// 非deamon模式，直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := os.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}
		// gspt.SetProcTitle("ygo app")

		fmt.Println("app serve port: ", appPort)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

// 重新启动一个app服务
var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重新启动一个app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		// 不存在pid文件 直接daemon方式启动apps
		if !util.Exists(serverPidFile) {
			appDaemon = true
			return appStartCommand.RunE(c, args)
		}

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				// 杀死进程
				if err := killPid(pid); err != nil {
					return err
				}
				// 获取closeWait
				closeWait := 5
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}

				// 确认进程已经关闭,每秒检测一次， 最多检测closeWait * 2秒
				for i := 0; i < closeWait*2; i++ {
					if !util.CheckProcessExist(pid) {
						break
					}
					time.Sleep(1 * time.Second)
				}

				// 如果进程等待了2*closeWait之后还没结束，返回错误，不进行后续的操作
				if util.CheckProcessExist(pid) {
					fmt.Println("结束进程失败:"+strconv.Itoa(pid), "请查看原因")
					return errors.New("结束进程失败")
				}
				if err := os.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}

				fmt.Println("结束进程成功:" + strconv.Itoa(pid))
			}
		}

		appDaemon = true
		// 直接daemon方式启动apps
		return appStartCommand.RunE(c, args)
	},
}

// 停止一个已经启动的app服务
var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			// 杀死进程
			if err := killPid(pid); err != nil {
				return err
			}
			if err := os.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("停止进程:", pid)
		}
		return nil
	},
}

// 获取启动的app的pid
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app服务已经启动, pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}

// 启动AppServer, 这个函数会将当前goroutine阻塞
func startAppServe(server *http.Server, c framework.Container) error {
	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal, 1)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	closeWait := 5
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	return nil
}

func killPid(pid int) error {
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
