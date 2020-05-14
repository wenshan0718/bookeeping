package entity

/**
记录实体类
*/
type Record struct {
	//id主键
	Id int `json:"id"`
	//金额
	Monery float32 `json:"monery"`
	//一级分类
	OneClassify     int    `json:"oneClassify"`
	OneClassifyName string `json:"oneClassifyName"`
	//二级分类
	TwoClassify     int    `json:"twoClassify"`
	TwoClassifyName string `json:"twoClassifyName"`
	//记录时间
	Time string `json:"time"`
	//备注
	Remark string `json:"remark"`
	//支出还是收入1支出2收入
	Outin int `json:"outin"`
	//用户id
	UserId int `json:"userId"`
	//查询条件开始时间
	StartTime string `json:"startTime`
	//查询条件结束时间
	EndTime string `json:"endTime"`
}
