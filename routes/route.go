/*
 * @PackageName: routes
 * @FileName: route.go
 * @Description: 路由配置
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:41:28
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:41:28
 */

package routes

import (
	c "easy_dfs/app/controllers"
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/middlewares"
	"easy_dfs/pkg/http/http_response"
	"github.com/gin-gonic/gin"
)

// 路由注册
func RegisterRoutes(r *gin.Engine) {
	// 默认路由
	r.GET("/", func(ctx *gin.Context) {
		http_response.Response(ctx, response_code.REQUEST_SUCCESS, true, "操作成功", "欢迎使用本系统", nil)
	})

	// 忽略/favicon.ico
	r.Any("/favicon.ico", func(ctx *gin.Context) {
		ctx.Next()
	})

	// storage路由组，并应用中间件
	sg := r.Group("/storage").Use(middlewares.FileUrlCheck())
	{
		sc := new(c.StorageController)
		sg.Any("/*path", sc.GetFile)
	}

	// 存储桶路由
	br := r.Group("/bucket").Use(middlewares.AccessKeyCheck())
	{
		bc := new(c.BucketController)
		br.POST("/create", bc.CreateBucket)
		br.GET("/list", bc.ListBuckets)
		br.GET("/info", bc.GetBucketInfo)
		br.DELETE("/delete", bc.DeleteBucket)
	}

	// 访问密钥路由
	akr := r.Group("/access_key")
	{
		akc := new(c.AccessKeyController)
		akr.POST("/create", akc.CreateAccessKey)
		akr.GET("/list", akc.ListAccessKeys)
		akr.GET("/info", akc.GetAccessKeyInfo)
		akr.DELETE("/delete", akc.DeleteAccessKey)
	}

	// 文件路由
	fr := r.Group("/file").Use(middlewares.AccessKeyCheck())
	{
		fc := new(c.FileController)
		fr.POST("/upload", fc.UploadFile)
		fr.GET("/download", fc.DownloadFile)
		fr.GET("/list", fc.ListFiles)
		fr.GET("/list-all", fc.ListAllFiles)
		fr.GET("/info", fc.GetFileInfo)
		fr.DELETE("/delete", fc.DeleteFile)
	}
}
