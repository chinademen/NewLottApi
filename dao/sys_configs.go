package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SysConfigs struct {
	Id                   int       `orm:"column(id);auto"`
	ParentId             uint16    `orm:"column(parent_id);null"`
	Parent               string    `orm:"column(parent);size(100);null"`
	Item                 string    `orm:"column(item);size(255)"`
	Value                string    `orm:"column(value);size(1024)"`
	DefaultValue         string    `orm:"column(default_value);size(1024)"`
	Title                string    `orm:"column(title);size(100)"`
	DataType             string    `orm:"column(data_type);size(10)" description:"int,float,string"`
	FormFace             string    `orm:"column(form_face);size(20)"`
	Validation           string    `orm:"column(validation);size(10)"`
	DataSource           string    `orm:"column(data_source);size(1024)"`
	Description          string    `orm:"column(description);size(1024)"`
	FormatedValue        string    `orm:"column(formated_value);size(1024);null"`
	FormatedDefaultValue string    `orm:"column(formated_default_value);size(1024);null"`
	Sequence             uint      `orm:"column(sequence);null"`
	UpdatedAt            time.Time `orm:"column(updated_at);type(timestamp);auto_now_add"`
}

func (t *SysConfigs) TableName() string {
	return "sys_configs"
}

func init() {
	orm.RegisterModel(new(SysConfigs))
}

// AddSysConfigs insert a new SysConfigs into database and returns
// last inserted Id on success.
func AddSysConfigs(m *SysConfigs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSysConfigs mutil insert a new SysConfigs into database
func AddMultiSysConfigs(mlist []*SysConfigs) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSysConfigsById retrieves SysConfigs by Id. Returns error if
// Id doesn't exist
func GetSysConfigsById(id int) (v *SysConfigs, err error) {
	o := orm.NewOrm()
	v = &SysConfigs{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetSysConfigsByName(sField string) (v *SysConfigs, err error) {
	o := orm.NewOrm()
	v = &SysConfigs{Item: sField}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSysConfigs retrieves all SysConfigs matches certain condition. Returns empty list if
// no records exist
func GetAllSysConfigs(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []SysConfigs, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SysConfigs))
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

// UpdateSysConfigs updates SysConfigs by Id and returns error if
// the record to be updated doesn't exist
func UpdateSysConfigsById(m *SysConfigs) (err error) {
	o := orm.NewOrm()
	v := SysConfigs{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSysConfigs deletes SysConfigs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSysConfigs(id int) (err error) {
	o := orm.NewOrm()
	v := SysConfigs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SysConfigs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
