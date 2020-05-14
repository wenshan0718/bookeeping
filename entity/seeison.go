package entity

import (
	"bookkeeping/util"
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

//session 类型
type Session struct {
	SessionId string
	UserID    int
	Name      string
}

/**
通过sessionID 获取redis中的session数据
*/
func GetSessionByCookie(sessionId string) (*Session, error) {
	conn := util.Pool.Get()
	defer conn.Close()
	mes, err := redis.String(conn.Do("get", sessionId))
	if err != nil {
		if err == redis.ErrNil {
			err = REDIS_SELECT_RESULT_NIL
		}
		return nil, err
	}
	nowTime := time.Now().Unix()
	//跟新过期时间10分钟
	_, err = conn.Do("EXPIREAT", sessionId, nowTime+60*10)
	if err != nil {
		return nil, err
	}
	//将数据反序列化成session类型
	var session Session
	err = json.Unmarshal([]byte(mes), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil

}

/**
将sesison保存到redis中
*/
func (session *Session) SaveSession() error {
	//将session数据序列化转成字符串然后保存到redis中
	byte := make([]byte, 1024)
	byte, err := json.Marshal(session)
	if err != nil {
		return err
	}
	conn := util.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("set", session.SessionId, string(byte))
	if err != nil {
		return err
	}

	nowTime := time.Now().Unix()
	//设置过期时间10分钟
	_, err = conn.Do("EXPIREAT", session.SessionId, nowTime+60*10)
	if err != nil {
		return err
	}
	return nil
}
