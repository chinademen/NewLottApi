package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type TransferInfo struct {
	Id          int       `orm:"column(id);auto"`
	MerchantId  uint      `orm:"column(merchant_id)" description:"商户ID"`
	UserId      uint      `orm:"column(user_id)" description:"用户id"`
	TypeId      uint8     `orm:"column(type_id)" description:"交易类型 1-转入 2-转出"`
	Merchant    string    `orm:"column(merchant);size(20)" description:"商户名"`
	BillNo      string    `orm:"column(bill_no);size(40)" description:"订单号码"`
	OrderNumber string    `orm:"column(order_number);size(40)" description:"商户方订单号码"`
	Amount      float64   `orm:"column(amount)" description:"金额"`
	Status      int8      `orm:"column(status)" description:"订单状态"`
	AcceptedAt  time.Time `orm:"column(accepted_at);type(datetime);null" description:"接受时间"`
	RejectedAt  time.Time `orm:"column(rejected_at);type(datetime);null" description:"拒绝时间"`
	CanceledAt  time.Time `orm:"column(canceled_at);type(datetime);null" description:"取消时间"`
	RefundedAt  time.Time `orm:"column(refunded_at);type(datetime);null"`
	ExpireAt    time.Time `orm:"column(expire_at);type(datetime);null" description:"过期时间"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);auto_now_add"`
}

func (t *TransferInfo) TableName() string {
	return "transfer_info"
}

func init() {
	orm.RegisterModel(new(TransferInfo))
}

// AddTransferInfo insert a new TransferInfo into database and returns
// last inserted Id on success.
func AddTransferInfo(m *TransferInfo) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiTransferInfo mutil insert a new TransferInfo into database
func AddMultiTransferInfo(mlist []*TransferInfo) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetTransferInfoById retrieves TransferInfo by Id. Returns error if
// Id doesn't exist
func GetTransferInfoById(id int) (v *TransferInfo, err error) {
	o := orm.NewOrm()
	v = &TransferInfo{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTransferInfo retrieves all TransferInfo matches certain condition. Returns empty list if
// no records exist
func GetAllTransferInfo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*TransferInfo, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(TransferInfo))
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

// UpdateTransferInfo updates TransferInfo by Id and returns error if
// the record to be updated doesn't exist
func UpdateTransferInfoById(m *TransferInfo) (err error) {
	o := orm.NewOrm()
	v := TransferInfo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTransferInfo deletes TransferInfo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTransferInfo(id int) (err error) {
	o := orm.NewOrm()
	v := TransferInfo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TransferInfo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
