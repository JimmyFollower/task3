package crud

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

//题目1：基本CRUD操作,
// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、
// grade （学生年级，字符串类型）。

// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。,
// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。,
// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。,
// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

// 学生结构体包含id 姓名 年龄 年级
type Student struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Age   int
	Grade string
}

func Run(db *gorm.DB) {
	// 创建表
	db.AutoMigrate(&Student{})
	// 插入数据
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	db.Create(&student)
	// 查询数据
	var students []Student
	res := db.Where("age > ?", 18).Find(&students)
	if res.Error != nil {
		log.Println("查询数据失败：", res.Error)
		return
	}
	for _, student := range students {
		fmt.Println("大于18岁的学生信息：", student)
	}
	// 更新数据
	res2 := db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
	fmt.Println("更新数据影响的结果：", res2.RowsAffected)
	// 删除数据
	res3 := db.Where("age < ?", 15).Delete(&Student{})
	fmt.Println("删除数据影响结果：", res3.RowsAffected)
}
