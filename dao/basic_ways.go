package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type BasicWays struct {
	Id          int       `orm:"column(id);auto" json:"id,string"`
	LotteryType uint8     `orm:"column(lottery_type)" json:"lottery_type,string" description:"彩票类型: 1-数字排列类型 2-乐透类型"`
	Name        string    `orm:"column(name);size(10)" json:"name,string" description:"名称"`
	Function    string    `orm:"column(function);size(64)" json:"function,string" description:"计奖方法"`
	Description string    `orm:"column(description);size(255);null" json:"description,string" description:"描述"`
	Sequence    uint      `orm:"column(sequence);null" json:"sequence,string" description:"排序"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *BasicWays) TableName() string {
	return "basic_ways"
}

func init() {
	orm.RegisterModel(new(BasicWays))
}

// AddBasicWays insert a new BasicWays into database and returns
// last inserted Id on success.
func AddBasicWays(m *BasicWays) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiBasicWays mutil insert a new BasicWays into database
func AddMultiBasicWays(mlist []*BasicWays) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetBasicWaysById retrieves BasicWays by Id. Returns error if
// Id doesn't exist
func GetBasicWaysById(id int) (v *BasicWays, err error) {
	o := orm.NewOrm()
	v = &BasicWays{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBasicWays retrieves all BasicWays matches certain condition. Returns empty list if
// no records exist
func GetAllBasicWays(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*BasicWays, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(BasicWays))
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

// UpdateBasicWays updates BasicWays by Id and returns error if
// the record to be updated doesn't exist
func UpdateBasicWaysById(m *BasicWays) (err error) {
	o := orm.NewOrm()
	v := BasicWays{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBasicWays deletes BasicWays by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBasicWays(id int) (err error) {
	o := orm.NewOrm()
	v := BasicWays{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BasicWays{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
