/*
 * @PackageName: services
 * @FileName: file_service.go
 * @Description:
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:22:26
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:22:26
 */

package services

import (
	"easy_dfs/pkg/config"
	"easy_dfs/pkg/filesystem"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// FileService 文件服务
type FileService struct {
	Storage       filesystem.FileSystemStorage // 文件系统存储接口，用于具体的文件操作
	BucketService BucketService                // 存储桶服务，用于获取存储桶信息
	mu            sync.Mutex                   // 互斥锁，用于保护文件操作的并发访问
}

const BasePath = "storage/" // 保存文件的基础路径

// SaveFile 将数据保存到指定的存储桶和文件名中
func (fs *FileService) SaveFile(bucket, filename string, data io.Reader) error {
	// 判断bucket是否合法
	if bucket == "" {
		return errors.New("bucket不能为空")
	}

	// bucket需在配置文件中存在
	_, err := fs.BucketService.FindBucketInfo(bucket)
	if err != nil {
		return err
	}

	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	filePath := fs.getFilePath(bucket, filename) // 获取文件保存路径
	return fs.Storage.Save(filePath, data)       // 调用存储接口保存文件
}

// LoadFile 加载指定存储桶和文件名的文件内容
func (fs *FileService) LoadFile(bucket, filename string) (io.Reader, error) {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	filePath := fs.getFilePath(bucket, filename) // 获取文件加载路径
	return fs.Storage.Load(filePath)             // 调用存储接口加载文件内容
}

// // LoadFileByPath 通过文件路径加载文件内容
// func (fs *FileService) LoadFileByPath(bucket, filePath string) (io.Reader, error) {
// 	fs.mu.Lock()                                     // 加锁以保护文件操作
// 	defer fs.mu.Unlock()                             // 延迟释放锁
// 	realFilePath := fs.getFilePath(bucket, filePath) // 获取文件加载路径
// 	return fs.Storage.Load(realFilePath)             // 调用存储接口加载文件内容
// }

func (fs *FileService) LoadFileByPath(bucket, filePath string) (io.Reader, error) {
	fs.mu.Lock()                                     // 加锁以保护文件操作
	defer fs.mu.Unlock()                             // 延迟释放锁
	realFilePath := fs.getFilePath(bucket, filePath) // 获取文件加载路径

	// // 打印调试信息
	// fmt.Println("Loading file from:", realFilePath)

	return fs.Storage.Load(realFilePath) // 调用存储接口加载文件内容
}

// 获取指定存储桶下的所有文件
func (fs *FileService) ListFiles(bucket string) ([]string, error) {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	return fs.Storage.ListFiles(bucket) // 调用存储接口列出文件
}

// 获取所有存储桶下的所有文件
func (fs *FileService) ListAllFiles() ([]filesystem.ResponseFileList, error) {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	return fs.Storage.ListAllFiles() // 调用存储接口列出所有文件

}

// 获取文件信息
func (fs *FileService) GetFileInfo(bucket, filename string) (filesystem.FileInfo, error) {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	filePath := fs.getFilePath(bucket, filename) // 获取文件信息路径
	return fs.Storage.GetFileInfo(filePath)      // 调用存储接口获取文件信息
}

// DeleteFile 删除指定存储桶和文件名的文件
func (fs *FileService) DeleteFile(bucket, filename string) error {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	filePath := fs.getFilePath(bucket, filename) // 获取文件删除路径
	return fs.Storage.Delete(filePath)           // 调用存储接口删除文件
}

// FileExists 检查指定存储桶和文件名的文件是否存在
func (fs *FileService) FileExists(bucket, filename string) (bool, error) {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	filePath := fs.getFilePath(bucket, filename) // 获取文件路径
	_, err := fs.Storage.Load(filePath)          // 尝试加载文件内容

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil // 文件不存在，返回 false
		}
		return false, err // 其他错误，返回错误信息
	}

	return true, nil // 文件存在，返回 true
}

// GetFileSize 获取指定存储桶和文件名的文件大小
func (fs *FileService) GetFileSize(bucket, filename string) (int64, error) {
	fs.mu.Lock()         // 加锁以保护文件操作
	defer fs.mu.Unlock() // 延迟释放锁

	filePath := fs.getFilePath(bucket, filename)  // 获取文件路径
	info, err := fs.Storage.GetFileInfo(filePath) // 获取文件信息

	if err != nil {
		return 0, err // 返回错误信息
	}

	return info.FileSize, nil // 返回文件大小
}

// // getFilePath 根据环境和存储桶名称生成文件路径
// func (fs *FileService) getFilePath(bucket, filename string) string {
// 	basePath := BasePath // 默认使用基础路径
//
// 	// 如果不是生产环境，使用临时路径
// 	if config.Get("app.env") != "prod" {
// 		basePath = "tmp/" + BasePath
// 	}
// 	return basePath + bucket + "/" + filename // 返回完整的文件路径
// }

// getFilePath 根据环境和存储桶名称生成文件路径
func (fs *FileService) getFilePath(bucket, filename string) string {
	basePath := BasePath // 默认使用基础路径

	// 如果不是生产环境，使用临时路径
	if config.Get("app.env") != "prod" {
		basePath = filepath.Join("tmp", BasePath)
	}

	// 使用 filepath.Join 来构建路径
	return filepath.Join(basePath, bucket, filename)
}
