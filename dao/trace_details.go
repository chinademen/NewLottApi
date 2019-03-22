package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type TraceDetails struct {
	Id         int       `orm:"column(id);auto" json:"id,string"`
	MerchantId int       `orm:"column(merchant_id)" json:"merchant_id,string" description:"商户"`
	UserId     int64     `orm:"column(user_id)" json:"user_id,string" description:"用户id"`
	AccountId  int64     `orm:"column(account_id)" json:"account_id,string" description:"用户账户id"`
	TraceId    uint64    `orm:"column(trace_id)" json:"trace_id,string" description:"追号id"`
	LotteryId  uint8     `orm:"column(lottery_id)" json:"lottery_id,string" description:"彩种id"`
	Issue      string    `orm:"column(issue);size(15)" json:"issue,string" description:"奖期"`
	EndTime    int       `orm:"column(end_time);size(16);null" json:"end_time,string" description:"截至时间"`
	Multiple   string    `orm:"column(multiple);size(16)" json:"multiple,string" description:"倍数"`
	Amount     float64   `orm:"column(amount);digits(14);decimals(4)" json:"amount,string" description:"投注金额"`
	ProjectId  uint64    `orm:"column(project_id);null" json:"project_id,string" description:"注单id"`
	Status     int8      `orm:"column(status)" json:"status,string" description:"状态"`
	BoughtAt   string    `orm:"column(bought_at);type(datetime);null" json:"bought_at,string" description:"投注时间"`
	CanceledAt time.Time `orm:"column(canceled_at);type(datetime);null" json:"canceled_at,string" description:"取消时间"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *TraceDetails) TableName() string {
	return "trace_details"
}

func init() {
	orm.RegisterModel(new(TraceDetails))
}

// AddTraceDetails insert a new TraceDetails into database and returns
// last inserted Id on success.
func AddTraceDetails(o orm.Ormer, m *TraceDetails) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// AddMultiTraceDetails mutil insert a new TraceDetails into database
func AddMultiTraceDetails(mlist []*TraceDetails) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetTraceDetailsById retrieves TraceDetails by Id. Returns error if
// Id doesn't exist
func GetTraceDetailsById(id int) (v *TraceDetails, err error) {
	o := orm.NewOrm()
	v = &TraceDetails{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTraceDetails retrieves all TraceDetails matches certain condition. Returns empty list if
// no records exist
func GetAllTraceDetails(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*TraceDetails, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(TraceDetails))
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

// UpdateTraceDetails updates TraceDetails by Id and returns error if
// the record to be updated doesn't exist
func UpdateTraceDetailsById(o orm.Ormer, m *TraceDetails) (err error) {
	v := TraceDetails{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTraceDetails deletes TraceDetails by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTraceDetails(id int) (err error) {
	o := orm.NewOrm()
	v := TraceDetails{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TraceDetails{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
