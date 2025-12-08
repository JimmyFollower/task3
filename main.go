package main

import (
	blog "study/gorm"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//T1 grom框架的使用
	// db, err := gorm.Open(mysql.Open("root:123123@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	// if err != nil {
	// 	log.Fatal("数据库连接失败", err)
	// }
	//crud.Run(db)
	//transaction.Run(db)

	//T2 sqlx框架使用
	// dsn := "root:123123@tcp(127.0.0.1:3306)/sqlx?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := sqlx.Connect("mysql", dsn)
	// if err != nil {
	// 	log.Fatal("数据库连接失败", err)
	// }
	// mysqlxtest_test.Run(db)

	db, err := gorm.Open(mysql.Open("root:123123@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	blog.Run(db)
	if err != nil {
		panic(err)
	}
	//T3 进阶gorm 使用

}
