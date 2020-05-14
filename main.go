package main

import (
	"bookkeeping/config"
	"bookkeeping/util"
	"fmt"
	"net/http"
	"time"
)

func main() {
	//初始化数据库
	err := util.MysqlInit()
	if err != nil {
		return
	}
	fmt.Println("初始化数据库完毕")
	//加载路由
	config.Router()
	//系统开始时初始化redis连接池
	err = util.InitPool("localhost:6379", 16, 0, 300*time.Second)
	if err != nil {
		return
	}
	fmt.Println("初始化redis完毕")

	//创建服务器参数
	server := &http.Server{
		Addr:           ":5031",          // 监听的TCP地址，如果为空字符串会使用":http"
		Handler:        nil,              // 调用的处理器，如为nil会调用http.DefaultServeMux
		ReadTimeout:    10 * time.Second, // 请求的读取操作在超时前的最大持续时间
		WriteTimeout:   10 * time.Second, // 回复的写入操作在超时前的最大持续时间
		MaxHeaderBytes: 0,                // 请求的头域最大长度，如为0则用DefaultMaxHeaderBytes
	}
	fmt.Println("开始启动服务")
	server.ListenAndServe()

}
