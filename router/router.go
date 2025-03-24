package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sztu/mutli-table/controller"
	"github.com/sztu/mutli-table/logger"
	"github.com/sztu/mutli-table/settings"
	"go.uber.org/zap"
)

// SetupRouter 初始化 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 日志中间件
	r.Use(logger.GinLogger(zap.L()), logger.GinRecovery(zap.L(), true))

	// 根据配置设置 Gin 的模式
	switch settings.GetConfig().Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 创建 API v1 路由组
	v1 := r.Group("/api/v1").Use(
		controller.LimitBodySizeMiddleware(),
		controller.TimeoutMiddleware(),
		controller.CorsMiddleware(
			controller.WithAllowOrigins([]string{"localhost"}),
		),
	)

	// 设置路由
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 用户登录相关路由
	v1.POST("/login", controller.LoginHandler)
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/logout", controller.LogoutHandler)
	v1.Use(controller.JWTAuthMiddleware())
	{
		// // 文档管理
		// v1.POST("/document", controller.CreateDocumentHandler)
		// v1.GET("/documents", controller.ListDocumentsHandler)
		// v1.GET("/document/:id", controller.GetDocumentHandler)
		// v1.PUT("/document/:id", controller.UpdateDocumentHandler)
		// v1.DELETE("/document/:id", controller.DeleteDocumentHandler)
		// v1.POST("/document/:id/share", controller.ShareDocumentHandler)

		// 工作表管理
		// v1.POST("/document/:id/sheet", controller.CreateSheetHandler)
		// v1.GET("/document/:id/sheet", controller.ListSheetsHandler)
		// v1.GET("/document/:id/sheet/:sheet_id", controller.GetSheetHandler)
		// v1.PUT("/document/:id/sheet/:sheet_id", controller.UpdateSheetHandler)
		// v1.DELETE("/document/:id/sheet/:sheet_id", controller.DeleteSheetHandler)
		// 工作表管理
		v1.POST("/sheet", controller.CreateSheetHandler)
		v1.GET("/sheet", controller.ListSheetsHandler)
		v1.GET("/sheet/:sheet_id", controller.GetSheetHandler)
		v1.PUT("/sheet/:sheet_id", controller.UpdateSheetHandler)
		v1.DELETE("/sheet/:sheet_id", controller.DeleteSheetHandler)

		// 单元格管理
		v1.GET("/sheet/:sheet_id/cell", controller.GetCellsHandler)
		v1.PUT("/sheet/:sheet_id/cell", controller.UpdateCellHandler)

		// 新增待拖动单元格管理
		v1.POST("/sheet/:id/drag-item", controller.CreateDragCellHandler)                 // 创建待拖动单元格
		v1.GET("/sheet/:id/drag-item", controller.ListDragCellsHandler)                   // 列出所有待拖动单元格
		v1.GET("/sheet/:id/drag-item/:drag_item_id", controller.GetDragCellHandler)       // 获取单个待拖动单元格
		v1.PUT("/sheet/:id/drag-item/:drag_item_id", controller.UpdateDragCellHandler)    // 更新待拖动单元格
		v1.DELETE("/sheet/:id/drag-item/:drag_item_id", controller.DeleteDragCellHandler) // 删除待拖动单元格

		// 新增拖放操作接口
		v1.PUT("/sheet/:sheet_id/drag-item/:drag_item_id/move", controller.MoveDragItemHandler)

		// // 实时协作（WebSocket）
		// v1.GET("/document/:id/ws", controller.DocumentWebSocketHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "请求的资源不存在",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "请求方式非法",
		})
	})

	return r
}
