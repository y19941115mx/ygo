package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/go-github/v39/github"
	"github.com/spf13/cast"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/util"
)

type AppProject struct {
	name    string
	folder  string
	mod     string
	version string
}

// new相关的名称
func initNewCommand() *cobra.Command {
	return newCommand
}

// 根据用户输入初始化项目的 realease 版本信息
func initAppProjectByAsk(app *AppProject) (*github.RepositoryRelease, error) {
	currentPath := util.GetExecDirectory()
	var release *github.RepositoryRelease

	var version string
	{
		prompt := &survey.Input{
			Message: "请输入项目名称：",
		}
		err := survey.AskOne(prompt, &app.name)
		if err != nil {
			return nil, err
		}

		folder := filepath.Join(currentPath, app.name)
		if util.Exists(folder) {
			return nil, errors.New("目录已存在，创建应用失败")
		}
		app.folder = folder
	}
	{
		prompt := &survey.Input{
			Message: "请输入模块名称(go.mod中的module, 默认为项目名称)：",
		}
		err := survey.AskOne(prompt, &app.mod)
		if err != nil {
			return nil, err
		}
		if app.mod == "" {
			app.mod = app.name
		}
	}
	{
		// 获取ygo的版本
		client := github.NewClient(nil)
		prompt := &survey.Input{
			Message: "请输入版本名称(参考 https://github.com/y19941115mx/ygo/releases，默认为最新版本)：",
		}
		err := survey.AskOne(prompt, &version)
		if err != nil {
			return nil, err
		}
		if version != "" {
			// 确认版本是否正确
			release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "y19941115mx", "ygo", version)
			if err != nil || release == nil {
				fmt.Println("版本不存在，创建应用失败，请参考 https://github.com/y19941115mx/ygo/releases")
				return nil, nil
			}
		}
		if version == "" {
			release, _, err = client.Repositories.GetLatestRelease(context.Background(), "y19941115mx", "ygo")
			if err != nil || release == nil {
				return nil, err
			}
			version = release.GetTagName()
			fmt.Println("使用ygo框架最新版本：" + version)
		}
	}
	app.version = version

	fmt.Println("====================================================")
	fmt.Println("开始进行创建应用操作")
	fmt.Println("应用名称：", app.name)
	fmt.Println("模块名称：", app.mod)
	fmt.Println("创建目录：", app.folder)
	fmt.Println("ygo框架版本：", app.version)

	return release, nil
}

// 创建一个新应用
var newCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个新的应用",
	RunE: func(c *cobra.Command, args []string) error {
		app := AppProject{}
		release, err := initAppProjectByAsk(&app)

		if err != nil {
			return err
		}

		templateFolder := filepath.Join(util.GetExecDirectory(), "template-ygo-"+app.version+"-"+cast.ToString(time.Now().Unix()))
		os.Mkdir(templateFolder, os.ModePerm)
		fmt.Println("创建临时目录", templateFolder)

		// 拷贝template项目
		url := release.GetZipballURL()
		err = util.DownloadFile(filepath.Join(templateFolder, "template.zip"), url)

		if err != nil {
			return err
		}
		fmt.Println("下载zip包到template.zip")

		_, err = util.Unzip(filepath.Join(templateFolder, "template.zip"), templateFolder)
		if err != nil {
			return err
		}

		fInfos, err := os.ReadDir(templateFolder)
		if err != nil {
			return err
		}
		for _, fInfo := range fInfos {
			// 找到解压后的文件夹
			if fInfo.IsDir() && strings.Contains(fInfo.Name(), "ygo") {
				if err := os.Rename(filepath.Join(templateFolder, fInfo.Name()), app.folder); err != nil {
					return err
				}
			}
		}
		fmt.Println("解压zip包")

		if err := os.RemoveAll(templateFolder); err != nil {
			return err
		}
		fmt.Println("删除临时文件夹", templateFolder)

		os.RemoveAll(path.Join(app.folder, ".git"))
		fmt.Println("删除.git目录")

		// 删除framework 目录
		os.RemoveAll(path.Join(app.folder, "framework"))
		fmt.Println("删除framework目录")

		filepath.Walk(app.folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			c, _ := os.ReadFile(path)

			// 修改 go.mod 文件  修改 module 名称  增加对框架的引用
			if path == filepath.Join(app.folder, "go.mod") {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("module github.com/y19941115mx/ygo"), []byte("module "+app.mod))
				c = bytes.ReplaceAll(c, []byte("require ("), []byte("require (\n\tgithub.com/y19941115mx/ygo "+app.version))
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
				return nil
			}
			// 替换 app 文件夹的引用位置
			isContain := bytes.Contains(c, []byte("github.com/y19941115mx/ygo/app"))
			if isContain {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("github.com/y19941115mx/ygo/app"), []byte(app.mod+"/app"))
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
			}

			return nil
		})

		fmt.Println("创建应用结束，目录：", app.folder)
		fmt.Println("====================================================")
		return nil
	},
}
