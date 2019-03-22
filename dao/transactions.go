package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Transactions struct {
	Id                int     `orm:"column(id);auto" json:"id,string"`
	MerchantId        uint    `orm:"column(merchant_id)" json:"merchant_id,string" description:"商户ＩＤ"`
	TerminalId        uint8   `orm:"column(terminal_id);null" json:"terminal_id,string" description:"终端ＩＤ"`
	SerialNumber      string  `orm:"column(serial_number);size(255)" json:"serial_number,string" description:"帐变序列号"`
	UserId            uint64  `orm:"column(user_id)" json:"user_id,string" description:"用户ＩＤ"`
	Username          string  `orm:"column(username);size(16);null" json:"username,string" description:"用户名"`
	IsTester          int8    `orm:"column(is_tester);null" json:"is_tester,string" description:"是否测试 0-否 1-是"`
	AccountId         uint64  `orm:"column(account_id)" json:"account_id,string" description:"账户ＩＤ"`
	TypeId            uint32  `orm:"column(type_id)" json:"type_id,string" description:"帐变类型ＩＤ"`
	IsIncome          uint8   `orm:"column(is_income)" json:"is_income,string" description:"是否入账 0-是 1-否"`
	TraceId           uint64  `orm:"column(trace_id);null" json:"trace_id,string" description:"追号ＩＤ"`
	LotteryId         uint8   `orm:"column(lottery_id);null" json:"lottery_id,string" description:"彩种ＩＤ"`
	Issue             string  `orm:"column(issue);size(20);null" json:"issue,string" description:"奖期"`
	MethodId          uint8   `orm:"column(method_id);null" json:"method_id,string" description:"玩法ID"`
	WayId             uint    `orm:"column(way_id);null" json:"way_id,string" description:"投注方式"`
	Coefficient       float64 `orm:"column(coefficient);null;digits(4);decimals(3)" json:"coefficient,string" description:"金额模式"`
	Description       string  `orm:"column(description);size(50)" json:"description,string" description:"描述（来源于帐变类型）"`
	ProjectId         uint64  `orm:"column(project_id);null" json:"project_id,string" description:"注单ＩＤ"`
	ProjectNo         string  `orm:"column(project_no);size(32);null" json:"project_no,string" description:"注单号"`
	Amount            float64 `orm:"column(amount);digits(16);decimals(6)" json:"amount,string" description:"金额"`
	Note              string  `orm:"column(note);size(100);null" json:"note,string" description:"备注"`
	PreviousBalance   float64 `orm:"column(previous_balance);digits(16);decimals(6)" json:"previous_balance,string" description:"帐变前总额度"`
	PreviousFrozen    float64 `orm:"column(previous_frozen);digits(16);decimals(6)" json:"previous_frozen,string" description:"帐变前冻结额度"`
	PreviousAvailable float64 `orm:"column(previous_available);digits(16);decimals(6)" json:"previous_available,string" description:"帐变前可用额度"`
	Balance           float64 `orm:"column(balance);digits(16);decimals(6)" json:"balance,string" description:"帐变后总额"`
	Frozen            float64 `orm:"column(frozen);digits(16);decimals(6)" json:"frozen,string" description:"帐变后冻结金额"`
	Available         float64 `orm:"column(available);digits(16);decimals(6)" json:"available,string" description:"帐变后可用额度"`
	Tag               string  `orm:"column(tag);size(30);null" json:"tag,string" description:"额外标记"`
	AdminUserId       uint    `orm:"column(admin_user_id);null" json:"admin_user_id,string"`
	Administrator     string  `orm:"column(administrator);size(16);null" json:"administrator,string"`
	Ip                string  `orm:"column(ip);size(15);null" json:"ip,string"`
	ProxyIp           string  `orm:"column(proxy_ip);size(15);null" json:"proxy_ip,string"`
	Safekey           string  `orm:"column(safekey);size(32)" json:"safekey,string"`
	CreatedAt         string  `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt         string  `orm:"column(updated_at);type(datetime);null;auto_now_add"`
	ExtraData         string  `orm:"column(extra_data);size(1024);null" json:"extra_data,string"`
}

func (t *Transactions) TableName() string {
	return "transactions"
}

func init() {
	orm.RegisterModel(new(Transactions))
}

// AddTransactions insert a new Transactions into database and returns
// last inserted Id on success.
func AddTransactions(o orm.Ormer, m *Transactions) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// AddMultiTransactions mutil insert a new Transactions into database
func AddMultiTransactions(mlist []*Transactions) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetTransactionsById retrieves Transactions by Id. Returns error if
// Id doesn't exist
func GetTransactionsById(id int) (v *Transactions, err error) {
	o := orm.NewOrm()
	v = &Transactions{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTransactions retrieves all Transactions matches certain condition. Returns empty list if
// no records exist
func GetAllTransactions(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Transactions, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transactions))
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

// UpdateTransactions updates Transactions by Id and returns error if
// the record to be updated doesn't exist
func UpdateTransactionsById(m *Transactions) (err error) {
	o := orm.NewOrm()
	v := Transactions{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTransactions deletes Transactions by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTransactions(id int) (err error) {
	o := orm.NewOrm()
	v := Transactions{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Transactions{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
