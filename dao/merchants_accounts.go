package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type MerchantsAccounts struct {
	Id            int       `orm:"column(id);auto"`
	MerchantId    int       `orm:"column(merchant_id)" description:"商户id"`
	Name          string    `orm:"column(name);size(16)" description:"商户名"`
	Balance       float64   `orm:"column(balance);digits(16);decimals(6)" description:"信用额度，0=不限"`
	TransferLimit float64   `orm:"column(transfer_limit);digits(16);decimals(6)" description:"转账限额，每次用户转入的金额最大值"`
	BetLimit      float64   `orm:"column(bet_limit);digits(16);decimals(6)" description:"投注限额，最大单笔可投注金额"`
	BonusLimit    float64   `orm:"column(bonus_limit);digits(16);decimals(6)" description:"单期限红，也就是每期接收最大投注"`
	ProfitScale   float64   `orm:"column(profit_scale);digits(16);decimals(6)" description:"抽水比例"`
	ProfitType    int       `orm:"column(profit_type)" description:"计费方式：1=负盈利模式，不计算限额，只是从负盈利报表中读取;2=信用额度减扣，每次投注从信用额度中扣除信用额度"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null;auto_now_add" description:"创建时间"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add" description:"更新时间"`
	BackupMadeAt  time.Time `orm:"column(backup_made_at);type(datetime);null;auto_now" description:"数据库更新时间"`
}

func (t *MerchantsAccounts) TableName() string {
	return "merchants_accounts"
}

func init() {
	orm.RegisterModel(new(MerchantsAccounts))
}

// AddMerchantsAccounts insert a new MerchantsAccounts into database and returns
// last inserted Id on success.
func AddMerchantsAccounts(m *MerchantsAccounts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiMerchantsAccounts mutil insert a new MerchantsAccounts into database
func AddMultiMerchantsAccounts(mlist []*MerchantsAccounts) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetMerchantsAccountsById retrieves MerchantsAccounts by Id. Returns error if
// Id doesn't exist
func GetMerchantsAccountsById(id int) (v *MerchantsAccounts, err error) {
	o := orm.NewOrm()
	v = &MerchantsAccounts{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMerchantsAccounts retrieves all MerchantsAccounts matches certain condition. Returns empty list if
// no records exist
func GetAllMerchantsAccounts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []MerchantsAccounts, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MerchantsAccounts))
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

// UpdateMerchantsAccounts updates MerchantsAccounts by Id and returns error if
// the record to be updated doesn't exist
func UpdateMerchantsAccountsById(m *MerchantsAccounts) (err error) {
	o := orm.NewOrm()
	v := MerchantsAccounts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMerchantsAccounts deletes MerchantsAccounts by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMerchantsAccounts(id int) (err error) {
	o := orm.NewOrm()
	v := MerchantsAccounts{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MerchantsAccounts{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
