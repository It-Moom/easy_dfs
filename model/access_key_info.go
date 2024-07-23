/*
 * @PackageName: model
 * @FileName: access_key_info.go
 * @Description: 访问密钥信息
 * @Author: gabbymrh
 * @Date: 2024-07-19 17:25:23
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-19 17:25:23
 */

package model

// AccessKeyInfo 访问密钥信息
type AccessKeyInfo struct {
	// 访问密钥名称
	Name string `json:"name"`
	// 访问密钥
	AccessKey string `json:"accessKey"`
	// 秘钥密码
	SecretKey string `json:"secretKey"`
	// 过期时间
	ExpireTime string `json:"expireTime"`
	// 状态:1=启用,-1=禁用
	Status int `json:"status"`
}
