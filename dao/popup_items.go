package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PopupItems struct {
	Id        int       `orm:"column(id);auto"`
	PopupId   uint      `orm:"column(popup_id)"`
	Field     string    `orm:"column(field);size(32)"`
	Label     string    `orm:"column(label);size(32)"`
	Type      string    `orm:"column(type);size(16)"`
	Required  int8      `orm:"column(required)"`
	MinLength uint      `orm:"column(min_length);null"`
	MaxLength uint      `orm:"column(max_length);null"`
	Sequence  uint      `orm:"column(sequence);null"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *PopupItems) TableName() string {
	return "popup_items"
}

func init() {
	orm.RegisterModel(new(PopupItems))
}

// AddPopupItems insert a new PopupItems into database and returns
// last inserted Id on success.
func AddPopupItems(m *PopupItems) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPopupItems mutil insert a new PopupItems into database
func AddMultiPopupItems(mlist []*PopupItems) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPopupItemsById retrieves PopupItems by Id. Returns error if
// Id doesn't exist
func GetPopupItemsById(id int) (v *PopupItems, err error) {
	o := orm.NewOrm()
	v = &PopupItems{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPopupItems retrieves all PopupItems matches certain condition. Returns empty list if
// no records exist
func GetAllPopupItems(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []PopupItems, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PopupItems))
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

// UpdatePopupItems updates PopupItems by Id and returns error if
// the record to be updated doesn't exist
func UpdatePopupItemsById(m *PopupItems) (err error) {
	o := orm.NewOrm()
	v := PopupItems{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePopupItems deletes PopupItems by Id and returns error if
// the record to be deleted doesn't exist
func DeletePopupItems(id int) (err error) {
	o := orm.NewOrm()
	v := PopupItems{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PopupItems{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
