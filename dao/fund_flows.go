package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type FundFlows struct {
	Id          int       `orm:"column(id);auto"`
	Description string    `orm:"column(description);size(50)"`
	Balance     int8      `orm:"column(balance)"`
	Available   int8      `orm:"column(available)"`
	Frozen      int8      `orm:"column(frozen)"`
	DeletedAt   time.Time `orm:"column(deleted_at);type(datetime);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *FundFlows) TableName() string {
	return "fund_flows"
}

func init() {
	orm.RegisterModel(new(FundFlows))
}

// AddFundFlows insert a new FundFlows into database and returns
// last inserted Id on success.
func AddFundFlows(m *FundFlows) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiFundFlows mutil insert a new FundFlows into database
func AddMultiFundFlows(mlist []*FundFlows) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetFundFlowsById retrieves FundFlows by Id. Returns error if
// Id doesn't exist
func GetFundFlowsById(id int) (v *FundFlows, err error) {
	o := orm.NewOrm()
	v = &FundFlows{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFundFlows retrieves all FundFlows matches certain condition. Returns empty list if
// no records exist
func GetAllFundFlows(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []FundFlows, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(FundFlows))
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

// UpdateFundFlows updates FundFlows by Id and returns error if
// the record to be updated doesn't exist
func UpdateFundFlowsById(m *FundFlows) (err error) {
	o := orm.NewOrm()
	v := FundFlows{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFundFlows deletes FundFlows by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFundFlows(id int) (err error) {
	o := orm.NewOrm()
	v := FundFlows{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&FundFlows{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
