/*
 * @PackageName: controllers
 * @FileName: file_controller.go
 * @Description: 文件控制器
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:18:22
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:18:22
 */
package controllers

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/enum/system_default"
	"easy_dfs/app/services"
	"easy_dfs/pkg/config"
	"easy_dfs/pkg/http/http_response"
	"easy_dfs/pkg/utils/str_util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// 文件控制器
type FileController struct {
	FileService services.FileService
}

// 文件上传返回数据结构体
type UploadResponse struct {
	Bucket       string `json:"bucket"`
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
	FileUrl      string `json:"fileUrl"`
	FileExt      string `json:"fileExt"`
	FileSize     int64  `json:"fileSize"`
}

// 上传文件
func (fc *FileController) UploadFile(c *gin.Context) {
	bucket := c.PostForm("bucket")
	if bucket == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("bucket不能为空"))
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, err)
		return
	}
	defer file.Close()

	// 文件后缀
	ext := filepath.Ext(header.Filename)

	// 保存路径及文件名
	savePath := c.PostForm("savePath")
	if savePath == "" {
		savePath = ""
	}
	filename := savePath + str_util.SimpleUUID() + ext
	// 保存文件名
	saveName := c.PostForm("saveName")
	if saveName != "" {
		filename = savePath + saveName + ext
	}
	if strings.Contains(filename, system_default.STORAGE_PATH) {
		http_response.Response(c, response_code.REQUEST_DENIED, false, "操作失败", nil, errors.New("文件保存名称非法"))
		return
	}

	if err := fc.FileService.SaveFile(bucket, filename, file); err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}

	fileURL := fmt.Sprintf("%s/storage/%s/%s?bucket=%s", config.Get("app.url"), bucket, filename, bucket)
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "上传成功", UploadResponse{
		Bucket:       bucket,
		OriginalName: header.Filename,
		FileName:     filename, // 使用保存后的文件名
		FileUrl:      fileURL,
		FileExt:      ext,
		FileSize:     header.Size,
	}, nil)
}

// 所有文件列表
func (fc *FileController) ListAllFiles(c *gin.Context) {
	files, err := fc.FileService.ListAllFiles()
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", files, nil)
}

// 文件列表
func (fc *FileController) ListFiles(c *gin.Context) {
	bucket := c.Query("bucket")
	if bucket == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("bucket不能为空"))
		return
	}
	files, err := fc.FileService.ListFiles(bucket)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", files, nil)

}

// 文件信息
func (fc *FileController) GetFileInfo(c *gin.Context) {
	bucket := c.Query("bucket")
	if bucket == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("bucket不能为空"))
		return
	}
	filename := c.Query("filename")
	if filename == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("filename不能为空"))
		return
	}
	fileSize, err := fc.FileService.GetFileInfo(bucket, filename)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", fileSize, nil)

}

// 下载文件
func (fc *FileController) DownloadFile(c *gin.Context) {
	bucket := c.PostForm("bucket")
	if bucket == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("bucket不能为空"))
		return
	}
	filename := c.Param("filename")
	reader, err := fc.FileService.LoadFile(bucket, filename)
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, reader)
}

// 删除文件
func (fc *FileController) DeleteFile(c *gin.Context) {
	bucket := c.Query("bucket")
	if bucket == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("bucket不能为空"))
		return
	}
	filename := c.Query("filename")
	if err := fc.FileService.DeleteFile(bucket, filename); err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "删除成功", nil, nil)
}
