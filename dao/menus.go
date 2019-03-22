package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Menus struct {
	Id              int       `orm:"column(id);auto" description:"菜单ID"`
	Title           string    `orm:"column(title);size(64)"`
	ParentId        uint      `orm:"column(parent_id);null"`
	Parent          string    `orm:"column(parent);size(50)"`
	ForefatherIds   string    `orm:"column(forefather_ids);size(100);null"`
	Forefathers     string    `orm:"column(forefathers);size(10240);null"`
	FunctionalityId uint      `orm:"column(functionality_id);null"`
	Description     string    `orm:"column(description);size(255);null"`
	Controller      string    `orm:"column(controller);size(40);null"`
	Action          string    `orm:"column(action);size(40);null"`
	Realm           uint8     `orm:"column(realm);null"`
	Params          string    `orm:"column(params);size(100);null"`
	NewWindow       int8      `orm:"column(new_window)"`
	Disabled        uint8     `orm:"column(disabled)" description:"菜单是否启用（0 正常 1关闭）"`
	Sequence        uint      `orm:"column(sequence)" description:"菜单排序"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *Menus) TableName() string {
	return "menus"
}

func init() {
	orm.RegisterModel(new(Menus))
}

// AddMenus insert a new Menus into database and returns
// last inserted Id on success.
func AddMenus(m *Menus) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiMenus mutil insert a new Menus into database
func AddMultiMenus(mlist []*Menus) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetMenusById retrieves Menus by Id. Returns error if
// Id doesn't exist
func GetMenusById(id int) (v *Menus, err error) {
	o := orm.NewOrm()
	v = &Menus{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMenus retrieves all Menus matches certain condition. Returns empty list if
// no records exist
func GetAllMenus(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Menus, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menus))
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

// UpdateMenus updates Menus by Id and returns error if
// the record to be updated doesn't exist
func UpdateMenusById(m *Menus) (err error) {
	o := orm.NewOrm()
	v := Menus{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMenus deletes Menus by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMenus(id int) (err error) {
	o := orm.NewOrm()
	v := Menus{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Menus{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
