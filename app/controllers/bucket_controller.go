/*
 * @PackageName: controllers
 * @FileName: bucket_controller.go
 * @Description: 文件桶控制器
 * @Author: gabbymrh
 * @Date: 2024-07-18 11:35:24
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 11:35:24
 */

package controllers

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/enum/system_default"
	"easy_dfs/app/services"
	"easy_dfs/model"
	"easy_dfs/pkg/http/http_response"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// 存储桶控制器
type BucketController struct {
	BucketService services.BucketService
}

// CreateBucket 创建存储桶
func (bc *BucketController) CreateBucket(c *gin.Context) {
	// 获取存储桶名称
	var bucketInfo model.BucketInfo
	if err := c.ShouldBindJSON(&bucketInfo); err != nil {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("参数有误"))
		return
	}

	if bucketInfo.Name == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("存储桶名称不能为空"))
		return
	}
	if strings.Contains(bucketInfo.Name, system_default.STORAGE_PATH) {
		http_response.Response(c, response_code.REQUEST_DENIED, false, "操作失败", nil, errors.New("存储桶名称非法"))
		return
	}
	if bucketInfo.AccessPolicy == "" {
		bucketInfo.AccessPolicy = "private"
	}
	err := bc.BucketService.CreateBucket(bucketInfo.Name, bucketInfo.AccessPolicy)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "新增成功", nil, nil)
}

// 获取存储桶列表
func (bc *BucketController) ListBuckets(c *gin.Context) {
	bucketList, err := bc.BucketService.GetBucketList()
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", bucketList, nil)
}

// 获取存储桶信息
func (bc *BucketController) GetBucketInfo(c *gin.Context) {
	bucketName := c.Query("bucket")
	if bucketName == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("存储桶名称不能为空"))
		return
	}
	bucketInfo, err := bc.BucketService.FindBucketInfo(bucketName)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", bucketInfo, nil)

}

// 删除存储桶
func (bc *BucketController) DeleteBucket(c *gin.Context) {
	// 获取存储桶名称
	bucketName := c.Query("bucket")
	if bucketName == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("存储桶名称不能为空"))
		return
	}
	err := bc.BucketService.DeleteBucket(bucketName)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "删除成功", nil, nil)
}
