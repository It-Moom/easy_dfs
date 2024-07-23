# EasyDFS
> 简洁易用的对象存储系统

## 使用说明
### 1.准备目录
创建或指定一个目录作为存放该系统的根目录
### 2.准备配置
将 `easy_dfs` 可执行文件和 `app.yml` 配置文件复制到该目录下，可自定义配置 `app.yml` 内的配置信息，比如端口
### 3.运行系统
直接执行 `easy_dfs` 可执行文件即可，比如 `./easy_dfs.exe` 或 `./easy_dfs`
### 4.生成密钥
根据API接口中的`AccessKey`目录对应方法生成密钥对
### 5.创建Bucket
根据API接口中的`Bucket`目录创建对应存储桶
### 5.上传/删除文件
根据API接口中的`File`目录调取对应接口上传或删除文件

## 开发说明
- 拉取代码到本地,并将`app.yml.example`复制为`app.yml`
- 安装依赖 `go mod tidy`
- 运行 `go run main.go` 或 `go build` 编译后运行
- 访问 `http://localhost:18088` 查看接口文档
- 代码中的 `app.yml` 是配置文件，可自定义配置

## 其他说明
- 非生产环境时上传的文件和配置都在 `tmp/` 下
- 生产环境上传的文件在 `storage/` 下,配置文件在 `config/` 下