package model

import (
	"github.com/pkg/errors"
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jianfengye/collection"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/provider/orm"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func selectForTableName(db *gorm.DB) ([]string, error) {
	// 获取所有表
	dbTables, err := db.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	tables := make([]string, 0, len(dbTables)+1)
	tables = append(tables, "*")
	tables = append(tables, dbTables...)

	genTables := make([]string, 0, len(dbTables)+1)
	{
		// 展示数据库中的表格列表，供用户选择：
		prompt := &survey.MultiSelect{
			Message: "请选择要生成模型的表格：",
			Options: tables,
		}
		survey.AskOne(prompt, &genTables)

		if collection.NewStrCollection(genTables).Contains("*") {
			genTables = dbTables
		}
	}

	if len(genTables) == 0 {
		return nil, errors.New("未选择任何需要生成模型的数据表")
	}
	return genTables, nil
}

// modelGenCommand 生成数据库模型的model代码文件
var modelGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成模型",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		logger := container.MustMake(contract.LogKey).(contract.Log)
		logger.SetLevel(contract.ErrorLevel)
		// 创建数据库gorm服务
		gormService := container.MustMake(contract.ORMKey).(contract.ORMService)
		db, err := gormService.GetDB(orm.WithConfigPath(database))
		if err != nil {
			fmt.Println("数据库连接：" + database + "失败，请检查配置")
			return err
		}

		// 创建交互命令，列表选择
		genTables, err := selectForTableName(db)
		if err != nil {
			return err
		}

		// 第二步确认要生成的目录文件位置，默认为 app/model文件夹：
		if !filepath.IsAbs(output) {
			absOutput, err := filepath.Abs(output)
			if err != nil {
				return err
			}
			output = absOutput
		}

		genFileNames := make([]string, 0, len(genTables))
		for _, genTable := range genTables {
			genFileNames = append(genFileNames, genTable+".gen.go")
		}

		fmt.Println("将生成下列文件：")
		for _, genFileName := range genFileNames {
			fmt.Println(genFileName + "（新文件）")
		}

		// 第三步选择后是一个生成模型的选项：
		selectRuleTips := []string{}
		ruleTips := map[string]string{
			"FieldNullable":     "FieldNullable, 对于数据库的可null字段设置指针",
			"FieldCoverable":    "FieldCoverable, 根据数据库的Default设置字段的默认值",
			"FieldWithIndexTag": "FieldWithIndexTag, 根据数据库的索引关系设置索引标签",
			"FieldWithTypeTag":  "FieldWithTypeTag, 根据数据库字段类型生成type类型字段",
		}
		tips := make([]string, 0, len(ruleTips))
		for _, val := range ruleTips {
			tips = append(tips, val)
		}
		promptRules := &survey.MultiSelect{
			Message: "请选择生成的模型规则：",
			Options: tips,
		}
		survey.AskOne(promptRules, &selectRuleTips)
		isSelectRule := func(key string, selectRuleTips []string, allRuleTips map[string]string) bool {
			tip := allRuleTips[key]
			selectRuleTipsColl := collection.NewStrCollection(selectRuleTips)
			return selectRuleTipsColl.Contains(tip)
		}

		// 生成模型文件
		g := gen.NewGenerator(gen.Config{
			OutPath: output,

			FieldNullable:     isSelectRule("FieldNullable", selectRuleTips, ruleTips),
			FieldCoverable:    isSelectRule("FieldCoverable", selectRuleTips, ruleTips),
			FieldWithIndexTag: isSelectRule("FieldWithIndexTag", selectRuleTips, ruleTips),
			FieldWithTypeTag:  isSelectRule("FieldWithTypeTag", selectRuleTips, ruleTips),

			Mode: gen.WithDefaultQuery,
		})

		g.UseDB(db)

		for _, table := range genTables {
			g.GenerateModel(table)
		}

		g.Execute()

		fmt.Println("生成模型成功")
		return nil
	},
}
