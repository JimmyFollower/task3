package main

import (
	"BLOG/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:123123@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // 必须传gorm.Config{}（GORM v2要求）

	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	//db.AutoMigrate(&models.Users{}, &models.Posts{}, &models.Comments{}, models.Token{})
	_ = db
	r := gin.Default()
	controller.UserControllerInit(r)
	controller.PostControllerInit(r)
	controller.CommentControllerInit(r)
	r.Run(":8080")
	//defer db.Close()

}
