package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type InstantIssues struct {
	Id        int       `orm:"column(id);auto"`
	UserId    uint      `orm:"column(user_id)" description:"用户ID"`
	LotteryId uint8     `orm:"column(lottery_id)" description:"彩种"`
	Issue     string    `orm:"column(issue);size(15)" description:"奖期"`
	WnNumber  string    `orm:"column(wn_number);size(60)" description:"中奖号码"`
	EncodedAt time.Time `orm:"column(encoded_at);type(datetime);null" description:"录号时间"`
	Status    uint8     `orm:"column(status)" description:"状态"`
	Tag       string    `orm:"column(tag);size(50);null" description:"附加信息"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *InstantIssues) TableName() string {
	return "instant_issues"
}

func init() {
	orm.RegisterModel(new(InstantIssues))
}

// AddInstantIssues insert a new InstantIssues into database and returns
// last inserted Id on success.
func AddInstantIssues(m *InstantIssues) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiInstantIssues mutil insert a new InstantIssues into database
func AddMultiInstantIssues(mlist []*InstantIssues) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetInstantIssuesById retrieves InstantIssues by Id. Returns error if
// Id doesn't exist
func GetInstantIssuesById(id int) (v *InstantIssues, err error) {
	o := orm.NewOrm()
	v = &InstantIssues{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllInstantIssues retrieves all InstantIssues matches certain condition. Returns empty list if
// no records exist
func GetAllInstantIssues(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []InstantIssues, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(InstantIssues))
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

// UpdateInstantIssues updates InstantIssues by Id and returns error if
// the record to be updated doesn't exist
func UpdateInstantIssuesById(m *InstantIssues) (err error) {
	o := orm.NewOrm()
	v := InstantIssues{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteInstantIssues deletes InstantIssues by Id and returns error if
// the record to be deleted doesn't exist
func DeleteInstantIssues(id int) (err error) {
	o := orm.NewOrm()
	v := InstantIssues{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&InstantIssues{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
