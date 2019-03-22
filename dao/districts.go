package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Districts struct {
	Id          int    `orm:"column(id);auto"`
	ProvinceId  uint   `orm:"column(province_id);null"`
	ParentId    uint   `orm:"column(parent_id);null" description:"上一级id"`
	Lft         int    `orm:"column(lft);null"`
	Rght        int    `orm:"column(rght);null"`
	Name        string `orm:"column(name);size(50)" description:"地区名称"`
	EnglishName string `orm:"column(english_name);size(50)" description:"英文名称"`
	Fullname    string `orm:"column(fullname);size(255)" description:"全名"`
	Zipcode     string `orm:"column(zipcode);size(6)" description:"邮编"`
	Telecode    string `orm:"column(telecode);size(5)" description:"电话区号"`
	Ext         string `orm:"column(ext);size(10)" description:"扩展说明"`
	Disabled    uint8  `orm:"column(disabled)" description:"是否禁用"`
	Sequence    uint   `orm:"column(sequence);null" description:"排序"`
}

func (t *Districts) TableName() string {
	return "districts"
}

func init() {
	orm.RegisterModel(new(Districts))
}

// AddDistricts insert a new Districts into database and returns
// last inserted Id on success.
func AddDistricts(m *Districts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiDistricts mutil insert a new Districts into database
func AddMultiDistricts(mlist []*Districts) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetDistrictsById retrieves Districts by Id. Returns error if
// Id doesn't exist
func GetDistrictsById(id int) (v *Districts, err error) {
	o := orm.NewOrm()
	v = &Districts{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDistricts retrieves all Districts matches certain condition. Returns empty list if
// no records exist
func GetAllDistricts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Districts, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Districts))
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

// UpdateDistricts updates Districts by Id and returns error if
// the record to be updated doesn't exist
func UpdateDistrictsById(m *Districts) (err error) {
	o := orm.NewOrm()
	v := Districts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDistricts deletes Districts by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDistricts(id int) (err error) {
	o := orm.NewOrm()
	v := Districts{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Districts{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
