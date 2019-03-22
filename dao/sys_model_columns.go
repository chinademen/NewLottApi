package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type SysModelColumns struct {
	Id              int    `orm:"column(id);auto"`
	Name            string `orm:"column(name);size(64)"`
	SysModelId      uint   `orm:"column(sys_model_id);null"`
	SysModelName    string `orm:"column(sys_model_name);size(64);null"`
	Db              string `orm:"column(db);size(64)"`
	Tablename       string `orm:"column(table_name);size(64)"`
	OrdinalPosition uint64 `orm:"column(ordinal_position)" description:"原始顺序"`
	ColumnDefault   string `orm:"column(column_default);null" description:"默认值"`
	IsNullable      string `orm:"column(is_nullable);size(3)"`
	DataType        string `orm:"column(data_type);size(64)"`
	MaxLength       uint64 `orm:"column(max_length);null"`
	CharsetName     string `orm:"column(charset_name);size(32);null"`
	ColumnType      string `orm:"column(column_type)" description:"字段类型"`
	ColumnComment   string `orm:"column(column_comment);size(255)"`
	Note            string `orm:"column(note);size(1024)"`
}

func (t *SysModelColumns) TableName() string {
	return "sys_model_columns"
}

func init() {
	orm.RegisterModel(new(SysModelColumns))
}

// AddSysModelColumns insert a new SysModelColumns into database and returns
// last inserted Id on success.
func AddSysModelColumns(m *SysModelColumns) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSysModelColumns mutil insert a new SysModelColumns into database
func AddMultiSysModelColumns(mlist []*SysModelColumns) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSysModelColumnsById retrieves SysModelColumns by Id. Returns error if
// Id doesn't exist
func GetSysModelColumnsById(id int) (v *SysModelColumns, err error) {
	o := orm.NewOrm()
	v = &SysModelColumns{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSysModelColumns retrieves all SysModelColumns matches certain condition. Returns empty list if
// no records exist
func GetAllSysModelColumns(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []SysModelColumns, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SysModelColumns))
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

// UpdateSysModelColumns updates SysModelColumns by Id and returns error if
// the record to be updated doesn't exist
func UpdateSysModelColumnsById(m *SysModelColumns) (err error) {
	o := orm.NewOrm()
	v := SysModelColumns{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSysModelColumns deletes SysModelColumns by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSysModelColumns(id int) (err error) {
	o := orm.NewOrm()
	v := SysModelColumns{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SysModelColumns{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
