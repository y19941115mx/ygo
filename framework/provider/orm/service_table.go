package orm

import (
	"context"

	"github.com/jianfengye/collection"
	"github.com/y19941115mx/ygo/framework/contract"
	"gorm.io/gorm"
)

func (app *YgoGorm) GetTables(ctx context.Context, db *gorm.DB) ([]string, error) {
	return db.Migrator().GetTables()
}

func (app *YgoGorm) HasTable(ctx context.Context, db *gorm.DB, table string) (bool, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return false, err
	}
	tableColl := collection.NewStrCollection(tables)
	isContain := tableColl.Contains(table)
	return isContain, nil
}

func (app *YgoGorm) GetTableColumns(ctx context.Context, db *gorm.DB, table string) ([]contract.TableColumn, error) {
	// 执行原始的SQL语句
	var columns []contract.TableColumn
	result := db.Raw("SHOW COLUMNS FROM " + table).Scan(&columns)
	if result.Error != nil {
		// 处理错误
		return nil, result.Error
	}
	return columns, nil
}
