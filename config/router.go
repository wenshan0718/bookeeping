package config

import (
	"bookkeeping/model"
	"net/http"
)

/**
路由
*/
func Router() {
	//设置处理静态资源
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//登录页面
	http.HandleFunc("/login", model.LoginIndex)
	//登录
	http.HandleFunc("/login/commit", model.Commit)
	//注册
	http.HandleFunc("/login/register", model.Register)

	//主页
	http.HandleFunc("/master", model.Master)

	//设置页面
	http.HandleFunc("/toSetUp", model.SetUpIndex)
	//保存分类
	http.HandleFunc("/setUp/classifySave", model.ClassifySave)
	//获取分类
	http.HandleFunc("/setUp/SelectClassify", model.SelectClassify)
	//获取分组成支出和收入的分类
	http.HandleFunc("/setUp/getGroupClassify", model.GetGroupClassify)
	//编辑分类
	http.HandleFunc("/setUp/classifyAddEdit", model.ClassifyAddEdit)
	//删除分类
	http.HandleFunc("/setUp/classifyAddDel", model.ClassifyAddDel)
	//修改密码
	http.HandleFunc("/setUp/alterSave", model.AlterSave)

	//记录页面
	http.HandleFunc("/record/index", model.RecordIndex)
	//保存记录
	http.HandleFunc("/record/SaveRecord", model.SaveRecord)
	//修改记录
	http.HandleFunc("/record/editRecord", model.EditRecord)
	//删除记录
	http.HandleFunc("/record/delRecord", model.DelRecord)
	//流水页面
	http.HandleFunc("/record/runningWater", model.RunningWater)
	//查询记录
	http.HandleFunc("/record/searchRecord", model.SearchRecord)
	//根据id查询记录
	http.HandleFunc("/record/SearchRecordById", model.SearchRecordById)

	///时间统计页面
	http.HandleFunc("/statistics/StatisticsTimeIndex", model.StatisticsTimeIndex)
	//根据时间统计数据
	http.HandleFunc("/statistics/SearchRecordOfTime", model.SearchRecordOfTime)
	///分类统计页面
	http.HandleFunc("/statistics/StatisticsClassifyIndex", model.StatisticsClassifyIndex)
	//按分类统计记录
	http.HandleFunc("/statistics/SearchRecordOfClassify", model.SearchRecordOfClassify)
}
