package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Releases struct {
	Id          int       `orm:"column(id);auto"`
	MerchantId  int       `orm:"column(merchant_id);null" description:"商户id"`
	TerminalId  uint8     `orm:"column(terminal_id)" description:"终端id"`
	Version     string    `orm:"column(version);size(12)" description:"版本"`
	Filename    string    `orm:"column(filename);size(200)" description:"文件名称"`
	Description string    `orm:"column(description)" description:"描述"`
	IsForce     int8      `orm:"column(is_force)" description:"更新，1=强制更新"`
	StartTime   time.Time `orm:"column(start_time);type(datetime)" description:"生效时间"`
	Status      int8      `orm:"column(status)" description:"状态，1=可用"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *Releases) TableName() string {
	return "releases"
}

func init() {
	orm.RegisterModel(new(Releases))
}

// AddReleases insert a new Releases into database and returns
// last inserted Id on success.
func AddReleases(m *Releases) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiReleases mutil insert a new Releases into database
func AddMultiReleases(mlist []*Releases) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetReleasesById retrieves Releases by Id. Returns error if
// Id doesn't exist
func GetReleasesById(id int) (v *Releases, err error) {
	o := orm.NewOrm()
	v = &Releases{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllReleases retrieves all Releases matches certain condition. Returns empty list if
// no records exist
func GetAllReleases(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Releases, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Releases))
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

// UpdateReleases updates Releases by Id and returns error if
// the record to be updated doesn't exist
func UpdateReleasesById(m *Releases) (err error) {
	o := orm.NewOrm()
	v := Releases{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteReleases deletes Releases by Id and returns error if
// the record to be deleted doesn't exist
func DeleteReleases(id int) (err error) {
	o := orm.NewOrm()
	v := Releases{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Releases{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
