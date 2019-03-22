package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserLoginLogs struct {
	Id            int       `orm:"column(id);auto"`
	MerchantId    uint      `orm:"column(merchant_id)"`
	UserId        uint      `orm:"column(user_id)"`
	Username      string    `orm:"column(username);size(16);null"`
	IsTester      int8      `orm:"column(is_tester);null"`
	TerminalId    uint8     `orm:"column(terminal_id);null"`
	Ip            string    `orm:"column(ip);size(15)"`
	SignedTime    uint      `orm:"column(signed_time)"`
	SessionId     string    `orm:"column(session_id);size(64);null"`
	HttpUserAgent string    `orm:"column(http_user_agent);size(10240);null"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *UserLoginLogs) TableName() string {
	return "user_login_logs"
}

func init() {
	orm.RegisterModel(new(UserLoginLogs))
}

// AddUserLoginLogs insert a new UserLoginLogs into database and returns
// last inserted Id on success.
func AddUserLoginLogs(m *UserLoginLogs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiUserLoginLogs mutil insert a new UserLoginLogs into database
func AddMultiUserLoginLogs(mlist []*UserLoginLogs) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetUserLoginLogsById retrieves UserLoginLogs by Id. Returns error if
// Id doesn't exist
func GetUserLoginLogsById(id int) (v *UserLoginLogs, err error) {
	o := orm.NewOrm()
	v = &UserLoginLogs{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserLoginLogs retrieves all UserLoginLogs matches certain condition. Returns empty list if
// no records exist
func GetAllUserLoginLogs(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []UserLoginLogs, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserLoginLogs))
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

// UpdateUserLoginLogs updates UserLoginLogs by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserLoginLogsById(m *UserLoginLogs) (err error) {
	o := orm.NewOrm()
	v := UserLoginLogs{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserLoginLogs deletes UserLoginLogs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserLoginLogs(id int) (err error) {
	o := orm.NewOrm()
	v := UserLoginLogs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserLoginLogs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}