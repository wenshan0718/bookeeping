package util

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//数据库连接
var DB *sql.DB

/**
初始化数据库
*/
func MysqlInit() (err error) {
	fmt.Println("初始化数据库")
	//创建数据库连接
	db, err := sql.Open("mysql", "wen:123456@(localhost:3306)/bookkeeping")
	if err != nil {
		fmt.Println("初始化数据库连接失败", err)
		return
	}
	//连接数据库
	err = db.Ping()
	if err != nil {
		fmt.Println("数据库连接失败", err)
		return
	}
	//将连接赋值给全局变量
	DB = db
	return
}
