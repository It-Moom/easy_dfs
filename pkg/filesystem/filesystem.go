/*
 * @PackageName: filesystem
 * @FileName: filesystem.go
 * @Description: 文件系统
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:17:58
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:17:58
 */

package filesystem

import (
	"easy_dfs/pkg/config"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemStorage struct {
	BaseDir string
}

type ResponseFileList struct {
	Bucket   string   `json:"bucket"`
	FileList []string `json:"fileList"`
}

func (s *FileSystemStorage) initBaseDir() {
	// 根据环境判断是否使用tmp目录
	if config.Get("app.env") != "prod" {
		s.BaseDir = filepath.Join("tmp", "storage")
	} else {
		s.BaseDir = "storage"
	}
}

// Save 保存文件
func (s *FileSystemStorage) Save(filename string, data io.Reader) error {
	path := filepath.Join(s.BaseDir, filename)
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, data)
	return err
}

// Load 加载文件
func (s *FileSystemStorage) Load(filename string) (io.Reader, error) {
	path := filepath.Join(s.BaseDir, filename)
	return os.Open(path)
}

// Delete 删除文件
func (s *FileSystemStorage) Delete(filename string) error {
	path := filepath.Join(s.BaseDir, filename)
	return os.Remove(path)
}

// DeleteByPath 根据bucketName和filePath删除文件
func (s *FileSystemStorage) DeleteByPath(bucketName, filePath string) error {
	path := filepath.Join(s.BaseDir, bucketName, filePath)
	return os.Remove(path)
}

// ListFiles 根据bucketName列出文件
func (s *FileSystemStorage) ListFiles(bucketName string) ([]string, error) {
	s.initBaseDir()
	return s.listFilesRecursive(bucketName, bucketName)
}

// ListAllFiles 列出所有文件
func (s *FileSystemStorage) ListAllFiles() ([]ResponseFileList, error) {
	s.initBaseDir()
	buckets, err := os.ReadDir(s.BaseDir)
	if err != nil {
		return nil, err
	}

	var responseFileLists []ResponseFileList
	for _, bucket := range buckets {
		if bucket.IsDir() {
			bucketName := bucket.Name()
			files, err := s.listFilesRecursive(bucketName, bucketName)
			if err != nil {
				return nil, err
			}
			responseFileLists = append(responseFileLists, ResponseFileList{
				Bucket:   bucketName,
				FileList: files,
			})
		}
	}
	return responseFileLists, nil
}

// 列出文件
func (s *FileSystemStorage) listFilesRecursive(basePath, bucketName string) ([]string, error) {
	path := filepath.Join(s.BaseDir, basePath)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		fullPath := filepath.Join(basePath, file.Name())
		if file.IsDir() {
			subFiles, err := s.listFilesRecursive(fullPath, bucketName)
			if err != nil {
				return nil, err
			}
			fileNames = append(fileNames, subFiles...)
		} else {
			relativePath := strings.TrimPrefix(fullPath, s.BaseDir+"/")
			fileNames = append(fileNames, relativePath)
		}
	}
	return fileNames, nil
}

type FileInfo struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileExt  string `json:"fileExt"`
	FileSize int64  `json:"fileSize"`
}

// GetFileInfo 获取文件信息
func (s *FileSystemStorage) GetFileInfo(path string) (FileInfo, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, nil
	}
	info := FileInfo{
		FileName: filepath.Base(path),
		FilePath: path,
		FileExt:  filepath.Ext(path),
		FileSize: fileInfo.Size(),
	}
	return info, nil
}
