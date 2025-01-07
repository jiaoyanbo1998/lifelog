package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lifelog-grpc/pkg/loggerx"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetMysql(l loggerx.Logger) *gorm.DB {
	once.Do(func() {
		db = initMysql(l)
	})
	return db
}

// initMysql 初始化数据库连接
func initMysql(l loggerx.Logger) *gorm.DB {
	// 数据库连接信息
	/*
		用户名：root
		密码：123456
		主机地址：127.0.0.1
		端口：3306
		数据库名：webook
		字符集：utf8mb4
		时间解析：True
		时区：Local
	*/
	// 数据库连接信息
	dsn := "root:123456@tcp(127.0.0.1:3306)/lifelog?charset=utf8mb4&parseTime=True&loc=Local"
	// 使用gorm的mysql驱动，建立数据库连接
	d, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return d
}
