/*
 * @PackageName: str_util
 * @FileName: str_util.go
 * @Description: 字符串工具
 * @Author: gabbymrh
 * @Date: 2024-07-18 16:15:39
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 16:15:39
 */

package str_util

import (
	"easy_dfs/pkg/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/sony/sonyflake"
	"github.com/spf13/cast"
	"strings"
)

// 雪花ID实例
var sf *sonyflake.Sonyflake

// 初始化
func init() {
	var st sonyflake.Settings
	// st.MachineID = awsutil.AmazonEC2MachineID
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		logger.ErrorString("id_generator", "init", "雪花ID生成器初始化失败")
	}
}

// SimpleUUID 简洁UUID生成器
func SimpleUUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// SnowFlakeID 雪花ID生成器(19位)
func SnowFlakeID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		logger.ErrorString("SnowFlakeID", "ID生成失败", err.Error())
	}
	// id 转字符串拼接上1后转为uint64
	formatId := "1" + cast.ToString(id)
	return cast.ToUint64(formatId)

}
