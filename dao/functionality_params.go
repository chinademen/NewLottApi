package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type FunctionalityParams struct {
	Id              int       `orm:"column(id);auto"`
	FunctionalityId uint      `orm:"column(functionality_id)"`
	Name            string    `orm:"column(name);size(32)"`
	Type            string    `orm:"column(type);size(20)" description:"数据类型"`
	DefaultValue    string    `orm:"column(default_value);size(200);null"`
	LimitWhenNull   int8      `orm:"column(limit_when_null)" description:"为空时是否设定"`
	Sequence        uint      `orm:"column(sequence);null" description:"排序"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *FunctionalityParams) TableName() string {
	return "functionality_params"
}

func init() {
	orm.RegisterModel(new(FunctionalityParams))
}

// AddFunctionalityParams insert a new FunctionalityParams into database and returns
// last inserted Id on success.
func AddFunctionalityParams(m *FunctionalityParams) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiFunctionalityParams mutil insert a new FunctionalityParams into database
func AddMultiFunctionalityParams(mlist []*FunctionalityParams) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetFunctionalityParamsById retrieves FunctionalityParams by Id. Returns error if
// Id doesn't exist
func GetFunctionalityParamsById(id int) (v *FunctionalityParams, err error) {
	o := orm.NewOrm()
	v = &FunctionalityParams{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFunctionalityParams retrieves all FunctionalityParams matches certain condition. Returns empty list if
// no records exist
func GetAllFunctionalityParams(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []FunctionalityParams, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(FunctionalityParams))
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

// UpdateFunctionalityParams updates FunctionalityParams by Id and returns error if
// the record to be updated doesn't exist
func UpdateFunctionalityParamsById(m *FunctionalityParams) (err error) {
	o := orm.NewOrm()
	v := FunctionalityParams{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFunctionalityParams deletes FunctionalityParams by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFunctionalityParams(id int) (err error) {
	o := orm.NewOrm()
	v := FunctionalityParams{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&FunctionalityParams{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
