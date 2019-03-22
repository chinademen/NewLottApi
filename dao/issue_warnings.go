package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type IssueWarnings struct {
	Id               int       `orm:"column(id);auto"`
	CodecenterId     int8      `orm:"column(codecenter_id)"`
	LotteryId        int8      `orm:"column(lottery_id)"`
	Issue            string    `orm:"column(issue);size(25)"`
	Number           string    `orm:"column(number);size(60);null"`
	WarningType      string    `orm:"column(warning_type);size(6)" description:"告警类型"`
	ErrCode          string    `orm:"column(err_code);size(6);null"`
	ErrMsg           string    `orm:"column(err_msg);size(300);null"`
	EarliestDrawTime time.Time `orm:"column(earliest_draw_time);type(datetime);null"`
	RecordId         int       `orm:"column(record_id)"`
	CorrectTime      time.Time `orm:"column(correct_time);type(datetime);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);null"`
	Status           int8      `orm:"column(status);null" description:"0:未处理,1:已处理"`
}

func (t *IssueWarnings) TableName() string {
	return "issue_warnings"
}

func init() {
	orm.RegisterModel(new(IssueWarnings))
}

// AddIssueWarnings insert a new IssueWarnings into database and returns
// last inserted Id on success.
func AddIssueWarnings(m *IssueWarnings) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiIssueWarnings mutil insert a new IssueWarnings into database
func AddMultiIssueWarnings(mlist []*IssueWarnings) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetIssueWarningsById retrieves IssueWarnings by Id. Returns error if
// Id doesn't exist
func GetIssueWarningsById(id int) (v *IssueWarnings, err error) {
	o := orm.NewOrm()
	v = &IssueWarnings{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllIssueWarnings retrieves all IssueWarnings matches certain condition. Returns empty list if
// no records exist
func GetAllIssueWarnings(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []IssueWarnings, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(IssueWarnings))
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

// UpdateIssueWarnings updates IssueWarnings by Id and returns error if
// the record to be updated doesn't exist
func UpdateIssueWarningsById(m *IssueWarnings) (err error) {
	o := orm.NewOrm()
	v := IssueWarnings{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteIssueWarnings deletes IssueWarnings by Id and returns error if
// the record to be deleted doesn't exist
func DeleteIssueWarnings(id int) (err error) {
	o := orm.NewOrm()
	v := IssueWarnings{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&IssueWarnings{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
