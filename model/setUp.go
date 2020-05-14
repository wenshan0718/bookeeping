//设置
package model

import (
	"bookkeeping/entity"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

/**
设置页面
*/
func SetUpIndex(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	t, err := template.ParseFiles("./static/view/setUp.html")
	if err != nil {
		fmt.Println("返回设置页面错误", err)
		resMes := entity.NewResMes(500, "返回设置页面错误", "")
		resMes.ResMarshal(res)
		return
	}
	t.Execute(res, nil)
}

/**
保存分类
*/
func ClassifySave(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	//使用json解析器获取前端post传过来的参数
	decode := json.NewDecoder(req.Body)
	var class entity.Classify
	err := decode.Decode(&class)
	if err != nil {
		fmt.Println("保存分类发生错误", err)
		resMes := entity.NewResMes(500, "保存分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	//将用户id放入class里面
	class.UserId = session.UserID
	err = SaveClassifyDao(&class)
	if err != nil {
		fmt.Println("保存分类发生错误", err)
		resMes := entity.NewResMes(500, "保存分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "保存分类成功", nil)
	resMes.ResMarshal(res)
}

/**
获取分类
*/
func SelectClassify(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	//使用json解析器获取前端post传过来的参数
	decode := json.NewDecoder(req.Body)
	var class entity.Classify
	err := decode.Decode(&class)
	if err != nil {
		fmt.Println("获取分类发生错误", err)
		resMes := entity.NewResMes(500, "获取分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	//将用户id放入class里面
	class.UserId = session.UserID
	calssifys, err := SelectClassifyDao(&class)
	if err != nil {
		fmt.Println("获取分类发生错误", err)
		resMes := entity.NewResMes(500, "获取分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "获取分类成功", calssifys)
	resMes.ResMarshal(res)
}

/**
获取分组成支出和收入的分类
*/
func GetGroupClassify(res http.ResponseWriter, req *http.Request) {
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	data, err := ClassifyGroup(session.UserID)
	if err != nil {
		fmt.Println("获取分类好的分类")
		resMes := entity.NewResMes(500, "获取分类好的分类错误", "")
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "获取成功", data)
	resMes.ResMarshal(res)
}

/**
编辑分类
*/
func ClassifyAddEdit(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	//使用json解析器获取前端post传过来的参数
	decode := json.NewDecoder(req.Body)
	var class entity.Classify
	err := decode.Decode(&class)
	if err != nil {
		fmt.Println("修改分类发生错误", err)
		resMes := entity.NewResMes(500, "修改分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	err = UpdateClassify(&class)
	if err != nil {
		fmt.Println("修改分类发生错误", err)
		resMes := entity.NewResMes(500, "修改分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "修改分类成功", nil)
	resMes.ResMarshal(res)
}

/**
删除分类
*/
func ClassifyAddDel(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, _ := IfLogin(res, req)
	if !ifLogin {
		return
	}
	//使用json解析器获取前端post传过来的参数
	decode := json.NewDecoder(req.Body)
	var class entity.Classify
	err := decode.Decode(&class)
	if err != nil {
		fmt.Println("删除分类发生错误", err)
		resMes := entity.NewResMes(500, "删除分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	err = DelClassify(&class)
	if err != nil {
		fmt.Println("删除分类发生错误", err)
		resMes := entity.NewResMes(500, "删除分类发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "删除分类成功", nil)
	resMes.ResMarshal(res)
}

/**
修改密码
*/
func AlterSave(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	//使用json解析器获取前端post传过来的参数
	decode := json.NewDecoder(req.Body)
	var dataMap map[string]string = make(map[string]string)
	err := decode.Decode(&dataMap)
	if err != nil {
		fmt.Println("修改密码发生错误", err)
		resMes := entity.NewResMes(500, "修改密码发生错误", nil)
		resMes.ResMarshal(res)
		return
	}

	user, err := GetUserId(session.UserID)
	if err != nil {
		fmt.Println("修改密码发生错误", err)
		resMes := entity.NewResMes(500, "修改密码发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	//检验密码是否相等
	if user.Password != dataMap["oldPassword"] {
		resMes := entity.NewResMes(500, "修改密码失败原密码错误", nil)
		resMes.ResMarshal(res)
		return
	}
	err = AlterUserPassword(session.UserID, dataMap["newPassword"])
	if err != nil {
		fmt.Println("修改密码发生错误", err)
		resMes := entity.NewResMes(500, "修改密码发生错误", nil)
		resMes.ResMarshal(res)
		return
	}
	resMes := entity.NewResMes(200, "修改密码成功", nil)
	resMes.ResMarshal(res)
}
