package model

//记录
import (
	"bookkeeping/entity"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//流水页面
func RunningWater(res http.ResponseWriter, req *http.Request) {
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	t, err := template.ParseFiles("./static/view/runningWater.html")
	if err != nil {
		fmt.Println("返回登录页面错误", err)
		resMes := entity.NewResMes(500, "返回登录页面错误", "")
		resMes.ResMarshal(res)
		return
	}
	t.Execute(res, nil)
}

//记录页面
func RecordIndex(res http.ResponseWriter, req *http.Request) {
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	t, err := template.ParseFiles("./static/view/record.html")
	if err != nil {
		fmt.Println("返回登录页面错误", err)
		resMes := entity.NewResMes(500, "返回登录页面错误", "")
		resMes.ResMarshal(res)
		return
	}
	var data map[string]string = make(map[string]string)
	data["type"] = "add"
	rawQuery := req.URL.RawQuery
	index := strings.Index(rawQuery, "id=")
	if index != -1 {
		data["type"] = "edit"
		start := index + 3
		end := strings.Index(string([]byte(rawQuery)[start:]), "&")
		if end == -1 {
			data["id"] = string([]byte(rawQuery)[start:])
		} else {
			data["id"] = string([]byte(rawQuery)[start : end+start])
		}
	}
	t.Execute(res, data)
}

//保存记录
func SaveRecord(res http.ResponseWriter, req *http.Request) {
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var record entity.Record
	err := decoder.Decode(&record)
	if err != nil {
		fmt.Println("保存记录发生错误", err)
		resMes := entity.NewResMes(500, "保存发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	record.UserId = session.UserID
	err = InsertRecord(&record)
	if err != nil {
		fmt.Println("保存记录发生错误", err)
		resMes := entity.NewResMes(500, "保存发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "保存成功", nil)
	resMes.ResMarshal(res)
}

//修改记录
func EditRecord(res http.ResponseWriter, req *http.Request) {
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var record entity.Record
	err := decoder.Decode(&record)
	if err != nil {
		fmt.Println("保存记录发生错误", err)
		resMes := entity.NewResMes(500, "保存发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	err = UpdateRecord(&record)
	if err != nil {
		fmt.Println("保存记录发生错误", err)
		resMes := entity.NewResMes(500, "保存发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "保存成功", nil)
	resMes.ResMarshal(res)
}

//查询记录
func SearchRecord(res http.ResponseWriter, req *http.Request) {

	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var mapdata map[string]int = make(map[string]int)
	err := decoder.Decode(&mapdata)
	if err != nil {
		fmt.Println("查询记录发生错误", err)
		resMes := entity.NewResMes(500, "查询发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	result, err := SelectRecord(session.UserID, mapdata["page"], mapdata["pageSize"])
	if err != nil {
		fmt.Println("查询记录发生错误", err)
		resMes := entity.NewResMes(500, "查询发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "查询成功", result)
	resMes.ResMarshal(res)
}

//删除记录
func DelRecord(res http.ResponseWriter, req *http.Request) {

	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var mapdata map[string]int = make(map[string]int)
	err := decoder.Decode(&mapdata)
	if err != nil {
		fmt.Println("删除记录发生错误", err)
		resMes := entity.NewResMes(500, "删除发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	err = DeleteRecordById(mapdata["id"])
	if err != nil {
		fmt.Println("删除记录发生错误", err)
		resMes := entity.NewResMes(500, "删除查询发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "删除成功", nil)
	resMes.ResMarshal(res)
}

//根据id查询记录
func SearchRecordById(res http.ResponseWriter, req *http.Request) {

	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var mapdata map[string]int = make(map[string]int)
	err := decoder.Decode(&mapdata)
	if err != nil {
		fmt.Println("查询记录发生错误", err)
		resMes := entity.NewResMes(500, "查询发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	result, err := SelectRecordById(session.UserID, mapdata["id"])
	if err != nil {
		fmt.Println("查询记录发生错误", err)
		resMes := entity.NewResMes(500, "查询发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "查询成功", result)
	resMes.ResMarshal(res)
}
