package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminUsers struct {
	Id            int       `orm:"column(id);auto"`
	Username      string    `orm:"column(username);size(16)"`
	Name          string    `orm:"column(name);size(50);null"`
	Email         string    `orm:"column(email);size(200);null"`
	Password      string    `orm:"column(password);size(60)"`
	Language      string    `orm:"column(language);size(10);null"`
	MenuLink      int8      `orm:"column(menu_link)" description:"记录级链接"`
	MenuContext   int8      `orm:"column(menu_context)" description:"记录级右键菜单"`
	Rights        string    `orm:"column(rights);size(10240);null"`
	Actived       int8      `orm:"column(actived)" description:"激活"`
	RememberToken string    `orm:"column(remember_token);size(200);null"`
	SigninAt      time.Time `orm:"column(signin_at);type(timestamp);null"`
	PwdSetedAt    time.Time `orm:"column(pwd_seted_at);type(datetime);null"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *AdminUsers) TableName() string {
	return "admin_users"
}

func init() {
	orm.RegisterModel(new(AdminUsers))
}

// AddAdminUsers insert a new AdminUsers into database and returns
// last inserted Id on success.
func AddAdminUsers(m *AdminUsers) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdminUsers mutil insert a new AdminUsers into database
func AddMultiAdminUsers(mlist []*AdminUsers) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdminUsersById retrieves AdminUsers by Id. Returns error if
// Id doesn't exist
func GetAdminUsersById(id int) (v *AdminUsers, err error) {
	o := orm.NewOrm()
	v = &AdminUsers{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdminUsers retrieves all AdminUsers matches certain condition. Returns empty list if
// no records exist
func GetAllAdminUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdminUsers, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdminUsers))
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

// UpdateAdminUsers updates AdminUsers by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdminUsersById(m *AdminUsers) (err error) {
	o := orm.NewOrm()
	v := AdminUsers{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdminUsers deletes AdminUsers by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdminUsers(id int) (err error) {
	o := orm.NewOrm()
	v := AdminUsers{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdminUsers{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
