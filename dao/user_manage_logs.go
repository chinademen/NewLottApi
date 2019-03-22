package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserManageLogs struct {
	Id             int       `orm:"column(id);auto"`
	UserId         uint      `orm:"column(user_id)"`
	AdminId        uint      `orm:"column(admin_id)" description:"操作管理员id"`
	Admin          string    `orm:"column(admin);size(100)" description:"操作管理员"`
	CommentAdminId uint      `orm:"column(comment_admin_id);null" description:"填写备注的管理员id"`
	CommentAdmin   string    `orm:"column(comment_admin);size(100);null" description:"填写备注的管理员"`
	Comment        string    `orm:"column(comment);null" description:"备注"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null" description:"添加时间"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(timestamp);null" description:"更新时间"`
}

func (t *UserManageLogs) TableName() string {
	return "user_manage_logs"
}

func init() {
	orm.RegisterModel(new(UserManageLogs))
}

// AddUserManageLogs insert a new UserManageLogs into database and returns
// last inserted Id on success.
func AddUserManageLogs(m *UserManageLogs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiUserManageLogs mutil insert a new UserManageLogs into database
func AddMultiUserManageLogs(mlist []*UserManageLogs) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetUserManageLogsById retrieves UserManageLogs by Id. Returns error if
// Id doesn't exist
func GetUserManageLogsById(id int) (v *UserManageLogs, err error) {
	o := orm.NewOrm()
	v = &UserManageLogs{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserManageLogs retrieves all UserManageLogs matches certain condition. Returns empty list if
// no records exist
func GetAllUserManageLogs(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []UserManageLogs, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserManageLogs))
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

// UpdateUserManageLogs updates UserManageLogs by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserManageLogsById(m *UserManageLogs) (err error) {
	o := orm.NewOrm()
	v := UserManageLogs{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserManageLogs deletes UserManageLogs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserManageLogs(id int) (err error) {
	o := orm.NewOrm()
	v := UserManageLogs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserManageLogs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
