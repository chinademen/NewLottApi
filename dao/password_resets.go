package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PasswordResets struct {
	Id        int       `orm:"column(id);auto"`
	Username  string    `orm:"column(username);size(255)"`
	Token     string    `orm:"column(token);size(255)"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);auto_now"`
}

func (t *PasswordResets) TableName() string {
	return "password_resets"
}

func init() {
	orm.RegisterModel(new(PasswordResets))
}

// AddPasswordResets insert a new PasswordResets into database and returns
// last inserted Id on success.
func AddPasswordResets(m *PasswordResets) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPasswordResets mutil insert a new PasswordResets into database
func AddMultiPasswordResets(mlist []*PasswordResets) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPasswordResetsById retrieves PasswordResets by Id. Returns error if
// Id doesn't exist
func GetPasswordResetsById(id int) (v *PasswordResets, err error) {
	o := orm.NewOrm()
	v = &PasswordResets{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPasswordResets retrieves all PasswordResets matches certain condition. Returns empty list if
// no records exist
func GetAllPasswordResets(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []PasswordResets, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PasswordResets))
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

// UpdatePasswordResets updates PasswordResets by Id and returns error if
// the record to be updated doesn't exist
func UpdatePasswordResetsById(m *PasswordResets) (err error) {
	o := orm.NewOrm()
	v := PasswordResets{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePasswordResets deletes PasswordResets by Id and returns error if
// the record to be deleted doesn't exist
func DeletePasswordResets(id int) (err error) {
	o := orm.NewOrm()
	v := PasswordResets{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PasswordResets{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
