package model

//主页
import (
	"bookkeeping/entity"
	"fmt"
	"html/template"
	"net/http"
)

func Master(res http.ResponseWriter, req *http.Request) {
	//先判读是否有进行登录
	ifLogin, session := IfLogin(res, req)
	if !ifLogin {
		return
	}
	template, err := template.ParseFiles("./static/view/master.html")
	if err != nil {
		fmt.Println("解析主页模板错误")
		resMes := entity.NewResMes(500, "访问错误", "")
		resMes.ResMarshal(res)
		return
	}
	var result = map[string]float32{
		"dayOut":   0,
		"dayIn":    0,
		"monthOut": 0,
		"monthIn":  0,
		"yearOut":  0,
		"yearIn":   0,
	}
	//获取本日数据
	data, err := NowTimeSelectRecord(1, session.UserID)
	if err == nil {
		for _, v := range *data {
			if v.Outin == 1 {
				result["dayOut"] = v.Monery
			} else if v.Outin == 2 {
				result["dayIn"] = v.Monery
			}
		}
	} else {
		fmt.Println("主页获取本日数据错误", err)
	}
	//获取本月数据
	data, err = NowTimeSelectRecord(2, session.UserID)
	if err == nil {
		for _, v := range *data {
			if v.Outin == 1 {
				result["monthOut"] = v.Monery
			} else if v.Outin == 2 {
				result["monthIn"] = v.Monery
			}
		}
	} else {
		fmt.Println("主页获取本月数据错误", err)
	}
	//获取本年数据
	data, err = NowTimeSelectRecord(3, session.UserID)
	if err == nil {
		for _, v := range *data {
			if v.Outin == 1 {
				result["yearOut"] = v.Monery
			} else if v.Outin == 2 {
				result["yearIn"] = v.Monery
			}
		}
	} else {
		fmt.Println("主页获取本年数据错误", err)
	}

	template.Execute(res, result)
}
