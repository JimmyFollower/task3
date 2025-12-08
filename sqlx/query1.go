package mysqlxtest_test

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
*

	假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	要求 ：
	编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

*

*/
/**

	假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	要求 ：
	定义一个 Book 结构体，包含与 books 表对应的字段。
	编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

**/
type employee struct {
	ID         int
	Name       string
	Department string
	Salary     int
}

func CreateTable(db *sqlx.DB) {
	createTableSql := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		department VARCHAR(255),
		salary INT
	);`

	createTableSql2 := `CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255),
		author VARCHAR(255),
		price FLOAT
	);`
	_, err := db.Exec(createTableSql)
	if err != nil {
		fmt.Println("创建表失败：", err)
		return
	}
	fmt.Println("创建表成功")

	_, err2 := db.Exec(createTableSql2)
	if err2 != nil {
		fmt.Println("创建表失败：", err2)
		return
	}
	fmt.Println("创建表成功")
}

func InsertData(db *sqlx.DB) {
	employeeList := []*employee{
		{1, "name1", "技术部", 6000},
		{2, "name2", "技术部", 9000},
		{3, "name3", "技术部", 10000},
		{4, "name4", "销售部", 6000},
		{5, "name5", "销售部", 7000},
	}
	for _, employee := range employeeList {
		insertSql := `
		INSERT INTO employees (name, department, salary) VALUES (?, ?, ?);`
		_, err := db.Exec(insertSql, employee.Name, employee.Department, employee.Salary)
		if err != nil {
			fmt.Println("插入数据失败：", err)
			return
		}
	}

}

// 自定义的 Employee 结构体切片
func QueryByDepartment(db *sqlx.DB) ([]employee, error) {
	var employeeList []employee
	querySql := `
	SELECT * FROM employees WHERE department = ?;`
	err := db.Select(&employeeList, querySql, "技术部")
	return employeeList, err

}
func QueryBySalary(db *sqlx.DB) (employee, error) {
	var employee employee
	querySql := `
	SELECT * FROM employees WHERE salary > ? ORDER BY salary DESC LIMIT 1;`
	err := db.Get(&employee, querySql, 50)
	return employee, err
}

func QueryByprice(db *sqlx.DB) ([]Books, error) {
	var selecbooks []Books

	err1 := db.Select(&selecbooks, "SELECT * FROM books where price>50")
	if err1 != nil {
		return nil, err1
	}

	return selecbooks, nil

}
func Run(db *sqlx.DB) {
	CreateTable(db)
	//InsertData(db)
	// employeeList, err := QueryByDepartment(db)
	// if err != nil {
	// 	fmt.Println("查询数据失败：", err)
	// 	return
	// }
	// fmt.Println("查询数据成功：", employeeList)
	// maxe1, err := QueryBySalary(db)
	// if err != nil {
	// 	fmt.Println("查询数据失败：", err)
	// 	return
	// }
	// fmt.Println("最高薪资的员工：", maxe1)

	//CreateTableBooks(db)

	bookList, err := QueryByprice(db)
	if err != nil {
		fmt.Println("查询数据失败：", err)
		return
	}
	fmt.Println("查询数据成功：", bookList)

}

//books

type Books struct {
	ID     int
	Title  string
	Author string
	Price  float64
}

func CreateTableBooks(db *sqlx.DB) {
	booksIns := []*Books{
		{1, "title1", "author1", 10.0},
		{2, "title2", "author2", 20.0},
		{3, "title3", "author3", 9.9},
		{4, "title4", "author4", 100.0},
		{5, "title5", "author5", 99999.99},
	}
	_, er := db.NamedExec("INSERT INTO books (id, title, author, price) VALUES (:id, :title, :author, :price)", booksIns)
	if er != nil {
		fmt.Println("批量插入测试数据异常：", er)
		return
	}
}
