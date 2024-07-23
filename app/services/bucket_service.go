/*
 * @PackageName: services
 * @FileName: bucket_service.go
 * @Description:
 * @Author: gabbymrh
 * @Date: 2024-07-18 17:23:26
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 17:23:26
 */

package services

import (
	"easy_dfs/model"
	"easy_dfs/pkg/config"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// BucketService 存储桶服务
type BucketService struct {
	mu sync.Mutex
}

// getConfigPath 根据应用环境返回存储桶配置文件的路径
func (bs *BucketService) getConfigPath() string {
	bucketConfigPath := "config/bucket.json"
	if config.Get("app.env") != "prod" {
		bucketConfigPath = "tmp/config/bucket.json"
	}
	return bucketConfigPath
}

// readBucketConfig 从文件中读取存储桶配置并返回解析后的数据
func (bs *BucketService) readBucketConfig() ([]model.BucketInfo, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bucketConfigPath := bs.getConfigPath()

	byteValue, err := os.ReadFile(bucketConfigPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var bucketConfig []model.BucketInfo
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &bucketConfig)
		if err != nil {
			return nil, err
		}
	}

	return bucketConfig, nil
}

// writeBucketConfig 将存储桶配置写入文件
func (bs *BucketService) writeBucketConfig(bucketConfig []model.BucketInfo) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bucketConfigPath := bs.getConfigPath()

	bucketConfigJson, err := json.Marshal(bucketConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(bucketConfigPath, bucketConfigJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CreateBucket 创建新的存储桶，如果已存在则返回错误
func (bs *BucketService) CreateBucket(bucketName string, accessPolicy string) error {
	bucketConfig, err := bs.readBucketConfig()
	if err != nil {
		return err
	}

	for _, v := range bucketConfig {
		if v.Name == bucketName {
			return errors.New("bucket 已存在")
		}
	}

	bucketConfig = append(bucketConfig, model.BucketInfo{Name: bucketName, AccessPolicy: accessPolicy})

	return bs.writeBucketConfig(bucketConfig)
}

// GetBucketList 返回所有存储桶配置的列表
func (bs *BucketService) GetBucketList() ([]model.BucketInfo, error) {
	return bs.readBucketConfig()
}

// FindBucketInfo 根据存储桶名称查找并返回存储桶配置信息
func (bs *BucketService) FindBucketInfo(bucketName string) (*model.BucketInfo, error) {
	bucketConfig, err := bs.readBucketConfig()
	if err != nil {
		return nil, err
	}

	for _, v := range bucketConfig {
		if v.Name == bucketName {
			return &v, nil
		}
	}

	return nil, errors.New("bucket 不存在")
}

// DeleteBucket 根据存储桶名称删除存储桶配置
func (bs *BucketService) DeleteBucket(bucketName string) error {
	bucketConfig, err := bs.readBucketConfig()
	if err != nil {
		return err
	}

	for k, v := range bucketConfig {
		if v.Name == bucketName {
			bucketConfig = append(bucketConfig[:k], bucketConfig[k+1:]...)
			break
		}
	}

	return bs.writeBucketConfig(bucketConfig)
}
