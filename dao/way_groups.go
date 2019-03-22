package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type WayGroups struct {
	Id         int       `orm:"column(id);auto" json:"id:string"`
	SeriesId   uint8     `orm:"column(series_id)" json:"series_id:string" description:"系列id"`
	TerminalId uint8     `orm:"column(terminal_id)" json:"terminal_id:string" description:"终端id"`
	ParentId   uint      `orm:"column(parent_id);null" json:"parent_id:string" description:"上级ID"`
	Parent     string    `orm:"column(parent);size(20);null" json:"parent:string" description:"上级"`
	Title      string    `orm:"column(title);size(20)" json:"title:string" description:"标题"`
	EnTitle    string    `orm:"column(en_title);size(20)" json:"en_title:string" description:"英文标题"`
	ForDisplay uint8     `orm:"column(for_display)" json:"for_display:string" description:"显示"`
	Sequence   uint      `orm:"column(sequence)" json:"sequence:string" description:"排序"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *WayGroups) TableName() string {
	return "way_groups"
}

func init() {
	orm.RegisterModel(new(WayGroups))
}

// AddWayGroups insert a new WayGroups into database and returns
// last inserted Id on success.
func AddWayGroups(m *WayGroups) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiWayGroups mutil insert a new WayGroups into database
func AddMultiWayGroups(mlist []*WayGroups) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetWayGroupsById retrieves WayGroups by Id. Returns error if
// Id doesn't exist
func GetWayGroupsById(id int) (v *WayGroups, err error) {
	o := orm.NewOrm()
	v = &WayGroups{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWayGroups retrieves all WayGroups matches certain condition. Returns empty list if
// no records exist
func GetAllWayGroups(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []WayGroups, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(WayGroups))
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

// UpdateWayGroups updates WayGroups by Id and returns error if
// the record to be updated doesn't exist
func UpdateWayGroupsById(m *WayGroups) (err error) {
	o := orm.NewOrm()
	v := WayGroups{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWayGroups deletes WayGroups by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWayGroups(id int) (err error) {
	o := orm.NewOrm()
	v := WayGroups{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&WayGroups{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
