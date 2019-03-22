package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type RestSettings struct {
	Id        int       `orm:"column(id);auto"`
	LotteryId uint8     `orm:"column(lottery_id)" description:"彩种id"`
	Periodic  uint8     `orm:"column(periodic);null" description:"休市类型，0：一次；1：周期性"`
	BeginDate time.Time `orm:"column(begin_date);type(date);null" description:"休市开始日期，休市类型为1时有效"`
	EndDate   time.Time `orm:"column(end_date);type(date);null" description:"休市结束日期，休市类型为1时有效"`
	Week      string    `orm:"column(week);size(20);null"`
	BeginTime time.Time `orm:"column(begin_time);type(time);null" description:"休市开始时间，休市类型为2是有效"`
	EndTime   time.Time `orm:"column(end_time);type(time);null" description:"修饰结束时间，修饰类型为2是有效"`
	KeepIssue uint8     `orm:"column(keep_issue)"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *RestSettings) TableName() string {
	return "rest_settings"
}

func init() {
	orm.RegisterModel(new(RestSettings))
}

// AddRestSettings insert a new RestSettings into database and returns
// last inserted Id on success.
func AddRestSettings(m *RestSettings) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiRestSettings mutil insert a new RestSettings into database
func AddMultiRestSettings(mlist []*RestSettings) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetRestSettingsById retrieves RestSettings by Id. Returns error if
// Id doesn't exist
func GetRestSettingsById(id int) (v *RestSettings, err error) {
	o := orm.NewOrm()
	v = &RestSettings{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRestSettings retrieves all RestSettings matches certain condition. Returns empty list if
// no records exist
func GetAllRestSettings(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []RestSettings, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RestSettings))
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

// UpdateRestSettings updates RestSettings by Id and returns error if
// the record to be updated doesn't exist
func UpdateRestSettingsById(m *RestSettings) (err error) {
	o := orm.NewOrm()
	v := RestSettings{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRestSettings deletes RestSettings by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRestSettings(id int) (err error) {
	o := orm.NewOrm()
	v := RestSettings{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RestSettings{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
