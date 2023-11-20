package routers

import (
	"github.com/gin-gonic/gin"
	"time"
	"toes/global"
	"toes/internal/controllers"
	"toes/internal/middleware"
)

func InstallRouters(g *gin.Engine) error {
	// Web 页面
	g.StaticFile("/", "web/index.html")
	g.Static("/static", "web") // web 静态资源

	// 注册 /health handler.
	g.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := g.Group("/v1")

	if global.Cfg.CheckHeader.All {
		v1.Use(middleware.CheckHeader())
		//v1.Use(middleware.CheckPermission())
	}

	accountV1 := v1.Group("account")
	//accountV1.PUT(":name", accountController.Update)        // 更新
	accountV1.POST("", controllers.AccountCtrl.Create)            // 创建
	accountV1.DELETE(":username", controllers.AccountCtrl.Delete) // 删除
	//accountV1.GET(":name", accountController.Get)           // 获取用户详情
	//accountV1.POST("/list", accountController.List)         // 列表

	sysV1 := v1.Group("sys")
	sysV1.GET("/debug/pprof/", controllers.SystemCtrl.Pprof)
	sysV1.GET("/debug/pprof/:app([\\w]+)", controllers.SystemCtrl.Pprof)
	sysV1.GET("/router/list", controllers.SystemCtrl.RouterList)
	//获取系统信息，用第三方库
	sysV1.GET("/info", controllers.SystemCtrl.SysInfo)
	//sysV1.GET("/version", sysController.Version)
	sysV1.GET("/ws", controllers.SystemCtrl.Ws)

	SetRouters(g)

	return nil
}

func SetRouters(g *gin.Engine) {
	data := make([]map[string]string, 0)
	r := g.Routes()
	for _, v := range r {
		data = append(data, map[string]string{
			"method": v.Method,
			"path":   v.Path,
		})
	}
	global.Cache.Set("CacheRouterKey", data, time.Hour*24)
}
