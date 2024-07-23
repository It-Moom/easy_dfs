/*
 * @PackageName: controllers
 * @FileName: storage_controller.go
 * @Description: 存储控制器
 * @Author: gabbymrh
 * @Date: 2024-07-20 12:09:41
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-20 12:09:41
 */

package controllers

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/services"
	"easy_dfs/pkg/http/http_response"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"io"
	"net/url"
	"path/filepath"
	"strings"
)

// 存储控制器
type StorageController struct {
	FileService services.FileService
}

// // 获取文件
// func (sc *StorageController) GetFile(c *gin.Context) {
// 	// 解析查询参数
// 	queryParams, bcerr := url.ParseQuery(c.Request.URL.RawQuery)
// 	if bcerr != nil {
// 		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, bcerr)
// 		return
// 	}
//
// 	// 获取bucket
// 	bucket := queryParams.Get("bucket")
// 	// 获取URL路径
// 	path := c.Request.URL.Path
// 	// 移除 /storage/{bucket} 前缀
// 	path = path[len("/storage/"+bucket):]
//
// 	// 加载文件
// 	reader, err := sc.FileService.LoadFileByPath(bucket, path)
// 	if err != nil {
// 		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
// 		return
// 	}
//
// 	// 读取文件内容来检测MIME类型
// 	mtype, err := mimetype.DetectReader(reader)
// 	if err != nil {
// 		http_response.Response(c, response_code.REQUEST_FAILS, false, "无法检测文件类型", nil, err)
// 		return
// 	}
//
// 	// 重置reader的读取位置
// 	if seeker, ok := reader.(io.Seeker); ok {
// 		seeker.Seek(0, io.SeekStart)
// 	}
//
// 	// 判断MIME类型，设置Content-Type和Content-Disposition
// 	contentType := mtype.String()
// 	if strings.HasPrefix(contentType, "image/") || contentType == "application/pdf" {
// 		c.Header("Content-Disposition", "inline; filename="+filepath.Base(path))
// 	} else {
// 		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(path))
// 	}
// 	c.Header("Content-Type", contentType)
//
// 	// 将文件内容写入响应
// 	io.Copy(c.Writer, reader)
// }

func (sc *StorageController) GetFile(c *gin.Context) {
	// 解析查询参数
	queryParams, bcerr := url.ParseQuery(c.Request.URL.RawQuery)
	if bcerr != nil {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, bcerr)
		return
	}

	// 获取bucket
	bucket := queryParams.Get("bucket")
	if bucket == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "缺少bucket参数", nil, nil)
		return
	}

	// 获取URL路径
	path := c.Request.URL.Path

	// 移除 /storage/{bucket} 前缀
	if !strings.HasPrefix(path, "/storage/"+bucket) {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "路径不匹配", nil, nil)
		return
	}
	path = path[len("/storage/"+bucket):]

	// 确保路径安全
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "路径不能是绝对路径", nil, nil)
		return
	}

	// 加载文件
	reader, err := sc.FileService.LoadFileByPath(bucket, path)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}

	// 读取文件内容来检测MIME类型
	mtype, err := mimetype.DetectReader(reader)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "无法检测文件类型", nil, err)
		return
	}

	// 重置reader的读取位置
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// 判断MIME类型，设置Content-Type和Content-Disposition
	contentType := mtype.String()
	if strings.HasPrefix(contentType, "image/") || contentType == "application/pdf" {
		c.Header("Content-Disposition", "inline; filename="+filepath.Base(path))
	} else {
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(path))
	}
	c.Header("Content-Type", contentType)

	// 将文件内容写入响应
	io.Copy(c.Writer, reader)
}
