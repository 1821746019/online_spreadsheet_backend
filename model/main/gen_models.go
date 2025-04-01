package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
	})

	dsn := "root:5210165@tcp(127.0.0.1:3306)/MutliTable?charset=utf8mb4&parseTime=True&loc=Local"
	gormdb, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(gormdb)

	// 生成所有表对应的模型
	g.ApplyBasic(
		g.GenerateModel("user"),
		g.GenerateModel("class"),
		g.GenerateModel("sheet"),
		g.GenerateModel("cell"),
		g.GenerateModel("draggable_item"),
		g.GenerateModel("draggable_item_sheet"),
		g.GenerateModel("permission"),
	)

	g.Execute()
}
