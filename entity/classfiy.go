package entity

/**
分类实体类
*/
type Classify struct {
	//分类ID
	Id int `json:"id"`
	//分类名称
	Name string `json:"name"`
	//分类层级1位一级分类2为二级分类
	Layer int `json:"layer"`
	//二级时的父类id
	ParentId int `json:"parentId"`
	//排序
	Sort int `json:"sort"`
	//1支出2收入
	Outin int `json:"outin"`
	//用户id
	UserId int `json:"userid"`
	//子集
	Children []Classify `json:"children"`
}
