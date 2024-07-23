/*
 * @PackageName: model
 * @FileName: bucket_info.go
 * @Description: bucket信息
 * @Author: gabbymrh
 * @Date: 2024-07-18 17:21:31
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 17:21:31
 */

package model

type BucketInfo struct {
	// 存储桶名称
	Name string `json:"name"`
	// 访问策略
	AccessPolicy string `json:"accessPolicy"`
	// 存储类型
	StorageType string `json:"storageType"`
}
