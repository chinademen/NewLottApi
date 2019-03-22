package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SeriesMethods struct {
	Id            int       `orm:"column(id);auto"`
	SeriesId      uint8     `orm:"column(series_id)" description:"系列id"`
	Name          string    `orm:"column(name);size(30);null" description:"名称"`
	BasicMethodId uint32    `orm:"column(basic_method_id);null" description:"基础玩法id"`
	IsAdjacent    uint8     `orm:"column(is_adjacent);null"`
	Offset        int8      `orm:"column(offset);null" description:"起始位"`
	Position      string    `orm:"column(position);size(100);null"`
	Hidden        int8      `orm:"column(hidden)" description:"隐藏"`
	Open          int8      `orm:"column(open)" description:"开放"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *SeriesMethods) TableName() string {
	return "series_methods"
}

func init() {
	orm.RegisterModel(new(SeriesMethods))
}

// AddSeriesMethods insert a new SeriesMethods into database and returns
// last inserted Id on success.
func AddSeriesMethods(m *SeriesMethods) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSeriesMethods mutil insert a new SeriesMethods into database
func AddMultiSeriesMethods(mlist []*SeriesMethods) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSeriesMethodsById retrieves SeriesMethods by Id. Returns error if
// Id doesn't exist
func GetSeriesMethodsById(id int) (v *SeriesMethods, err error) {
	o := orm.NewOrm()
	v = &SeriesMethods{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSeriesMethods retrieves all SeriesMethods matches certain condition. Returns empty list if
// no records exist
func GetAllSeriesMethods(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []SeriesMethods, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SeriesMethods))
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

// UpdateSeriesMethods updates SeriesMethods by Id and returns error if
// the record to be updated doesn't exist
func UpdateSeriesMethodsById(m *SeriesMethods) (err error) {
	o := orm.NewOrm()
	v := SeriesMethods{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSeriesMethods deletes SeriesMethods by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSeriesMethods(id int) (err error) {
	o := orm.NewOrm()
	v := SeriesMethods{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SeriesMethods{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
