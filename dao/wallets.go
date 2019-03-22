package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Wallets struct {
	Id               int       `orm:"column(id);auto"`
	MerchantId       uint      `orm:"column(merchant_id)" description:"商户id"`
	MerchantIdentity string    `orm:"column(merchant_identity);size(50)" description:"商户唯一标识"`
	IsTester         int8      `orm:"column(is_tester)" description:"0: 真实商户, 1: 测试商户"`
	Balance          float64   `orm:"column(balance);digits(16);decimals(6)"`
	Status           uint8     `orm:"column(status)" description:"0: 禁用, 1: 启用"`
	Locked           uint64    `orm:"column(locked)"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *Wallets) TableName() string {
	return "wallets"
}

func init() {
	orm.RegisterModel(new(Wallets))
}

// AddWallets insert a new Wallets into database and returns
// last inserted Id on success.
func AddWallets(m *Wallets) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiWallets mutil insert a new Wallets into database
func AddMultiWallets(mlist []*Wallets) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetWalletsById retrieves Wallets by Id. Returns error if
// Id doesn't exist
func GetWalletsById(id int) (v *Wallets, err error) {
	o := orm.NewOrm()
	v = &Wallets{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWallets retrieves all Wallets matches certain condition. Returns empty list if
// no records exist
func GetAllWallets(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Wallets, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Wallets))
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

// UpdateWallets updates Wallets by Id and returns error if
// the record to be updated doesn't exist
func UpdateWalletsById(m *Wallets) (err error) {
	o := orm.NewOrm()
	v := Wallets{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWallets deletes Wallets by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWallets(id int) (err error) {
	o := orm.NewOrm()
	v := Wallets{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Wallets{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
