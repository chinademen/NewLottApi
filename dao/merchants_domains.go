package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type MerchantsDomains struct {
	Id         int       `orm:"column(id);auto" description:"域名ID"`
	MerchantId int       `orm:"column(merchant_id);null" description:"商户id"`
	Domain     string    `orm:"column(domain);size(60)" description:"域名内容"`
	Status     uint8     `orm:"column(status)" description:"域名状态(0未启用,1使用中,2删除)"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null;auto_now_add" description:"创建时间"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add" description:"更新时间"`
}

func (t *MerchantsDomains) TableName() string {
	return "merchants_domains"
}

func init() {
	orm.RegisterModel(new(MerchantsDomains))
}

// AddMerchantsDomains insert a new MerchantsDomains into database and returns
// last inserted Id on success.
func AddMerchantsDomains(m *MerchantsDomains) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiMerchantsDomains mutil insert a new MerchantsDomains into database
func AddMultiMerchantsDomains(mlist []*MerchantsDomains) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetMerchantsDomainsById retrieves MerchantsDomains by Id. Returns error if
// Id doesn't exist
func GetMerchantsDomainsById(id int) (v *MerchantsDomains, err error) {
	o := orm.NewOrm()
	v = &MerchantsDomains{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMerchantsDomains retrieves all MerchantsDomains matches certain condition. Returns empty list if
// no records exist
func GetAllMerchantsDomains(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []MerchantsDomains, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MerchantsDomains))
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

// UpdateMerchantsDomains updates MerchantsDomains by Id and returns error if
// the record to be updated doesn't exist
func UpdateMerchantsDomainsById(m *MerchantsDomains) (err error) {
	o := orm.NewOrm()
	v := MerchantsDomains{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMerchantsDomains deletes MerchantsDomains by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMerchantsDomains(id int) (err error) {
	o := orm.NewOrm()
	v := MerchantsDomains{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MerchantsDomains{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
