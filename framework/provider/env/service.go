package env

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
)

// YgoEnv 是 Env 的具体实现
type YgoEnv struct {
	contract.Env
	maps   map[string]string // 保存所有的环境变量
	IsTest bool              // 代表是否为本地测试环境
}

// NewYgoEnv 有一个参数，代表当前是否为测试环境
func NewYgoEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("NewYgoEnv param error")
	}

	container := params[0].(framework.Container)
	isTest := params[1].(bool)
	app := container.MustMake(contract.AppKey).(contract.App)
	// 实例化
	ygoEnv := &YgoEnv{
		IsTest: isTest,
		// 实例化环境变量，APP_ENV默认设置为开发环境
		maps: map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// 解析folder/.env文件
	file := path.Join(app.BaseFolder(), ".env")
	// 读取.env文件, 不管任意失败，都不影响后续

	// 打开文件.env
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		// 读取文件
		br := bufio.NewReader(fi)
		for {
			// 按照行进行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			// 按照等号解析
			s := bytes.SplitN(line, []byte{'='}, 2)
			// 如果不符合规范，则过滤
			if len(s) < 2 {
				continue
			}

			key := string(bytes.TrimSpace(s[0]))
			val := string(bytes.TrimSpace(s[1]))
			ygoEnv.maps[key] = val
		}
	}

	// 获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		ygoEnv.maps[pair[0]] = pair[1]
	}

	// 返回实例
	return ygoEnv, nil
}

// AppEnv 获取表示当前APP环境的变量APP_ENV
func (en *YgoEnv) AppEnv() string {
	if en.IsTest {
		return "testing"
	}
	return en.Get("APP_ENV")
}

// IsExist 判断一个环境变量是否有被设置
func (en *YgoEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// Get 获取某个环境变量，如果没有设置，返回""
func (en *YgoEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

// All 获取所有的环境变量，.env和运行环境变量融合后结果
func (en *YgoEnv) All() map[string]string {
	return en.maps
}
