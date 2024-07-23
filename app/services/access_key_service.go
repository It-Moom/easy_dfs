/*
 * @PackageName: services
 * @FileName: access_key_service.go
 * @Description: 访问密钥服务
 * @Author: gabbymrh
 * @Date: 2024-07-19 17:22:14
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-19 17:22:14
 */

package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"easy_dfs/model"
	"easy_dfs/pkg/config"
)

// AccessKeyService 访问密钥服务
type AccessKeyService struct {
	mu sync.Mutex
}

// getConfigPath 根据应用环境返回访问密钥配置文件的路径
func (aks *AccessKeyService) getConfigPath() string {
	accessKeyConfigPath := "config/access_key.json"
	if config.Get("app.env") != "prod" {
		accessKeyConfigPath = "tmp/config/access_key.json"
	}
	return accessKeyConfigPath
}

// readAccessKeyConfig 从文件中读取访问密钥配置并返回解析后的数据
func (aks *AccessKeyService) readAccessKeyConfig() ([]model.AccessKeyInfo, error) {
	aks.mu.Lock()
	defer aks.mu.Unlock()

	accessKeyConfigPath := aks.getConfigPath()

	byteValue, err := os.ReadFile(accessKeyConfigPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var accessKeyConfig []model.AccessKeyInfo
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &accessKeyConfig)
		if err != nil {
			return nil, err
		}
	}

	return accessKeyConfig, nil
}

// writeAccessKeyConfig 将访问密钥配置写入文件
func (aks *AccessKeyService) writeAccessKeyConfig(accessKeyConfig []model.AccessKeyInfo) error {
	aks.mu.Lock()
	defer aks.mu.Unlock()

	accessKeyConfigPath := aks.getConfigPath()

	accessKeyConfigJson, err := json.Marshal(accessKeyConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(accessKeyConfigPath, accessKeyConfigJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GenerateAccessKey 生成一个新的访问密钥和秘钥
func (aks *AccessKeyService) GenerateAccessKey() (string, string, error) {
	accessBytes := make([]byte, 16)
	secretBytes := make([]byte, 32)

	if _, err := rand.Read(accessBytes); err != nil {
		return "", "", err
	}
	if _, err := rand.Read(secretBytes); err != nil {
		return "", "", err
	}

	accessKey := hex.EncodeToString(accessBytes)
	secretKey := hex.EncodeToString(secretBytes)
	return accessKey, secretKey, nil
}

// SaveAccessKey 保存访问密钥
func (aks *AccessKeyService) SaveAccessKey(userID string, accessKeyInfo model.AccessKeyInfo) error {
	accessKeyConfig, err := aks.readAccessKeyConfig()
	if err != nil {
		return err
	}

	for _, v := range accessKeyConfig {
		if v.Name == userID {
			return errors.New("name 已存在")
		}
	}

	accessKeyConfig = append(accessKeyConfig, accessKeyInfo)

	return aks.writeAccessKeyConfig(accessKeyConfig)
}

// GetAccessKeyList 获取访问密钥列表
func (aks *AccessKeyService) GetAccessKeyList() ([]model.AccessKeyInfo, error) {
	return aks.readAccessKeyConfig()
}

// GetAccessKey 获取访问密钥
func (aks *AccessKeyService) GetAccessKey(userID string) (*model.AccessKeyInfo, error) {
	accessKeyConfig, err := aks.readAccessKeyConfig()
	if err != nil {
		return nil, err
	}

	for _, v := range accessKeyConfig {
		if v.Name == userID {
			return &v, nil
		}
	}

	return nil, errors.New("name 不存在")
}

// 校验访问密钥
func (aks *AccessKeyService) CheckAccessKey(accessKey string, secretKey string) bool {
	accessKeyConfig, err := aks.readAccessKeyConfig()
	if err != nil {
		return false
	}

	for _, v := range accessKeyConfig {
		if v.AccessKey == accessKey && v.SecretKey == secretKey {
			return true
		}
	}

	return false
}

// DeleteAccessKey 删除访问密钥
func (aks *AccessKeyService) DeleteAccessKey(userID string) error {
	accessKeyConfig, err := aks.readAccessKeyConfig()
	if err != nil {
		return err
	}

	for k, v := range accessKeyConfig {
		if v.Name == userID {
			accessKeyConfig = append(accessKeyConfig[:k], accessKeyConfig[k+1:]...)
			break
		}
	}

	return aks.writeAccessKeyConfig(accessKeyConfig)
}

// CreateAndSaveAccessKey 创建并保存访问密钥
func (aks *AccessKeyService) CreateAndSaveAccessKey(userID string, name string, expireTime string) (model.AccessKeyInfo, error) {
	accessKey, secretKey, err := aks.GenerateAccessKey()
	if err != nil {
		return model.AccessKeyInfo{}, err
	}

	accessKeyInfo := model.AccessKeyInfo{
		Name:       name,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		ExpireTime: expireTime,
		Status:     1,
	}

	err = aks.SaveAccessKey(userID, accessKeyInfo)
	if err != nil {
		return model.AccessKeyInfo{}, err
	}

	return accessKeyInfo, nil
}
