package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type LotteryWayBlackList struct {
	Id         int   `orm:"column(id);auto" json:"id,string"`
	LotteryId  int32 `orm:"column(lottery_id)" json:"lottery_id,string"`
	SeriesWay  int32 `orm:"column(series_way)" json:"series_way,string"`
	TerminalId uint8 `orm:"column(terminal_id);null" json:"terminal_id,string"`
}

func (t *LotteryWayBlackList) TableName() string {
	return "lottery_way_black_list"
}

func init() {
	orm.RegisterModel(new(LotteryWayBlackList))
}

// AddLotteryWayBlackList insert a new LotteryWayBlackList into database and returns
// last inserted Id on success.
func AddLotteryWayBlackList(m *LotteryWayBlackList) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiLotteryWayBlackList mutil insert a new LotteryWayBlackList into database
func AddMultiLotteryWayBlackList(mlist []*LotteryWayBlackList) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetLotteryWayBlackListById retrieves LotteryWayBlackList by Id. Returns error if
// Id doesn't exist
func GetLotteryWayBlackListById(id int) (v *LotteryWayBlackList, err error) {
	o := orm.NewOrm()
	v = &LotteryWayBlackList{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLotteryWayBlackList retrieves all LotteryWayBlackList matches certain condition. Returns empty list if
// no records exist
func GetAllLotteryWayBlackList(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []LotteryWayBlackList, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(LotteryWayBlackList))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&ml, fields...); err == nil {
		return ml, nil
	}
	return nil, err
}

// UpdateLotteryWayBlackList updates LotteryWayBlackList by Id and returns error if
// the record to be updated doesn't exist
func UpdateLotteryWayBlackListById(m *LotteryWayBlackList) (err error) {
	o := orm.NewOrm()
	v := LotteryWayBlackList{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLotteryWayBlackList deletes LotteryWayBlackList by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLotteryWayBlackList(id int) (err error) {
	o := orm.NewOrm()
	v := LotteryWayBlackList{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&LotteryWayBlackList{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
