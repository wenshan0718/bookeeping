package entity

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
后台返货前端数据类型
*/
type ResMes struct {
	Status int         `json:"status"` //状态
	Mes    string      `json:"mes"`    //返回信息
	Data   interface{} `json:"data"`   //返回的数据转换成json格式保存成字符串
}

/**
根据传入的参数创建RemMes 并返回
status  状态
mes 返回信息
data 返回的数据
*/
func NewResMes(status int, mes string, data interface{}) *ResMes {
	resMes := ResMes{
		Status: status,
		Mes:    mes,
		Data:   data,
	}
	return &resMes
}

/**
将数据序列化后返回给前端
*/
func (resMes *ResMes) ResMarshal(res http.ResponseWriter) {
	data, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("数据返回前端时序列化错误", err)
		return
	}
	res.Write(data)
}
