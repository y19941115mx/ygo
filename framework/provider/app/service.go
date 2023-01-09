package app

import (
	"errors"
	"flag"
	"path/filepath"

	"gitee.com/y19941115mx/ygo/framework"
	"gitee.com/y19941115mx/ygo/framework/contract"
	"gitee.com/y19941115mx/ygo/framework/util"
	"github.com/google/uuid"
)

type YgoApp struct {
	// 实现接口
	contract.App

	container  framework.Container // 服务容器
	baseFolder string              // 基础路径
	appId      string              // 表示当前这个app的唯一id, 可用于分布式锁等
	configMap  map[string]string   // 支持配置的方式修改默认设置
}

// Version 实现版本
func (app YgoApp) Version() string {
	return "0.0.1"
}

// BaseFolder 表示基础目录，可以代表开发场景的目录，也可以代表运行时候的目录
func (app YgoApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}

	// 如果没有设置，使用默认的当前路径
	return util.GetExecDirectory()
}

// SourceFolder 定义项目源码路径
func (app YgoApp) SourceFolder() string {
	if val, ok := app.configMap["source_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app")
}

// ConfigFolder  表示配置文件地址
func (app YgoApp) ConfigFolder() string {
	if val, ok := app.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (app YgoApp) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

func (app YgoApp) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.SourceFolder(), "http")
}

func (app YgoApp) ConsoleFolder() string {
	if val, ok := app.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(app.SourceFolder(), "console")
}

func (app YgoApp) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (app YgoApp) ProviderFolder() string {
	if val, ok := app.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(app.SourceFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (app YgoApp) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (app YgoApp) CommandFolder() string {
	if val, ok := app.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (app YgoApp) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (app YgoApp) TestFolder() string {
	if val, ok := app.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "test")
}

// NewYgoApp 初始化App
func NewYgoApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("app param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
		flag.Parse()
	}
	appId := uuid.New().String()
	configMap := map[string]string{}
	return &YgoApp{baseFolder: baseFolder, container: container, appId: appId, configMap: configMap}, nil
}

func (app YgoApp) AppID() string {
	return app.appId
}

func (app *YgoApp) LoadAppConfig(kv map[string]string) {
	app.configMap = kv
}
