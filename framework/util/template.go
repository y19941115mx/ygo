package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

//  folderNeedCreate 表明路径是否需要创建  folderPath/file 确定了目标文件的路径 tmp 为模板字符串， data为数据
func CreateFileTemlate(folderNeedCreate bool, folderPath string, file string, tmp string, data interface{}) error {

	if folderNeedCreate && Exists(folderPath) {
		fmt.Println("目录已经存在")
		return nil
	}

	if folderNeedCreate {
		if err := os.Mkdir(folderPath, 0700); err != nil {
			return err
		}
	}

	funcs := template.FuncMap{"title": strings.ToTitle}

	//  创建文件
	filePath := filepath.Join(folderPath, file)
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Cause(err)
	}

	t := template.Must(template.New(file).Funcs(funcs).Parse(tmp))
	if err := t.Execute(f, data); err != nil {
		return errors.Cause(err)
	}

	return nil
}
