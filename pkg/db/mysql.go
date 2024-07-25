package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type MySQLOptions struct {
	Host     string
	UserName string
	Password string
	Database string

	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	MaxConnectionIdleTime time.Duration

	LogLevel int
	//Logger   logger.Interface
}

// NewMySQLClient 通过MysqlOptions struct做传值
func NewMySQLClient(opts *MySQLOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		opts.UserName,
		opts.Password,
		opts.Host,
		opts.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	sqlDB.SetConnMaxIdleTime(opts.MaxConnectionIdleTime)

	return db, nil
}

// NewMySQLClientWith TODO：通过Options模式创建MySQL客户端
//func NewMySQLClientWith() {
//
//}
