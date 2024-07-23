/*
 * @PackageName: bootstrap
 * @FileName: route.go
 * @Description: 路由配置
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:44:04
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:44:04
 */

package bootstrap

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/middlewares"
	"easy_dfs/pkg/config"
	"easy_dfs/pkg/http/http_response"
	"easy_dfs/routes"
	"github.com/gin-gonic/gin"
)

// 引导安装路由
func SetupRoute() {
	// GIN设为生产模式,不输出日志
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 找不到路由时的返回
	router.NoRoute(func(ctx *gin.Context) {
		http_response.Response(ctx, response_code.QUERY_EMPTY, false, "资源不存在", nil, nil)
	})
	// 路由注册
	routes.RegisterRoutes(router)

	// 全局中间件
	// 跨域处理中间件,文件URL检查中间件
	router.Use(
		middlewares.Cors(),
		middlewares.FileUrlCheck(),
		middlewares.Logger(),
		middlewares.Recovery(),
	)
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		panic(err)
	}
}
