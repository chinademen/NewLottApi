package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type TerminalLotteries struct {
	Id         int       `orm:"column(id);auto" json:"id,string"`
	TerminalId uint      `orm:"column(terminal_id);null" json:"terminal_id,string" description:"终端id"`
	LotteryId  uint32    `orm:"column(lottery_id);null" json:"lottery_id,string" description:"彩种id"`
	Status     uint8     `orm:"column(status)" json:"status,string" description:"终端状态：1=正式;2=测试;0=关闭"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *TerminalLotteries) TableName() string {
	return "terminal_lotteries"
}

func init() {
	orm.RegisterModel(new(TerminalLotteries))
}

// AddTerminalLotteries insert a new TerminalLotteries into database and returns
// last inserted Id on success.
func AddTerminalLotteries(m *TerminalLotteries) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiTerminalLotteries mutil insert a new TerminalLotteries into database
func AddMultiTerminalLotteries(mlist []*TerminalLotteries) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetTerminalLotteriesById retrieves TerminalLotteries by Id. Returns error if
// Id doesn't exist
func GetTerminalLotteriesById(id int) (v *TerminalLotteries, err error) {
	o := orm.NewOrm()
	v = &TerminalLotteries{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTerminalLotteries retrieves all TerminalLotteries matches certain condition. Returns empty list if
// no records exist
func GetAllTerminalLotteries(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []TerminalLotteries, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(TerminalLotteries))
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

// UpdateTerminalLotteries updates TerminalLotteries by Id and returns error if
// the record to be updated doesn't exist
func UpdateTerminalLotteriesById(m *TerminalLotteries) (err error) {
	o := orm.NewOrm()
	v := TerminalLotteries{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTerminalLotteries deletes TerminalLotteries by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTerminalLotteries(id int) (err error) {
	o := orm.NewOrm()
	v := TerminalLotteries{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TerminalLotteries{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
