package model

import (
	"bookkeeping/entity"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

//统计

//时间统计页面
func StatisticsTimeIndex(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	template, err := template.ParseFiles("./static/view/statisticsTime.html")
	if err != nil {
		fmt.Println("时间统计页面错误", err)
		resMes := entity.NewResMes(500, "访问错误", "")
		resMes.ResMarshal(res)
		return
	}
	template.Execute(res, nil)
}

//按时间统计记录
func SearchRecordOfTime(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	var dataMap map[string]int = make(map[string]int)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&dataMap)
	if err != nil {
		fmt.Println("时间统计记录错误", err)
		resMes := entity.NewResMes(500, "获取数据错误", "")
		resMes.ResMarshal(res)
		return
	}
	data, err := TimeSelectRecord(dataMap["timeType"], dataMap["page"], dataMap["pageSize"], session.UserID)
	if err != nil {
		fmt.Println("时间统计记录错误", err)
		resMes := entity.NewResMes(500, "获取数据错误", "")
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "查询成功", data)
	resMes.ResMarshal(res)
}

//分类统计页面
func StatisticsClassifyIndex(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	template, err := template.ParseFiles("./static/view/statisticsClassify.html")
	if err != nil {
		fmt.Println("时间统计页面错误", err)
		resMes := entity.NewResMes(500, "访问错误", "")
		resMes.ResMarshal(res)
		return
	}
	template.Execute(res, nil)
}

//按分类统计记录
func SearchRecordOfClassify(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	var dataMap map[string]int = make(map[string]int)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&dataMap)
	if err != nil {
		fmt.Println("分类统计记录错误", err)
		resMes := entity.NewResMes(500, "获取数据错误", "")
		resMes.ResMarshal(res)
		return
	}
	data, err := ClassifySelectRecord(dataMap["classifyType"], dataMap["page"], dataMap["pageSize"], session.UserID)
	if err != nil {
		fmt.Println("分类统计记录错误", err)
		resMes := entity.NewResMes(500, "获取数据错误", "")
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "查询成功", data)
	resMes.ResMarshal(res)
}
