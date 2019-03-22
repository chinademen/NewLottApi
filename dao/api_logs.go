package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ApiLogs struct {
	Id          int       `orm:"column(id);auto"`
	MerchantId  string    `orm:"column(merchant_id);null" description:"商户id"`
	Ip          string    `orm:"column(ip);size(15);null" description:"呼叫的ip"`
	Url         string    `orm:"column(url);size(150);null" description:"呼叫的url"`
	Domain      string    `orm:"column(domain);size(50);null" description:"使用的域名"`
	RequestBody string    `orm:"column(request_body);null" description:"使用的参数"`
	ServerId    string    `orm:"column(server_id);size(50);null" description:"服务器的id"`
	StartTime   time.Time `orm:"column(start_time);type(datetime);null" description:"呼叫接收时间"`
	EndTime     time.Time `orm:"column(end_time);type(datetime);null" description:"呼叫响应时间"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *ApiLogs) TableName() string {
	return "api_logs"
}

func init() {
	orm.RegisterModel(new(ApiLogs))
}

// AddApiLogs insert a new ApiLogs into database and returns
// last inserted Id on success.
func AddApiLogs(m *ApiLogs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiApiLogs mutil insert a new ApiLogs into database
func AddMultiApiLogs(mlist []*ApiLogs) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetApiLogsById retrieves ApiLogs by Id. Returns error if
// Id doesn't exist
func GetApiLogsById(id int) (v *ApiLogs, err error) {
	o := orm.NewOrm()
	v = &ApiLogs{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApiLogs retrieves all ApiLogs matches certain condition. Returns empty list if
// no records exist
func GetAllApiLogs(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []ApiLogs, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ApiLogs))
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

// UpdateApiLogs updates ApiLogs by Id and returns error if
// the record to be updated doesn't exist
func UpdateApiLogsById(m *ApiLogs) (err error) {
	o := orm.NewOrm()
	v := ApiLogs{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApiLogs deletes ApiLogs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApiLogs(id int) (err error) {
	o := orm.NewOrm()
	v := ApiLogs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApiLogs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
