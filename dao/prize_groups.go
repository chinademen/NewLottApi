package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PrizeGroups struct {
	Id           int       `orm:"column(id);auto" json:"id,string"`
	SeriesId     uint8     `orm:"column(series_id)" json:"series_id,string" description:"系列id"`
	Type         uint8     `orm:"column(type)" json:"type,string" description:"任务类型"`
	Name         string    `orm:"column(name);size(20)" json:"name,string" description:"名称"`
	ClassicPrize uint32    `orm:"column(classic_prize)" json:"classic_prize,string" description:"经典奖金"`
	Water        float32   `orm:"column(water)" json:"water,string" description:"水率"`
	CreatedAt    time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *PrizeGroups) TableName() string {
	return "prize_groups"
}

func init() {
	orm.RegisterModel(new(PrizeGroups))
}

// AddPrizeGroups insert a new PrizeGroups into database and returns
// last inserted Id on success.
func AddPrizeGroups(m *PrizeGroups) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPrizeGroups mutil insert a new PrizeGroups into database
func AddMultiPrizeGroups(mlist []*PrizeGroups) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPrizeGroupsById retrieves PrizeGroups by Id. Returns error if
// Id doesn't exist
func GetPrizeGroupsById(id int) (v *PrizeGroups, err error) {
	o := orm.NewOrm()
	v = &PrizeGroups{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPrizeGroups retrieves all PrizeGroups matches certain condition. Returns empty list if
// no records exist
func GetAllPrizeGroups(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*PrizeGroups, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PrizeGroups))
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

// UpdatePrizeGroups updates PrizeGroups by Id and returns error if
// the record to be updated doesn't exist
func UpdatePrizeGroupsById(m *PrizeGroups) (err error) {
	o := orm.NewOrm()
	v := PrizeGroups{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePrizeGroups deletes PrizeGroups by Id and returns error if
// the record to be deleted doesn't exist
func DeletePrizeGroups(id int) (err error) {
	o := orm.NewOrm()
	v := PrizeGroups{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PrizeGroups{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
