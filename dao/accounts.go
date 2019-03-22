package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Accounts struct {
	Id           int       `orm:"column(id);auto" json:"id,string"`
	MerchantId   int       `orm:"column(merchant_id)" json:"merchant_id,string" description:"商户id"`
	UserId       int64     `orm:"column(user_id)" json:"user_id,string" description:"用户id"`
	Username     string    `orm:"column(username);size(16)" json:"username,string" description:"用户名"`
	IsTester     int8      `orm:"column(is_tester);null" json:"is_tester,string" description:"1=测试"`
	Balance      float64   `orm:"column(balance);digits(16);decimals(6)" json:"balance,string" description:"总额度"`
	Frozen       float64   `orm:"column(frozen);digits(16);decimals(6)" json:"frozen,string" description:"冻结额度"`
	Available    float64   `orm:"column(available);digits(16);decimals(6)" json:"available,string" description:"可用额度"`
	Status       uint8     `orm:"column(status)" json:"status,string" description:"状态:1=正常，-1=删除"`
	Locked       uint64    `orm:"column(locked)" json:"locked,string" description:"1=冻结"`
	CreatedAt    time.Time `orm:"column(created_at);type(datetime);null;auto_now_add" json:"created_at,string" description:"创建时间"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add" json:"updated_at,string" description:"更新时间"`
	BackupMadeAt time.Time `orm:"column(backup_made_at);type(datetime);null" json:"backup_made_at,string" description:"数据库更新时间"`
}

func (t *Accounts) TableName() string {
	return "accounts"
}

func init() {
	orm.RegisterModel(new(Accounts))
}

// AddAccounts insert a new Accounts into database and returns
// last inserted Id on success.
func AddAccounts(m *Accounts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAccounts mutil insert a new Accounts into database
func AddMultiAccounts(mlist []*Accounts) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAccountsById retrieves Accounts by Id. Returns error if
// Id doesn't exist
func GetAccountsById(id int) (v *Accounts, err error) {
	o := orm.NewOrm()
	v = &Accounts{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAccounts retrieves all Accounts matches certain condition. Returns empty list if
// no records exist
func GetAllAccounts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Accounts, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Accounts))
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

// UpdateAccounts updates Accounts by Id and returns error if
// the record to be updated doesn't exist
func UpdateAccountsById(o orm.Ormer, m *Accounts) (err error) {
	v := Accounts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAccounts deletes Accounts by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAccounts(id int) (err error) {
	o := orm.NewOrm()
	v := Accounts{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Accounts{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
