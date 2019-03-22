package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type TransactionTypes struct {
	Id            int       `orm:"column(id);auto" json:"id,string"`
	ParentId      uint      `orm:"column(parent_id);null" json:"parent_id,string" description:"父ＩＤ  关联本表中主键"`
	Parent        string    `orm:"column(parent);size(30);null" json:"parent,string"　description:"父级cn_title"`
	FundFlowId    uint      `orm:"column(fund_flow_id);null" json:"fund_flow_id,string"　description:"资金流 关联fund_flows中ＩＤ"`
	Description   string    `orm:"column(description);size(30);null" json:"description,string"　description:"描述"`
	CnTitle       string    `orm:"column(cn_title);size(30);null" json:"cn_title,string"　description:"中文标题"`
	Balance       int8      `orm:"column(balance);null" json:"balance,string"　description:"总额度"`
	Available     int8      `orm:"column(available);null" json:"available,string"　description:"可用额度"`
	Frozen        int8      `orm:"column(frozen);null" json:"frozen,string"　description:"冻结额度"`
	Credit        int8      `orm:"column(credit)" json:"credit,string"　description:"收入"`
	Debit         int8      `orm:"column(debit)" json:"debit,string"　description:"支出"`
	ProjectLinked int8      `orm:"column(project_linked);null" json:"project_linked,string"　description:"关联注单"`
	TraceLinked   int8      `orm:"column(trace_linked);null" json:"trace_linked,string"　description:"关联追号"`
	ReverseType   uint      `orm:"column(reverse_type);null" json:"reverse_type,string"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *TransactionTypes) TableName() string {
	return "transaction_types"
}

func init() {
	orm.RegisterModel(new(TransactionTypes))
}

// AddTransactionTypes insert a new TransactionTypes into database and returns
// last inserted Id on success.
func AddTransactionTypes(m *TransactionTypes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiTransactionTypes mutil insert a new TransactionTypes into database
func AddMultiTransactionTypes(mlist []*TransactionTypes) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetTransactionTypesById retrieves TransactionTypes by Id. Returns error if
// Id doesn't exist
func GetTransactionTypesById(id int) (v *TransactionTypes, err error) {
	o := orm.NewOrm()
	v = &TransactionTypes{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTransactionTypes retrieves all TransactionTypes matches certain condition. Returns empty list if
// no records exist
func GetAllTransactionTypes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*TransactionTypes, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(TransactionTypes))
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

// UpdateTransactionTypes updates TransactionTypes by Id and returns error if
// the record to be updated doesn't exist
func UpdateTransactionTypesById(m *TransactionTypes) (err error) {
	o := orm.NewOrm()
	v := TransactionTypes{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTransactionTypes deletes TransactionTypes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTransactionTypes(id int) (err error) {
	o := orm.NewOrm()
	v := TransactionTypes{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TransactionTypes{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
