package entity

import "errors"

//自定义错误

var (
	REDIS_SELECT_RESULT_NIL = errors.New("redis查询数据为空")
)
