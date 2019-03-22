package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdTypes struct {
	Id        int       `orm:"column(id);auto"`
	TypeName  string    `orm:"column(type_name);size(50)"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *AdTypes) TableName() string {
	return "ad_types"
}

func init() {
	orm.RegisterModel(new(AdTypes))
}

// AddAdTypes insert a new AdTypes into database and returns
// last inserted Id on success.
func AddAdTypes(m *AdTypes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdTypes mutil insert a new AdTypes into database
func AddMultiAdTypes(mlist []*AdTypes) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdTypesById retrieves AdTypes by Id. Returns error if
// Id doesn't exist
func GetAdTypesById(id int) (v *AdTypes, err error) {
	o := orm.NewOrm()
	v = &AdTypes{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdTypes retrieves all AdTypes matches certain condition. Returns empty list if
// no records exist
func GetAllAdTypes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdTypes, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdTypes))
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

// UpdateAdTypes updates AdTypes by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdTypesById(m *AdTypes) (err error) {
	o := orm.NewOrm()
	v := AdTypes{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdTypes deletes AdTypes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdTypes(id int) (err error) {
	o := orm.NewOrm()
	v := AdTypes{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdTypes{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
