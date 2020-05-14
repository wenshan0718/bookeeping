package model

//登录
import (
	"bookkeeping/entity"
	"bookkeeping/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

//登录页面
func LoginIndex(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./static/view/login.html")
	if err != nil {
		fmt.Println("返回登录页面错误", err)
		resMes := entity.NewResMes(500, "返回登录页面错误", "")
		resMes.ResMarshal(res)
		return
	}
	t.Execute(res, nil)
}

//登录
func Commit(res http.ResponseWriter, req *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(req.Body)
	var user entity.User
	// 解析参数 存入user
	decoder.Decode(&user)
	resMes := entity.NewResMes(200, "登录成功", "")
	//查询用户信息
	us, err := GetUserNP(user.Name, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			resMes = entity.NewResMes(500, "账户或密码错误", "")

		} else {
			fmt.Println("登录查询数据错误", err)
			resMes = entity.NewResMes(500, "登录查询数据错误", "")
		}
		resMes.ResMarshal(res)
		return
	}
	//将用户信息放入sesison保存到redis中
	//生成sessionId
	sesId := uuid.Must(uuid.NewV4()).String()
	//创建sesison
	session := entity.Session{
		SessionId: sesId,
		UserID:    us.Id,
		Name:      us.Name,
	}
	//放入redis中
	err = session.SaveSession()
	if err != nil {
		fmt.Println(err)
		resMes = entity.NewResMes(500, "保存session数据错误", "")
		resMes.ResMarshal(res)
		return
	}
	resMes.Data = sesId
	resMes.ResMarshal(res)
}

//注册
func Register(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var user entity.User
	decoder.Decode(&user)
	resMes := entity.NewResMes(200, "保存成功", "")
	err := SaveUser(&user)
	if err != nil {
		fmt.Println("注册用户失败", err)
		resMes = entity.NewResMes(500, "注册用户失败", "")
	}
	resMes.ResMarshal(res)
}

/**
通过cookie ,sessionId字段去Redis中查询session数据判断是否登录
登录返回true否则返回false
*/
func IfLogin(res http.ResponseWriter, req *http.Request) (bool, *entity.Session) {
	cookies := req.Header.Get("Cookie")
	sessionId := util.GetCookieByName(cookies, "sessionId")
	if sessionId == "" {
		//重定向到登录页面
		http.Redirect(res, req, "/login", 307)
		return false, nil
	}
	session, err := entity.GetSessionByCookie(sessionId)
	if err != nil {
		if err == entity.REDIS_SELECT_RESULT_NIL {
			//重定向到登录页面
			http.Redirect(res, req, "/login", 307)
			return false, nil
		} else {
			resMes := entity.NewResMes(500, "服务器错误", "")
			resMes.ResMarshal(res)
			fmt.Println("获取sessionId发生错误", err)
			return false, nil
		}
	}
	return true, session
}
