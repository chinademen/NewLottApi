package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdInfos struct {
	Id           int       `orm:"column(id);auto"`
	AdLocationId int       `orm:"column(ad_location_id);null" description:"广告位"`
	Name         string    `orm:"column(name);size(255);null" description:"广告名称"`
	Content      string    `orm:"column(content);null" description:"广告文本内容"`
	PicUrl       string    `orm:"column(pic_url)" description:"广告URL"`
	IsClosed     int8      `orm:"column(is_closed)" description:"状态, 1:禁用, 0:启用"`
	RedirectUrl  string    `orm:"column(redirect_url);size(255)" description:"广告图片"`
	CreatorId    uint      `orm:"column(creator_id)" description:"用户ID"`
	Creator      string    `orm:"column(creator);size(50)" description:"用户名称"`
	Sequence     uint      `orm:"column(sequence)"`
	CreatedAt    time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *AdInfos) TableName() string {
	return "ad_infos"
}

func init() {
	orm.RegisterModel(new(AdInfos))
}

// AddAdInfos insert a new AdInfos into database and returns
// last inserted Id on success.
func AddAdInfos(m *AdInfos) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAdInfos mutil insert a new AdInfos into database
func AddMultiAdInfos(mlist []*AdInfos) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAdInfosById retrieves AdInfos by Id. Returns error if
// Id doesn't exist
func GetAdInfosById(id int) (v *AdInfos, err error) {
	o := orm.NewOrm()
	v = &AdInfos{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdInfos retrieves all AdInfos matches certain condition. Returns empty list if
// no records exist
func GetAllAdInfos(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AdInfos, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AdInfos))
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

// UpdateAdInfos updates AdInfos by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdInfosById(m *AdInfos) (err error) {
	o := orm.NewOrm()
	v := AdInfos{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAdInfos deletes AdInfos by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAdInfos(id int) (err error) {
	o := orm.NewOrm()
	v := AdInfos{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AdInfos{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
