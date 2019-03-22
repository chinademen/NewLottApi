package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserLoginIps struct {
	Id         int       `orm:"column(id);auto"`
	MerchantId int       `orm:"column(merchant_id);null" description:"商户id"`
	UserId     int       `orm:"column(user_id);null" description:"用户id"`
	Username   string    `orm:"column(username);size(16);null" description:"用户名"`
	IsTester   int8      `orm:"column(is_tester);null" description:"1=测试"`
	Ip         string    `orm:"column(ip);size(15);null" description:"ip地址"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null" description:"创建时间"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null" description:"更新时间"`
}

func (t *UserLoginIps) TableName() string {
	return "user_login_ips"
}

func init() {
	orm.RegisterModel(new(UserLoginIps))
}

// AddUserLoginIps insert a new UserLoginIps into database and returns
// last inserted Id on success.
func AddUserLoginIps(m *UserLoginIps) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiUserLoginIps mutil insert a new UserLoginIps into database
func AddMultiUserLoginIps(mlist []*UserLoginIps) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetUserLoginIpsById retrieves UserLoginIps by Id. Returns error if
// Id doesn't exist
func GetUserLoginIpsById(id int) (v *UserLoginIps, err error) {
	o := orm.NewOrm()
	v = &UserLoginIps{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserLoginIps retrieves all UserLoginIps matches certain condition. Returns empty list if
// no records exist
func GetAllUserLoginIps(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []UserLoginIps, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserLoginIps))
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

// UpdateUserLoginIps updates UserLoginIps by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserLoginIpsById(m *UserLoginIps) (err error) {
	o := orm.NewOrm()
	v := UserLoginIps{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserLoginIps deletes UserLoginIps by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserLoginIps(id int) (err error) {
	o := orm.NewOrm()
	v := UserLoginIps{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserLoginIps{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
