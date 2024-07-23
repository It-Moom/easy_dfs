/*
 * @PackageName: controllers
 * @FileName: access_key_controller.go
 * @Description: 访问密钥控制器
 * @Author: gabbymrh
 * @Date: 2024-07-18 11:35:50
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 11:35:50
 */

package controllers

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/services"
	"easy_dfs/model"
	"easy_dfs/pkg/http/http_response"
	"errors"
	"github.com/gin-gonic/gin"
)

// AccessKeyController 访问密钥控制器
type AccessKeyController struct {
	AccessKeyService services.AccessKeyService
}

// CreateAccessKey 创建访问密钥
func (akc *AccessKeyController) CreateAccessKey(c *gin.Context) {
	var accessKeyInfo model.AccessKeyInfo
	if err := c.ShouldBindJSON(&accessKeyInfo); err != nil {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("参数有误"))
		return
	}

	if accessKeyInfo.Name == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("访问密钥名称不能为空"))
		return
	}
	// if accessKeyInfo.ExpireTime == "" {
	// 	http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("过期时间不能为空"))
	// 	return
	// }

	accessKey, secretKey, err := akc.AccessKeyService.GenerateAccessKey()
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}

	accessKeyInfo.AccessKey = accessKey
	accessKeyInfo.SecretKey = secretKey
	accessKeyInfo.Status = 1

	err = akc.AccessKeyService.SaveAccessKey(accessKeyInfo.Name, accessKeyInfo)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "新增成功", accessKeyInfo, nil)
}

// ListAccessKeys 获取访问密钥列表
func (akc *AccessKeyController) ListAccessKeys(c *gin.Context) {
	accessKeyList, err := akc.AccessKeyService.GetAccessKeyList()
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", accessKeyList, nil)
}

// GetAccessKeyInfo 获取访问密钥信息
func (akc *AccessKeyController) GetAccessKeyInfo(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("访问密钥名称不能为空"))
		return
	}
	accessKeyInfo, err := akc.AccessKeyService.GetAccessKey(name)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "获取成功", accessKeyInfo, nil)
}

// DeleteAccessKey 删除访问密钥
func (akc *AccessKeyController) DeleteAccessKey(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		http_response.Response(c, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("访问密钥名称不能为空"))
		return
	}
	err := akc.AccessKeyService.DeleteAccessKey(name)
	if err != nil {
		http_response.Response(c, response_code.REQUEST_FAILS, false, "操作失败", nil, err)
		return
	}
	http_response.Response(c, response_code.REQUEST_SUCCESS, true, "删除成功", nil, nil)
}
