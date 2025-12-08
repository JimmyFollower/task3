package utils

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBUtil struct {
}

func (d DBUtil) Connect() *gorm.DB {
	dsn := "root:123123@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // 必须传gorm.Config{}（GORM v2要求）

	if err != nil {
		panic(err)
	}
	return db
}
