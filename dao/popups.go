package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Popups struct {
	Id        int       `orm:"column(id);auto"`
	Name      string    `orm:"column(name);size(32)"`
	NeedForm  uint8     `orm:"column(need_form)"`
	Method    uint8     `orm:"column(method)"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *Popups) TableName() string {
	return "popups"
}

func init() {
	orm.RegisterModel(new(Popups))
}

// AddPopups insert a new Popups into database and returns
// last inserted Id on success.
func AddPopups(m *Popups) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPopups mutil insert a new Popups into database
func AddMultiPopups(mlist []*Popups) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPopupsById retrieves Popups by Id. Returns error if
// Id doesn't exist
func GetPopupsById(id int) (v *Popups, err error) {
	o := orm.NewOrm()
	v = &Popups{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPopups retrieves all Popups matches certain condition. Returns empty list if
// no records exist
func GetAllPopups(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Popups, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Popups))
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

// UpdatePopups updates Popups by Id and returns error if
// the record to be updated doesn't exist
func UpdatePopupsById(m *Popups) (err error) {
	o := orm.NewOrm()
	v := Popups{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePopups deletes Popups by Id and returns error if
// the record to be deleted doesn't exist
func DeletePopups(id int) (err error) {
	o := orm.NewOrm()
	v := Popups{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Popups{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
