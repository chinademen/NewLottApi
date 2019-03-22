package models

import (
	"fmt"
	"strconv"
)

/**
* 查询一行数据指定的字段
* @param		sWhere		查询条件，必填不要带where
* @param		sField		查询字段，必填
* return		返回map
 */
func GetOne(tableName, sWhere, sField string) map[string]string {
	res := map[string]string{}

	//如果查询条件和字段没有填写，则不允许查询
	if len(sWhere) < 1 || len(sField) < 1 {
		return res
	}
	q := fmt.Sprintf("select %s from %s where %s limit 1;", sField, tableName, sWhere)
	res = Mdb.GetRow(q)
	return res
}

/**
* 查询多行数据指定的字段,支持数据分页
* @param		sWhere		查询条件，必填不要带where
* @param		sField		查询字段，必填
* @param		sOrder		排序字段，必填
* @param		offset		数据集游标
* @param		limit		查询记录的行数,必填
* return		返回数组map
 */
func GetList(tableName, sWhere, sField, sOrder string, offset, limit int) []map[string]string {
	res := []map[string]string{}
	//如果查询条件和字段没有填写，则不允许查询
	if len(sWhere) < 1 || len(sField) < 1 {
		return res
	}
	if len(sOrder) > 0 {
		sOrder = " order by " + sOrder
	}
	q := fmt.Sprintf("select %s from %s where %s %s limit %d,%d;", sField, tableName, sWhere, sOrder, offset, limit)
	res = Mdb.GetRows(q)
	return res
}

/**
* 查询多行数据指定的字段,通过分组
* @param		sWhere		查询条件，必填不要带where
* @param		sGroup		分组条件，选填
* @param		sField		查询字段，必填
* @param		sOrder		排序字段，必填
* @param		offset		数据集游标
* @param		limit		查询记录的行数,必填
* return		返回数组map
 */
func GetListByGroup(tableName, sWhere, sGroup, sField, sOrder string, offset, limit int) []map[string]string {
	res := []map[string]string{}
	//如果查询条件和字段没有填写，则不允许查询
	if len(sWhere) < 1 || len(sField) < 1 {
		return res
	}
	if len(sGroup) > 0 {
		sGroup = " group by " + sGroup
	}
	if len(sOrder) > 0 {
		sOrder = " order by " + sOrder
	}
	q := fmt.Sprintf("select %s from %s where %s %s %s limit %d,%d;", sField, tableName, sWhere, sGroup, sOrder, offset, limit)

	res = Mdb.GetRows(q)
	return res
}

/**
* 查询符合条件的条数
* @param		sWhere		查询条件，必填不要带where
* return		返回int
 */
func GetCount(tableName, sWhere string) int {
	res := 0
	q := fmt.Sprintf("select count(0) as nums from %s where %s;", tableName, sWhere)
	row := Mdb.GetRow(q)
	if len(row) > 0 {
		tmp := row["nums"]
		res, _ = strconv.Atoi(tmp)
	}
	return res
}

/**
* 查询符合条件的和值
* @param		sWhere		查询条件，必填不要带where
* @param		field		需要字段
* return		返回float64
 */
func GetSum(tableName, sWhere, field string) float64 {
	res := 0.00
	q := fmt.Sprintf("select sum(%s) as sum from %s where %s;", field, tableName, sWhere)
	row := Mdb.GetRow(q)
	if len(row) > 0 {
		tmp := row["sum"]
		res, _ = strconv.ParseFloat(tmp, 64)
	}
	return res
}

/**
* 查询符合条件的平均值
* @param		sWhere		查询条件，必填不要带where
* @param		field		需要字段
* return		返回float64
 */
func GetAvg(tableName, sWhere, field string) float64 {
	res := 0.00
	q := fmt.Sprintf("select avg(%s) as avg from %s where %s;", field, tableName, sWhere)
	row := Mdb.GetRow(q)
	if len(row) > 0 {
		tmp := row["avg"]
		res, _ = strconv.ParseFloat(tmp, 64)
	}
	return res
}

/**
* 查询符合条件的最大值
* @param		sWhere		查询条件，必填不要带where
* @param		field		需要字段
* return		返回float64
 */
func GetMax(tableName, sWhere, field string) float64 {
	res := 0.00
	q := fmt.Sprintf("select max(%s) as max from %s where %s;", field, tableName, sWhere)
	row := Mdb.GetRow(q)
	if len(row) > 0 {
		tmp := row["max"]
		res, _ = strconv.ParseFloat(tmp, 64)
	}
	return res
}

/**
* 查询符合条件的最大值
* @param		sWhere		查询条件，必填不要带where
* @param		field		需要字段
* return		返回float64
 */
func GetMin(tableName, sWhere, field string) float64 {
	res := 0.00
	q := fmt.Sprintf("select min(%s) as min from %s where %s;", field, tableName, sWhere)
	row := Mdb.GetRow(q)
	if len(row) > 0 {
		tmp := row["min"]
		res, _ = strconv.ParseFloat(tmp, 64)
	}
	return res
}

/**
* 更新内容
* @param		mData		需要更新的map[string]string
* @param		sWhere		查询条件，必填不要带where
* return		int	返回受影响的行数
 */
func Update(mData map[string]string, tableName, sWhere string) int {
	row := Mdb.UpdateRow(tableName, mData, sWhere)

	return row
}

/**
* 插入数据
* return		int	返回受影响的行数
* return		id	返回id，主要用于集群生成的id
 */
func Insert(mData map[string]string, tableName string) (int, string) {
	row, id := Mdb.InsertRow(tableName, mData)
	return row, id
}

/**
* 删除数据
* @param		sWhere		查询条件，必填不要带where
* return		int	返回受影响的行数
 */
func Delete(tableName, sWhere string) int {
	row := Mdb.DeleteRowsByWhere(tableName, sWhere)
	return row
}

/*** 将插入数据还原成sql语句
* @param		mData
 */
func GetInsertSql(tableName string, mData map[string]string) string {
	sql := Mdb.GetInsertSql(tableName, mData)
	return sql
}

/**
* 将更新数据还原成sql语句
 */
func GetUpdateSql(tableName string, mData map[string]string, w string) string {
	sql := Mdb.GetUpdateSql(tableName, mData, w)
	return sql
}

/**
* 事务处理
* @param			sqls			多组sql语句
 */

func Transaction(sqls []string) error {
	err := Mdb.CommitSql(sqls)
	return err
}

/*** 将插入数据还原成sql语句 有重复id报错
* @param		mData
 */
func GetInsertTrueSql(tableName string, mData map[string]string) string {
	sql := Mdb.GetInsertTrueSql(tableName, mData)
	return sql
}
