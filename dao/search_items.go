package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SearchItems struct {
	Id             int       `orm:"column(id);auto"`
	SearchConfigId uint      `orm:"column(search_config_id);null"`
	Model          string    `orm:"column(model);size(40)"`
	Field          string    `orm:"column(field);size(40)"`
	Label          string    `orm:"column(label);size(40)"`
	Type           string    `orm:"column(type);size(20)"`
	DefaultValue   string    `orm:"column(default_value);size(250);null"`
	Source         string    `orm:"column(source)" description:"数据源"`
	Div            int8      `orm:"column(div)"`
	Empty          int8      `orm:"column(empty)"`
	EmptyText      string    `orm:"column(empty_text);size(50)"`
	MatchRule      string    `orm:"column(match_rule);null" description:"匹配规则"`
	Sequence       uint      `orm:"column(sequence);null"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *SearchItems) TableName() string {
	return "search_items"
}

func init() {
	orm.RegisterModel(new(SearchItems))
}

// AddSearchItems insert a new SearchItems into database and returns
// last inserted Id on success.
func AddSearchItems(m *SearchItems) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSearchItems mutil insert a new SearchItems into database
func AddMultiSearchItems(mlist []*SearchItems) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSearchItemsById retrieves SearchItems by Id. Returns error if
// Id doesn't exist
func GetSearchItemsById(id int) (v *SearchItems, err error) {
	o := orm.NewOrm()
	v = &SearchItems{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSearchItems retrieves all SearchItems matches certain condition. Returns empty list if
// no records exist
func GetAllSearchItems(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []SearchItems, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SearchItems))
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

// UpdateSearchItems updates SearchItems by Id and returns error if
// the record to be updated doesn't exist
func UpdateSearchItemsById(m *SearchItems) (err error) {
	o := orm.NewOrm()
	v := SearchItems{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSearchItems deletes SearchItems by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSearchItems(id int) (err error) {
	o := orm.NewOrm()
	v := SearchItems{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SearchItems{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
