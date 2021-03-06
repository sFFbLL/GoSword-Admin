package router

import (

	"net/http"

	admin "project/app/admin/router"
	"project/common/middleware"
	"project/utils"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Setup 路由设置
func Setup(mode string) *gin.Engine {
	if mode == string(utils.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Cors(), middleware.GinLogger(), middleware.GinRecovery(true), middleware.Sentinel(200))
	r.Static("/static", "./static")
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	// 注册业务路由
	admin.InitAdminRouter(r)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	pprof.Register(r) // 注册pprof相关路由
	return r
}
