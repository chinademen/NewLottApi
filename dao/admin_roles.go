package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminRoles struct {
	Id            int       `orm:"column(id);auto" description:"分组ID"`
	Name          string    `orm:"column(name);size(40)"`
	Description   string    `orm:"column(description);size(255);null"`
	Rights        string    `orm:"column(rights);size(10240);null"`
	Priority      uint      `orm:"column(priority)"`
	IsSystem      int8      `orm:"column(is_system)" description:"系统角色"`
	RightSettable int8      `orm:"column(right_settable)" description:"是否可设置权限"`
	UserSettable  int8      `orm:"column(user_settable)" description:"是否可设置用户"`
	Disabled      uint8     `orm:"column(disabled)" description:"是否禁用"`
	Sequence      uint      `orm:"column(sequence)" description:"排序值"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *AdminRoles) TableName() string {
	return "admin_roles"
}

func init() {
	orm.RegisterModel(new(AdminRoles))
}

// AddAdminRoles insert a new AdminRoles into database and returns
// last inserted Id on success.
func AddAdminRoles(m *AdminRoles) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdminRoles mutil insert a new AdminRoles into database
func AddMultiAdminRoles(mlist []*AdminRoles) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdminRolesById retrieves AdminRoles by Id. Returns error if
// Id doesn't exist
func GetAdminRolesById(id int) (v *AdminRoles, err error) {
	o := orm.NewOrm()
	v = &AdminRoles{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdminRoles retrieves all AdminRoles matches certain condition. Returns empty list if
// no records exist
func GetAllAdminRoles(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdminRoles, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdminRoles))
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

// UpdateAdminRoles updates AdminRoles by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdminRolesById(m *AdminRoles) (err error) {
	o := orm.NewOrm()
	v := AdminRoles{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdminRoles deletes AdminRoles by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdminRoles(id int) (err error) {
	o := orm.NewOrm()
	v := AdminRoles{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdminRoles{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
