package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Merchants struct {
	Id        int       `orm:"column(id);pk" json:"id,string"`
	Identity  string    `orm:"column(identity);size(50)" json:"identity,string" description:"商户唯一标识，前缀"`
	Name      string    `orm:"column(name);size(50);null" json:"name,string" description:"商户名"`
	WalletId  int       `orm:"column(wallet_id)" json:"wallet_id,string" description:"商户荷包id"`
	SafeKey   string    `orm:"column(safe_key);size(32)" json:"safe_key,string" description:"商户唯一密钥"`
	PostUrl   string    `orm:"column(post_url);size(200);null" json:"post_url,string" description:"推送数据url"`
	Status    int8      `orm:"column(status);null" json:"status,string" description:"0: 未激活, 1: 激活"`
	IsTester  int8      `orm:"column(is_tester)" json:"is_tester,string" description:"0: 真实商户, 1: 测试商户"`
	Template  uint8     `orm:"column(template)" json:"template,string" description:"针对不同商户提供不同的模板样式"`
	Remark    string    `orm:"column(remark);size(200);null" json:"remark,string" description:"备注"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *Merchants) TableName() string {
	return "merchants"
}

func init() {
	orm.RegisterModel(new(Merchants))
}

// AddMerchants insert a new Merchants into database and returns
// last inserted Id on success.
func AddMerchants(m *Merchants) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiMerchants mutil insert a new Merchants into database
func AddMultiMerchants(mlist []*Merchants) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetMerchantsById retrieves Merchants by Id. Returns error if
// Id doesn't exist
func GetMerchantsById(id int) (v *Merchants, err error) {
	o := orm.NewOrm()
	v = &Merchants{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMerchants retrieves all Merchants matches certain condition. Returns empty list if
// no records exist
func GetAllMerchants(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Merchants, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Merchants))
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

// UpdateMerchants updates Merchants by Id and returns error if
// the record to be updated doesn't exist
func UpdateMerchantsById(m *Merchants) (err error) {
	o := orm.NewOrm()
	v := Merchants{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMerchants deletes Merchants by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMerchants(id int) (err error) {
	o := orm.NewOrm()
	v := Merchants{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Merchants{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
