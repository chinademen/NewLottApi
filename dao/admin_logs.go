package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminLogs struct {
	Id                 int       `orm:"column(id);auto"`
	IsAdmin            uint8     `orm:"column(is_admin)" description:"是否管理员"`
	UserId             uint      `orm:"column(user_id)" description:"用户id"`
	Username           string    `orm:"column(username);size(16)" description:"用户名字"`
	FunctionalityId    uint      `orm:"column(functionality_id)"`
	FunctionalityTitle string    `orm:"column(functionality_title);size(50);null"`
	Controller         string    `orm:"column(controller);size(40)"`
	Action             string    `orm:"column(action);size(40)"`
	Ip                 string    `orm:"column(ip);size(15);null"`
	ProxyIp            string    `orm:"column(proxy_ip);size(15);null"`
	Domain             string    `orm:"column(domain);size(50);null"`
	Env                string    `orm:"column(env);null"`
	RequestUri         string    `orm:"column(request_uri);size(1024);null"`
	RequestData        string    `orm:"column(request_data);null"`
	CreatedAt          time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt          time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *AdminLogs) TableName() string {
	return "admin_logs"
}

func init() {
	orm.RegisterModel(new(AdminLogs))
}

// AddAdminLogs insert a new AdminLogs into database and returns
// last inserted Id on success.
func AddAdminLogs(m *AdminLogs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdminLogs mutil insert a new AdminLogs into database
func AddMultiAdminLogs(mlist []*AdminLogs) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdminLogsById retrieves AdminLogs by Id. Returns error if
// Id doesn't exist
func GetAdminLogsById(id int) (v *AdminLogs, err error) {
	o := orm.NewOrm()
	v = &AdminLogs{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdminLogs retrieves all AdminLogs matches certain condition. Returns empty list if
// no records exist
func GetAllAdminLogs(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdminLogs, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdminLogs))
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

// UpdateAdminLogs updates AdminLogs by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdminLogsById(m *AdminLogs) (err error) {
	o := orm.NewOrm()
	v := AdminLogs{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdminLogs deletes AdminLogs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdminLogs(id int) (err error) {
	o := orm.NewOrm()
	v := AdminLogs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdminLogs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
