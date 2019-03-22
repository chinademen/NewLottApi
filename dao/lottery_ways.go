package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type LotteryWays struct {
	Id          int       `orm:"column(id);auto"`
	SeriesId    uint8     `orm:"column(series_id)" description:"系列ID"`
	LotteryId   uint32    `orm:"column(lottery_id)" description:"彩种ＩＤ"`
	SeriesWayId uint      `orm:"column(series_way_id)" description:"系列投注方式"`
	Name        string    `orm:"column(name);size(30)"`
	ShortName   string    `orm:"column(short_name);size(30);null"`
	Status      int8      `orm:"column(status);null" description:"状态 0：正常  1：关闭"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *LotteryWays) TableName() string {
	return "lottery_ways"
}

func init() {
	orm.RegisterModel(new(LotteryWays))
}

// AddLotteryWays insert a new LotteryWays into database and returns
// last inserted Id on success.
func AddLotteryWays(m *LotteryWays) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiLotteryWays mutil insert a new LotteryWays into database
func AddMultiLotteryWays(mlist []*LotteryWays) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetLotteryWaysById retrieves LotteryWays by Id. Returns error if
// Id doesn't exist
func GetLotteryWaysById(id int) (v *LotteryWays, err error) {
	o := orm.NewOrm()
	v = &LotteryWays{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLotteryWays retrieves all LotteryWays matches certain condition. Returns empty list if
// no records exist
func GetAllLotteryWays(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []LotteryWays, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(LotteryWays))
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

// UpdateLotteryWays updates LotteryWays by Id and returns error if
// the record to be updated doesn't exist
func UpdateLotteryWaysById(m *LotteryWays) (err error) {
	o := orm.NewOrm()
	v := LotteryWays{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLotteryWays deletes LotteryWays by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLotteryWays(id int) (err error) {
	o := orm.NewOrm()
	v := LotteryWays{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&LotteryWays{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
