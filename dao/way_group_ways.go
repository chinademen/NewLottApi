package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type WayGroupWays struct {
	Id          int       `orm:"column(id);auto" json:"id,string"`
	SeriesId    uint8     `orm:"column(series_id)" json:"series_id,string" description:"系列"`
	TerminalId  uint8     `orm:"column(terminal_id)" json:"terminal_id,string" description:"终端"`
	GroupId     uint      `orm:"column(group_id)" json:"group_id,string" description:"组"`
	SeriesWayId uint      `orm:"column(series_way_id)" json:"series_way_id,string" description:"系列投注方式id"`
	Title       string    `orm:"column(title);size(20);null" json:"title,string" description:"标题"`
	EnTitle     string    `orm:"column(en_title);size(30);null" json:"en_title,string" description:"英文标题"`
	ForDisplay  uint8     `orm:"column(for_display)" json:"for_display,string"`
	ForSearch   uint8     `orm:"column(for_search)" json:"for_search,string"`
	ForMobile   int8      `orm:"column(for_mobile)" json:"for_mobile,string"`
	Sequence    uint      `orm:"column(sequence)" json:"sequence,string" description:"排序"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *WayGroupWays) TableName() string {
	return "way_group_ways"
}

func init() {
	orm.RegisterModel(new(WayGroupWays))
}

// AddWayGroupWays insert a new WayGroupWays into database and returns
// last inserted Id on success.
func AddWayGroupWays(m *WayGroupWays) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiWayGroupWays mutil insert a new WayGroupWays into database
func AddMultiWayGroupWays(mlist []*WayGroupWays) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetWayGroupWaysById retrieves WayGroupWays by Id. Returns error if
// Id doesn't exist
func GetWayGroupWaysById(id int) (v *WayGroupWays, err error) {
	o := orm.NewOrm()
	v = &WayGroupWays{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWayGroupWays retrieves all WayGroupWays matches certain condition. Returns empty list if
// no records exist
func GetAllWayGroupWays(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []WayGroupWays, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(WayGroupWays))
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

// UpdateWayGroupWays updates WayGroupWays by Id and returns error if
// the record to be updated doesn't exist
func UpdateWayGroupWaysById(m *WayGroupWays) (err error) {
	o := orm.NewOrm()
	v := WayGroupWays{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWayGroupWays deletes WayGroupWays by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWayGroupWays(id int) (err error) {
	o := orm.NewOrm()
	v := WayGroupWays{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&WayGroupWays{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
