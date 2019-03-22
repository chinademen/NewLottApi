package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AuditTypes struct {
	Id          int       `orm:"column(id);auto"`
	Name        string    `orm:"column(name);size(50)" description:"审核名称"`
	Controller  string    `orm:"column(controller);size(40)" description:"审核控制器"`
	Action      string    `orm:"column(action);size(40)" description:"审核方法"`
	Description string    `orm:"column(description);size(255)" description:"描述"`
	Sequence    uint      `orm:"column(sequence)" description:"排序"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *AuditTypes) TableName() string {
	return "audit_types"
}

func init() {
	orm.RegisterModel(new(AuditTypes))
}

// AddAuditTypes insert a new AuditTypes into database and returns
// last inserted Id on success.
func AddAuditTypes(m *AuditTypes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAuditTypes mutil insert a new AuditTypes into database
func AddMultiAuditTypes(mlist []*AuditTypes) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAuditTypesById retrieves AuditTypes by Id. Returns error if
// Id doesn't exist
func GetAuditTypesById(id int) (v *AuditTypes, err error) {
	o := orm.NewOrm()
	v = &AuditTypes{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAuditTypes retrieves all AuditTypes matches certain condition. Returns empty list if
// no records exist
func GetAllAuditTypes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AuditTypes, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AuditTypes))
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

// UpdateAuditTypes updates AuditTypes by Id and returns error if
// the record to be updated doesn't exist
func UpdateAuditTypesById(m *AuditTypes) (err error) {
	o := orm.NewOrm()
	v := AuditTypes{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAuditTypes deletes AuditTypes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAuditTypes(id int) (err error) {
	o := orm.NewOrm()
	v := AuditTypes{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AuditTypes{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
