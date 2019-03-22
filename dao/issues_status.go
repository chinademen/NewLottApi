package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type IssuesStatus struct {
	Id        int    `orm:"column(id);auto"`
	Type      int8   `orm:"column(type);null" description:"状态类型(计奖/派奖/追号)"`
	IssuesId  int    `orm:"column(issues_id);null" description:"奖期id"`
	LotteryId uint8  `orm:"column(lottery_id)" description:"彩种id"`
	Issue     string `orm:"column(issue);size(15)" description:"奖期"`
	Status    int8   `orm:"column(status);null" description:"状态值，1=完成"`
}

func (t *IssuesStatus) TableName() string {
	return "issues_status"
}

func init() {
	orm.RegisterModel(new(IssuesStatus))
}

// AddIssuesStatus insert a new IssuesStatus into database and returns
// last inserted Id on success.
func AddIssuesStatus(m *IssuesStatus) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiIssuesStatus mutil insert a new IssuesStatus into database
func AddMultiIssuesStatus(mlist []*IssuesStatus) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetIssuesStatusById retrieves IssuesStatus by Id. Returns error if
// Id doesn't exist
func GetIssuesStatusById(id int) (v *IssuesStatus, err error) {
	o := orm.NewOrm()
	v = &IssuesStatus{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllIssuesStatus retrieves all IssuesStatus matches certain condition. Returns empty list if
// no records exist
func GetAllIssuesStatus(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []IssuesStatus, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(IssuesStatus))
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

// UpdateIssuesStatus updates IssuesStatus by Id and returns error if
// the record to be updated doesn't exist
func UpdateIssuesStatusById(m *IssuesStatus) (err error) {
	o := orm.NewOrm()
	v := IssuesStatus{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteIssuesStatus deletes IssuesStatus by Id and returns error if
// the record to be deleted doesn't exist
func DeleteIssuesStatus(id int) (err error) {
	o := orm.NewOrm()
	v := IssuesStatus{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&IssuesStatus{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
