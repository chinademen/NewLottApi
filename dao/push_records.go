package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PushRecords struct {
	Id              int       `orm:"column(id);auto"`
	QueueId         uint64    `orm:"column(queue_id);null"`
	CustomerId      uint      `orm:"column(customer_id);null"`
	SetUrl          string    `orm:"column(set_url);size(255);null"`
	LotteryId       uint8     `orm:"column(lottery_id)"`
	RequestLottery  string    `orm:"column(request_lottery);size(10);null"`
	Issue           string    `orm:"column(issue);size(15)"`
	CustomerKey     string    `orm:"column(customer_key);size(32);null"`
	CodecenterId    uint8     `orm:"column(codecenter_id)"`
	CodecenterIp    string    `orm:"column(codecenter_ip);size(255);null"`
	CodecenterLogId uint      `orm:"column(codecenter_log_id);null"`
	SafeStr         string    `orm:"column(safe_str);size(32);null"`
	RequestTime     float64   `orm:"column(request_time);null"`
	AcceptTime      float64   `orm:"column(accept_time);null"`
	FinishTime      float64   `orm:"column(finish_time);null"`
	SpentTime       float32   `orm:"column(spent_time);null"`
	Code            string    `orm:"column(code);size(60);null"`
	Response        string    `orm:"column(response);size(100);null"`
	Status          uint      `orm:"column(status);null" description:"完成状态：值为接口文档中规定的值，可能是复合值"`
	RequestData     string    `orm:"column(request_data);null"`
	VerifyData      string    `orm:"column(verify_data);null"`
	VerifyResult    string    `orm:"column(verify_result);null"`
	CreatedAt       time.Time `orm:"column(created_at);type(datetime)"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *PushRecords) TableName() string {
	return "push_records"
}

func init() {
	orm.RegisterModel(new(PushRecords))
}

// AddPushRecords insert a new PushRecords into database and returns
// last inserted Id on success.
func AddPushRecords(m *PushRecords) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPushRecords mutil insert a new PushRecords into database
func AddMultiPushRecords(mlist []*PushRecords) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPushRecordsById retrieves PushRecords by Id. Returns error if
// Id doesn't exist
func GetPushRecordsById(id int) (v *PushRecords, err error) {
	o := orm.NewOrm()
	v = &PushRecords{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPushRecords retrieves all PushRecords matches certain condition. Returns empty list if
// no records exist
func GetAllPushRecords(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []PushRecords, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PushRecords))
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

// UpdatePushRecords updates PushRecords by Id and returns error if
// the record to be updated doesn't exist
func UpdatePushRecordsById(m *PushRecords) (err error) {
	o := orm.NewOrm()
	v := PushRecords{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePushRecords deletes PushRecords by Id and returns error if
// the record to be deleted doesn't exist
func DeletePushRecords(id int) (err error) {
	o := orm.NewOrm()
	v := PushRecords{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PushRecords{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
