package model

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/provider/orm"
	"github.com/y19941115mx/ygo/framework/util"
)

func mountRouterFunc(filePath string, funcStr string, importstr string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	var i int

	for i = len(lines) - 1; i >= 0; i-- {
		if lines[i] == "}" {
			lines = append(lines[:i], append([]string{"\t" + funcStr}, lines[i:]...)...)
			break
		}
	}

	for ; i >= 0; i-- {
		if lines[i] == ")" {
			lines = append(lines[:i], append([]string{"\t" + importstr}, lines[i:]...)...)
			break
		}
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	writer.Flush()
	return nil
}

var modelApiCommand = &cobra.Command{
	Use:          "api",
	Short:        "通过数据库生成api",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		container := c.GetContainer()
		logger := container.MustMake(contract.LogKey).(contract.Log)
		logger.SetLevel(contract.ErrorLevel)

		gormService := container.MustMake(contract.ORMKey).(contract.ORMService)
		db, err := gormService.GetDB(orm.WithConfigPath(database))
		if err != nil {
			fmt.Println("数据库连接：" + database + "失败，请检查配置")
			return err
		}

		// 检测数据库是否可以连接
		canConnected, err := gormService.CanConnect(ctx, db)
		if err != nil || !canConnected {
			fmt.Println("数据库连接：" + database + "失败，请检查配置")
			return err
		}

		tables, err := gormService.GetTables(ctx, db)
		if err != nil {
			return errors.Wrap(err, "获取数据库表格失败")
		}

		table := ""
		{
			// 第一步是一个交互命令行工具，首先展示要生成的表列表选择：
			prompt := &survey.Select{
				Message: "请选择要生成模型的表格：",
				Options: tables,
			}
			survey.AskOne(prompt, &table)
		}

		hasTable, err := gormService.HasTable(ctx, db, table)
		if err != nil {
			return fmt.Errorf("数据库连接失败，表格 %v, 错误 %v", table, err)
		}

		if !hasTable {
			return fmt.Errorf("表格 %v 不存在", table)
		}

		// 获取所有字段
		columns, err := gormService.GetTableColumns(ctx, db, table)
		if err != nil {
			return fmt.Errorf("获取表格 %v 列表字段失败: %v", table, err)
		}

		tableLower := strings.ToLower(table)

		// 生成接口代码
		app := container.MustMake(contract.AppKey).(contract.App)
		folder := filepath.Join(app.HttpModuleFolder(), output)

		modelFile := filepath.Join(folder, fmt.Sprintf("gen_%s_model.go", tableLower))
		routerFile := filepath.Join(folder, fmt.Sprintf("gen_%s_router.go", tableLower))
		apiCreateFile := filepath.Join(folder, fmt.Sprintf("gen_%s_api_create.go", tableLower))
		apiDeleteFile := filepath.Join(folder, fmt.Sprintf("gen_%s_api_delete.go", tableLower))
		apiListFile := filepath.Join(folder, fmt.Sprintf("gen_%s_api_list.go", tableLower))
		apiShowFile := filepath.Join(folder, fmt.Sprintf("gen_%s_api_show.go", tableLower))
		apiUpdateFile := filepath.Join(folder, fmt.Sprintf("gen_%s_api_update.go", tableLower))

		// 检测会重新生成如下文件
		{
			getFileTip := func(file string) string {
				if util.Exists(file) {
					return "[替换] " + file
				}
				return "[生成] " + file
			}
			fmt.Println("继续命令会做如下操作：")
			fmt.Println(getFileTip(modelFile))
			fmt.Println(getFileTip(routerFile))
			fmt.Println(getFileTip(apiCreateFile))
			fmt.Println(getFileTip(apiDeleteFile))
			fmt.Println(getFileTip(apiListFile))
			fmt.Println(getFileTip(apiShowFile))
			fmt.Println(getFileTip(apiUpdateFile))

			isContinue := false
			// 会生成的文件列表
			prompt := &survey.Confirm{
				Message: "是否继续操作",
			}
			_ = survey.AskOne(prompt, &isContinue)

			if !isContinue {
				fmt.Println("暂停操作")
				return nil
			}
		}

		apiGenerator := NewApiGenerator(table, columns)
		// get folder last string split by path separator
		apiGenerator.SetPackageName(strings.ToLower(filepath.Base(folder)))

		if !util.Exists(folder) {
			if err := os.Mkdir(folder, 0755); err != nil {
				return errors.Wrap(err, "create folder error")
			}
		}

		if err := apiGenerator.GenModelFile(ctx, modelFile); err != nil {
			return errors.Wrap(err, "GenModelFile error")
		}
		if err := apiGenerator.GenRouterFile(ctx, routerFile); err != nil {
			return errors.Wrap(err, "GenRouterFile error")
		}
		if err := apiGenerator.GenApiCreateFile(ctx, apiCreateFile); err != nil {
			return errors.Wrap(err, "GenApiCreateFile error")
		}
		if err := apiGenerator.GenApiDeleteFile(ctx, apiDeleteFile); err != nil {
			return errors.Wrap(err, "GenApiDeleteFile error")
		}
		if err := apiGenerator.GenApiListFile(ctx, apiListFile); err != nil {
			return errors.Wrap(err, "GenApiListFile error")
		}
		if err := apiGenerator.GenApiShowFile(ctx, apiShowFile); err != nil {
			return errors.Wrap(err, "GenApiShowFile error")
		}
		if err := apiGenerator.GenApiUpdateFile(ctx, apiUpdateFile); err != nil {
			return errors.Wrap(err, "GenApiUpdateFile error")
		}

		// 检测会重新生成如下文件
		{
			getFileTip := func(file string) string {
				return "[成功] " + file
			}
			fmt.Println(getFileTip(modelFile))
			fmt.Println(getFileTip(routerFile))
			fmt.Println(getFileTip(apiCreateFile))
			fmt.Println(getFileTip(apiDeleteFile))
			fmt.Println(getFileTip(apiListFile))
			fmt.Println(getFileTip(apiShowFile))
			fmt.Println(getFileTip(apiUpdateFile))
		}

		fmt.Println("=======================")
		tableTile := util.ToTitleCamel(tableLower)
		routerFuncStr := output + "." + tableTile + "ApiRegister(r)"
		importStr := "\"github.com/y19941115mx/ygo/app/http/module/" + output + "\""
		rootRouterFile := filepath.Join(app.HttpFolder(), "route.go")

		err = mountRouterFunc(rootRouterFile, routerFuncStr, importStr)
		if err != nil {
			return err
		}
		fmt.Println("路由挂载成功")
		fmt.Println("代码生成结束，请自行检查代码逻辑，如需生成接口文档请执行swagger gen命令")
		fmt.Println("=======================")
		return nil
	},
}
