//数据库操作
package model

import (
	"bookkeeping/entity"
	"bookkeeping/util"
	"fmt"
	"strconv"
	"sync"
	"time"
)

/**
根据用户名和密码查找用户
*/
func GetUserNP(name string, password string) (*entity.User, error) {
	var sql = " SELECT ID,NAME,PASSWORD FROM USER WHERE NAME=? and PASSWORD=?"
	row := util.DB.QueryRow(sql, name, password)
	user := entity.NewUser()
	err := row.Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {

		return nil, err
	}
	return user, nil

}

/**
根据Id查找用户
*/
func GetUserId(id int) (*entity.User, error) {
	var sql = " SELECT ID,NAME,PASSWORD FROM USER WHERE ID=? "
	row := util.DB.QueryRow(sql, id)
	user := entity.NewUser()
	err := row.Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/**
保存用户信息
*/
func SaveUser(user *entity.User) error {
	var sql = " INSERT INTO  USER (NAME,PASSWORD)  VALUES(?,?)"
	_, err := util.DB.Exec(sql, user.Name, user.Password)
	if err != nil {
		return err
	}
	return nil
}

/**
保存类别信息
*/
func SaveClassifyDao(classify *entity.Classify) error {
	var sql = " INSERT INTO CLASSIFY (NAME,LAYER,PARENT_ID,SORT,OUTIN,USER_ID) VALUES(?,?,?,?,?,?) "
	_, err := util.DB.Exec(sql, classify.Name, classify.Layer, classify.ParentId, classify.Sort, classify.Outin, classify.UserId)
	if err != nil {
		return err
	}
	return nil
}

/**
查询分类信息
*/
func SelectClassifyDao(classify *entity.Classify) (*[]entity.Classify, error) {
	sql := "SELECT ID,NAME,LAYER,PARENT_ID,SORT,OUTIN FROM CLASSIFY WHERE USER_ID=?"
	var param []interface{} = make([]interface{}, 0, 1)
	param = append(param, classify.UserId)
	if classify.Id != 0 {
		sql += " AND ID=? "
		param = append(param, classify.Id)
	}
	if classify.Name != "" {
		sql += " AND Name=? "
		param = append(param, classify.Name)
	}
	if classify.Layer != 0 {
		sql += " AND Layer=? "
		param = append(param, classify.Layer)
	}
	if classify.ParentId != 0 {
		sql += " AND PARENT_ID=? "
		param = append(param, classify.ParentId)
	}
	sql += " ORDER BY SORT DESC"
	rows, err := util.DB.Query(sql, param...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var classifys []entity.Classify = make([]entity.Classify, 0, 5)
	var cfs entity.Classify
	for rows.Next() {
		err = rows.Scan(&cfs.Id, &cfs.Name, &cfs.Layer, &cfs.ParentId, &cfs.Sort, &cfs.Outin)
		if err != nil {
			return nil, err
		}
		classifys = append(classifys, cfs)
	}

	return &classifys, err
}

/**
根据用户id查询出分类后进行分组,将二级分类放入一级分类的children字段中然后分成收入和支出
*/
func ClassifyGroup(userId int) ([2][]entity.Classify, error) {
	//最后整理出的数据第一个为支出第二个为收入
	var result [2][]entity.Classify
	result[0] = make([]entity.Classify, 0)
	result[1] = make([]entity.Classify, 0)
	sql := "SELECT ID,NAME,LAYER,PARENT_ID,SORT,OUTIN FROM CLASSIFY WHERE USER_ID=?"
	rows, err := util.DB.Query(sql, userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	//一级分类
	var oneClassifys []entity.Classify = make([]entity.Classify, 0, 10)
	//二级分类
	var twoClassifys []entity.Classify = make([]entity.Classify, 0, 10)
	//分类对象用于接收数据库查出来的数据
	var cfs entity.Classify
	for rows.Next() {
		err = rows.Scan(&cfs.Id, &cfs.Name, &cfs.Layer, &cfs.ParentId, &cfs.Sort, &cfs.Outin)
		if err != nil {
			return result, err
		}
		//层级是否为一
		if cfs.Layer == 1 {
			oneClassifys = append(oneClassifys, cfs)
			//层级是否为二
		} else if cfs.Layer == 2 {
			twoClassifys = append(twoClassifys, cfs)
		}
	}

	//支出的管道
	var oneClassChan chan entity.Classify = make(chan entity.Classify, len(oneClassifys))
	//收入的管道
	var twoClassChan chan entity.Classify = make(chan entity.Classify, len(oneClassifys))
	wg := sync.WaitGroup{}
	//开启协成变量oneClassifys将数据分别放入支出/收入的管道
	wg.Add(3)
	go func() {
		for _, v := range oneClassifys {
			if v.Outin == 1 {
				oneClassChan <- v
			} else if v.Outin == 2 {
				twoClassChan <- v
			}
		}
		close(oneClassChan)
		close(twoClassChan)
		wg.Done()
	}()
	//开启协成整理支出
	go groupClassfiyChan(&oneClassChan, &twoClassifys, &result, 0, &wg)
	//开启协成整理收入
	go groupClassfiyChan(&twoClassChan, &twoClassifys, &result, 1, &wg)
	wg.Wait()

	return result, err
}

/**
*读取管道中的classify遍历二级分类放入children中整理完毕后放入根据传入的下标放入result结果中
*ch {*} classify的管道
*twoClassifys{*} 二级分类
*result {*}最后的结果
*i  result的下标
*wt 锁
 */
func groupClassfiyChan(ch *chan entity.Classify, twoClassifys *[]entity.Classify, result *[2][]entity.Classify, i int, wt *sync.WaitGroup) {
	for v := range *ch {
		v.Children = make([]entity.Classify, 0)
		for _, it := range *twoClassifys {
			if v.Id == it.ParentId {
				v.Children = append(v.Children, it)
			}
		}
		result[i] = append(result[i], v)
	}
	wt.Done()
}

/**
更新分类信息
*/
func UpdateClassify(classify *entity.Classify) error {
	var sql = " UPDATE CLASSIFY SET NAME=?,SORT=?,OUTIN=? WHERE ID = ? "
	_, err := util.DB.Exec(sql, classify.Name, classify.Sort, classify.Outin, classify.Id)
	if err != nil {
		return err
	}
	//如果为一级分类同时更新子集分类的支出/收入字段
	if classify.Layer == 1 {
		var sql = " UPDATE CLASSIFY SET OUTIN=? WHERE PARENT_ID = ? "
		_, err := util.DB.Exec(sql, classify.Outin, classify.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
删除分类
*/
func DelClassify(classify *entity.Classify) error {
	var sql = " DELETE FROM  CLASSIFY WHERE ID=? "
	_, err := util.DB.Exec(sql, classify.Id)
	if err != nil {
		return err
	}
	//如果为一级分类同时删除子集分类
	if classify.Layer == 1 {
		var sql = " DELETE FROM  CLASSIFY WHERE PARENT_ID = ? "
		_, err := util.DB.Exec(sql, classify.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
修改用户密码
*/
func AlterUserPassword(id int, password string) error {
	var sql = " UPDATE   USER SET PASSWORD=? WHERE ID=? "
	_, err := util.DB.Exec(sql, password, id)
	if err != nil {
		return err
	}
	return nil
}

/**
保存记录
*/
func InsertRecord(r *entity.Record) error {
	sql := " INSERT INTO  RECORD (MONERY,ONE_CLASSIFY,TWO_CLASSIFY,TIME,REMARK,OUTIN,USER_ID) VALUES(?,?,?,?,?,?,?) "
	_, err := util.DB.Exec(sql, r.Monery, r.OneClassify, r.TwoClassify, r.Time, r.Remark, r.Outin, r.UserId)
	return err
}

/**
查询记录
r record查询条件
page 	分页
pageSize 每页数量
*/
func SelectRecord(userId int, page int, pageSize int) ([]entity.Record, error) {
	sql := "SELECT ID,MONERY,ONE_CLASSIFY,TWO_CLASSIFY,TIME,REMARK,OUTIN,USER_ID FROM RECORD WHERE USER_ID=? ORDER BY TIME DESC"
	//var param []interface{} = make([]interface{}, 1)
	//param[0] = r.UserId
	if pageSize != 0 {
		sql += " LIMIT " + strconv.Itoa(pageSize*(page-1)) + "," + strconv.Itoa(pageSize)
	}
	rows, err := util.DB.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	var result []entity.Record = make([]entity.Record, 0)
	ite := entity.Record{}
	for rows.Next() {
		rows.Scan(&ite.Id, &ite.Monery, &ite.OneClassify, &ite.TwoClassify, &ite.Time, &ite.Remark, &ite.Outin, &ite.UserId)
		result = append(result, ite)
	}
	rows.Close()

	return result, err
}

/**
根据id查询记录
*/
func SelectRecordById(userId int, id int) (*entity.Record, error) {
	sql := "SELECT ID,MONERY,ONE_CLASSIFY,TWO_CLASSIFY,TIME,REMARK,OUTIN,USER_ID FROM RECORD WHERE USER_ID=? AND ID=?"
	row := util.DB.QueryRow(sql, userId, id)
	ite := entity.Record{}
	err := row.Scan(&ite.Id, &ite.Monery, &ite.OneClassify, &ite.TwoClassify, &ite.Time, &ite.Remark, &ite.Outin, &ite.UserId)
	return &ite, err
}

/**
根据id删除记录
*/
func DeleteRecordById(id int) error {
	sql := "DELETE FROM RECORD WHERE  ID=?"
	_, err := util.DB.Exec(sql, id)
	return err
}

/**
修改记录
*/
func UpdateRecord(r *entity.Record) error {
	sql := "UPDATE RECORD SET  MONERY=?,ONE_CLASSIFY=?,TWO_CLASSIFY=?,TIME=?,REMARK=?,OUTIN=? WHERE ID=?"
	_, err := util.DB.Exec(sql, r.Monery, r.OneClassify, r.TwoClassify, r.Time, r.Remark, r.Outin, r.Id)
	return err
}

//按时间分组统计记录
func TimeSelectRecord(timeType int, page int, pageSize int, userid int) (*[]entity.Record, error) {
	var tiemTemplate string
	switch timeType {
	case 1: //按日分组
		tiemTemplate = "%Y-%m-%d"
	case 2: //按月分组
		tiemTemplate = "%Y-%m"
	case 3: //按年分组
		tiemTemplate = "%Y"
	}
	sql := "SELECT SUM(MONERY) ,OUTIN ,date_format(TIME,'" + tiemTemplate + "') FROM " +
		" record WHERE USER_ID=? GROUP BY OUTIN,date_format(TIME,'" + tiemTemplate + "')  ORDER BY TIME DESC " +
		" LIMIT " + strconv.Itoa(page*pageSize) + "," + strconv.Itoa(pageSize)
	rows, err := util.DB.Query(sql, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []entity.Record = make([]entity.Record, 0)
	var record = entity.Record{}
	for rows.Next() {
		rows.Scan(&record.Monery, &record.Outin, &record.Time)
		result = append(result, record)
	}
	return &result, err
}

//按类别分组统计记录
func ClassifySelectRecord(classifyType int, page int, pageSize int, userid int) (*[]entity.Record, error) {
	var classTemplate string
	switch classifyType {
	case 1: //按日分组
		classTemplate = "ONE_CLASSIFY"
	case 2: //按月分组
		classTemplate = "TWO_CLASSIFY"
	}
	sql := "SELECT SUM(T1.MONERY),T1.OUTIN,T2.NAME,T3.NAME FROM " +
		" record T1 LEFT JOIN classify T2 ON T2.ID=T1.ONE_CLASSIFY LEFT JOIN classify T3 ON T3.ID=T1.TWO_CLASSIFY WHERE T1.USER_ID=? GROUP BY T1.OUTIN,T1." + classTemplate +
		" LIMIT " + strconv.Itoa(page*pageSize) + "," + strconv.Itoa(pageSize)
	rows, err := util.DB.Query(sql, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []entity.Record = make([]entity.Record, 0)
	var record = entity.Record{}
	for rows.Next() {
		rows.Scan(&record.Monery, &record.Outin, &record.OneClassifyName, &record.TwoClassifyName)
		result = append(result, record)
	}
	return &result, err
}

//根据传入的参数查询本日,本月,本年
func NowTimeSelectRecord(timeType int, userid int) (*[]entity.Record, error) {
	var tiemTemplate string
	year, month, day := time.Now().Date()
	var dataTime string
	switch timeType {
	case 1: //按日分组
		tiemTemplate = "%Y-%m-%d"
		dataTime = strconv.Itoa(year) + "-"
		if month < 10 {
			dataTime += "0"
		}
		dataTime += strconv.Itoa(int(month)) + "-"
		if day < 10 {
			dataTime += "0"
		}
		dataTime += strconv.Itoa(day)
	case 2: //按月分组
		tiemTemplate = "%Y-%m"
		dataTime = strconv.Itoa(year) + "-"
		if month < 10 {
			dataTime += "0"
		}
		dataTime += strconv.Itoa(int(month))
	case 3: //按年分组
		dataTime = fmt.Sprintf("%d", year)
		tiemTemplate = "%Y"
	}
	sql := ` SELECT SUM(MONERY),OUTIN FROM record WHERE USER_ID=? AND date_format(TIME,?)=?
	 GROUP BY OUTIN `
	rows, err := util.DB.Query(sql, userid, tiemTemplate, dataTime)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var result []entity.Record = make([]entity.Record, 0)
	var record = entity.Record{}
	for rows.Next() {
		rows.Scan(&record.Monery, &record.Outin)
		result = append(result, record)
	}
	return &result, err
}
