package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type CodeCenters struct {
	Id              int       `orm:"column(id);auto"`
	Name            string    `orm:"column(name);size(20)"`
	CustomerId      uint      `orm:"column(customer_id);null"`
	Version         uint8     `orm:"column(version)"`
	Domain          string    `orm:"column(domain);size(50)"`
	Ip              string    `orm:"column(ip);size(100)"`
	SetUrl          string    `orm:"column(set_url);size(200)"`
	SetVerifyUrl    string    `orm:"column(set_verify_url);size(200)"`
	GetUrl          string    `orm:"column(get_url);size(200)"`
	ReviseUrl       string    `orm:"column(revise_url);size(200);null"`
	ReviseVerifyUrl string    `orm:"column(revise_verify_url);size(200);null"`
	AlarmUrl        string    `orm:"column(alarm_url);size(200);null"`
	AlarmVerifyUrl  string    `orm:"column(alarm_verify_url);size(200);null"`
	CustomerKey     string    `orm:"column(customer_key);size(32)"`
	CreatedAt       time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *CodeCenters) TableName() string {
	return "code_centers"
}

func init() {
	orm.RegisterModel(new(CodeCenters))
}

// AddCodeCenters insert a new CodeCenters into database and returns
// last inserted Id on success.
func AddCodeCenters(m *CodeCenters) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiCodeCenters mutil insert a new CodeCenters into database
func AddMultiCodeCenters(mlist []*CodeCenters) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetCodeCentersById retrieves CodeCenters by Id. Returns error if
// Id doesn't exist
func GetCodeCentersById(id int) (v *CodeCenters, err error) {
	o := orm.NewOrm()
	v = &CodeCenters{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCodeCenters retrieves all CodeCenters matches certain condition. Returns empty list if
// no records exist
func GetAllCodeCenters(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []CodeCenters, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CodeCenters))
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

// UpdateCodeCenters updates CodeCenters by Id and returns error if
// the record to be updated doesn't exist
func UpdateCodeCentersById(m *CodeCenters) (err error) {
	o := orm.NewOrm()
	v := CodeCenters{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCodeCenters deletes CodeCenters by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCodeCenters(id int) (err error) {
	o := orm.NewOrm()
	v := CodeCenters{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CodeCenters{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
