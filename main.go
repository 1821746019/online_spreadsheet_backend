package main

import (
	"fmt"
	"log"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/DAO/Redis"
	"github.com/sztu/mutli-table/logger"
	"github.com/sztu/mutli-table/pkg/snowflake"
	"github.com/sztu/mutli-table/router"
	"github.com/sztu/mutli-table/settings"
)

func main() {
	// 设置机器号
	snowflake.SetMachineID(1)

	if err := logger.SetupGlobalLogger(settings.GetConfig().LoggerConfig); err != nil {
		fmt.Printf("初始化日志库失败,错误原因: %v\n", err)
	}

	defer mysql.Close()
	defer Redis.Close()

	// 初始化路由
	r := router.SetupRouter()
	err := r.Run(fmt.Sprintf("%s:%d", settings.GetConfig().Host, settings.GetConfig().Port))
	if err != nil {
		fmt.Printf("启动失败,错误原因: %v\n", err)
	}
	uid, err := snowflake.GetID()
	if err != nil {
		log.Fatal("获取id出错")
	}
	fmt.Print(uid)
}
