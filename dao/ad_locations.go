package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdLocations struct {
	Id          int       `orm:"column(id);auto"`
	Name        string    `orm:"column(name);size(50)" description:"广告位置名称"`
	TypeId      uint      `orm:"column(type_id)" description:"广告位类型"`
	TypeName    string    `orm:"column(type_name);size(50);null"`
	Description string    `orm:"column(description);size(100)" description:"广告位描述"`
	TextLength  int       `orm:"column(text_length);null" description:"广告文本长度限制"`
	PicWidth    int       `orm:"column(pic_width);null" description:"广告图片width"`
	PicHeight   int       `orm:"column(pic_height);null" description:"广告图片height"`
	IsClosed    int8      `orm:"column(is_closed)" description:"状态, 1:禁用, 0:启用"`
	RollTime    int       `orm:"column(roll_time);null" description:"滚动时间"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *AdLocations) TableName() string {
	return "ad_locations"
}

func init() {
	orm.RegisterModel(new(AdLocations))
}

// AddAdLocations insert a new AdLocations into database and returns
// last inserted Id on success.
func AddAdLocations(m *AdLocations) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdLocations mutil insert a new AdLocations into database
func AddMultiAdLocations(mlist []*AdLocations) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdLocationsById retrieves AdLocations by Id. Returns error if
// Id doesn't exist
func GetAdLocationsById(id int) (v *AdLocations, err error) {
	o := orm.NewOrm()
	v = &AdLocations{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdLocations retrieves all AdLocations matches certain condition. Returns empty list if
// no records exist
func GetAllAdLocations(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdLocations, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdLocations))
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

// UpdateAdLocations updates AdLocations by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdLocationsById(m *AdLocations) (err error) {
	o := orm.NewOrm()
	v := AdLocations{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdLocations deletes AdLocations by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdLocations(id int) (err error) {
	o := orm.NewOrm()
	v := AdLocations{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdLocations{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
