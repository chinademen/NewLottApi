package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminPwdHistories struct {
	Id            int       `orm:"column(id);auto"`
	AdminUserId   uint      `orm:"column(admin_user_id)" description:"管理员ID"`
	AdminUsername string    `orm:"column(admin_username);size(16)" description:"管理员名"`
	Password      string    `orm:"column(password);size(60)" description:"密码"`
	OperatorId    int       `orm:"column(operator_id);null" description:"修改操作人ID"`
	Operator      string    `orm:"column(operator);size(16);null"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *AdminPwdHistories) TableName() string {
	return "admin_pwd_histories"
}

func init() {
	orm.RegisterModel(new(AdminPwdHistories))
}

// AddAdminPwdHistories insert a new AdminPwdHistories into database and returns
// last inserted Id on success.
func AddAdminPwdHistories(m *AdminPwdHistories) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdminPwdHistories mutil insert a new AdminPwdHistories into database
func AddMultiAdminPwdHistories(mlist []*AdminPwdHistories) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdminPwdHistoriesById retrieves AdminPwdHistories by Id. Returns error if
// Id doesn't exist
func GetAdminPwdHistoriesById(id int) (v *AdminPwdHistories, err error) {
	o := orm.NewOrm()
	v = &AdminPwdHistories{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdminPwdHistories retrieves all AdminPwdHistories matches certain condition. Returns empty list if
// no records exist
func GetAllAdminPwdHistories(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdminPwdHistories, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdminPwdHistories))
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

// UpdateAdminPwdHistories updates AdminPwdHistories by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdminPwdHistoriesById(m *AdminPwdHistories) (err error) {
	o := orm.NewOrm()
	v := AdminPwdHistories{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdminPwdHistories deletes AdminPwdHistories by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdminPwdHistories(id int) (err error) {
	o := orm.NewOrm()
	v := AdminPwdHistories{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdminPwdHistories{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
