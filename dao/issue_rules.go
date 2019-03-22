package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type IssueRules struct {
	Id              int       `orm:"column(id);auto"`
	LotteryId       uint32    `orm:"column(lottery_id)" description:"彩种ID"`
	BeginTime       time.Time `orm:"column(begin_time);type(time)" description:"开始时间"`
	EndTime         time.Time `orm:"column(end_time);type(time)" description:"截止时间"`
	Cycle           uint      `orm:"column(cycle);null" description:"周期"`
	NumberDelayTime uint      `orm:"column(number_delay_time)" description:"奖期开售延迟时间"`
	FirstTime       time.Time `orm:"column(first_time);type(time)" description:"首期截止时间"`
	StopAdjustTime  uint16    `orm:"column(stop_adjust_time)" description:"销售截止时间调整（提前多少秒）"`
	EncodeTime      uint16    `orm:"column(encode_time)" description:"录号延迟时间（官方开间时间之后多少妙可以录号）"`
	IssueCount      uint      `orm:"column(issue_count);null"`
	Enabled         uint8     `orm:"column(enabled)" description:"是否有效"`
	CreatedAt       time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *IssueRules) TableName() string {
	return "issue_rules"
}

func init() {
	orm.RegisterModel(new(IssueRules))
}

// AddIssueRules insert a new IssueRules into database and returns
// last inserted Id on success.
func AddIssueRules(m *IssueRules) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiIssueRules mutil insert a new IssueRules into database
func AddMultiIssueRules(mlist []*IssueRules) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetIssueRulesById retrieves IssueRules by Id. Returns error if
// Id doesn't exist
func GetIssueRulesById(id int) (v *IssueRules, err error) {
	o := orm.NewOrm()
	v = &IssueRules{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllIssueRules retrieves all IssueRules matches certain condition. Returns empty list if
// no records exist
func GetAllIssueRules(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []IssueRules, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(IssueRules))
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

// UpdateIssueRules updates IssueRules by Id and returns error if
// the record to be updated doesn't exist
func UpdateIssueRulesById(m *IssueRules) (err error) {
	o := orm.NewOrm()
	v := IssueRules{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteIssueRules deletes IssueRules by Id and returns error if
// the record to be deleted doesn't exist
func DeleteIssueRules(id int) (err error) {
	o := orm.NewOrm()
	v := IssueRules{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&IssueRules{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
