package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/contract"
	"gopkg.in/yaml.v2"
)

type YgoConfig struct {
	c        framework.Container    // 容器
	folder   string                 // 文件夹
	keyDelim string                 // 路径的分隔符，默认为点
	lock     sync.RWMutex           // 配置文件读写锁
	envMaps  map[string]string      // 所有的环境变量
	confMaps map[string]interface{} // 配置文件结构，key为文件名
	confRaws map[string][]byte      // 配置文件的原始信息, key为文件名
}

// 读取某个配置文件
func (conf *YgoConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	//  判断文件是否以yaml或者yml作为后缀
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		// 读取文件内容
		bf, err := os.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		// 直接针对文本做环境变量的替换
		bf = replace(bf, conf.envMaps)
		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		// 读取app.path中的信息，更新app对应的folder
		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := conf.c.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

// 删除文件的操作
func (conf *YgoConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func NewYgoConfig(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("config param error")
	}
	container := params[0].(framework.Container)

	appService := container.MustMake(contract.AppKey).(contract.App)
	envService := container.MustMake(contract.EnvKey).(contract.Env)

	// 配置文件夹地址
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	// 实例化
	ygoConf := &YgoConfig{
		c:        container,
		folder:   envFolder,
		envMaps:  envService.All(),
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		keyDelim: ".",
		lock:     sync.RWMutex{},
	}
	// 检查配置文件夹是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return ygoConf, nil
	}

	// 读取每个文件
	files, err := os.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, file := range files {
		fileName := file.Name()
		err := ygoConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					//判断事件发生的类型
					// Create 创建
					// Write 写入
					// Remove 删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建配置文件 : ", ev.Name)
						ygoConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("修改配置文件 : ", ev.Name)
						ygoConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除配置文件 : ", ev.Name)
						ygoConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	return ygoConf, nil
}

func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch next := next.(type) {
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return searchMap(next, path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

// 通过path来获取某个配置项
func (conf *YgoConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

// IsExist check setting is exist
func (conf *YgoConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

// Get 获取某个配置项
func (conf *YgoConfig) Get(key string) interface{} {
	return conf.find(key)
}

// GetBool 获取bool类型配置
func (conf *YgoConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

// GetInt 获取int类型配置
func (conf *YgoConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *YgoConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *YgoConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *YgoConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *YgoConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *YgoConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *YgoConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *YgoConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *YgoConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be an pointer
func (conf *YgoConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(conf.find(key))
}
