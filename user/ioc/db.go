package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lifelog-grpc/pkg/loggerx"
)

// InitMysql 初始化数据库连接
func InitMysql(l loggerx.Logger) *gorm.DB {
	// 数据库连接信息
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	config := Config{
		DSN: "root:123456@tcp(127.0.0.1:3306)/lifelog?charset=utf8mb4&parseTime=True&loc=Local",
	}
	// 将配置文件中，db.mysql下的所有配置项，绑定到结构体字段上
	err := viper.UnmarshalKey("db.mysql", &config)
	if err != nil {
		panic(err)
	}
	// 使用gorm的mysql驱动，建立数据库连接
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		// 只要初始化过程出错，程序直接结束
		panic(err)
	}
	return db
}
