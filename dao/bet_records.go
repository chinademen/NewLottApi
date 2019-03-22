package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type BetRecords struct {
	Id             int       `orm:"column(id);auto" json:"id,string"`
	MerchantId     uint      `orm:"column(merchant_id)" json:"merchant_id,string"`
	UserId         uint64    `orm:"column(user_id)" json:"user_id,string"`
	Username       string    `orm:"column(username);size(16)" json:"username,string"`
	IsTester       int8      `orm:"column(is_tester);null" json:"is_tester,string"`
	LotteryId      uint      `orm:"column(lottery_id)" json:"lottery_id,string"`
	BetCount       uint      `orm:"column(bet_count)" json:"bet_count,string"`
	IsTrace        int8      `orm:"column(is_trace);null" json:"is_trace,string"`
	BetData        string    `orm:"column(bet_data);null" json:"bet_data,string"`
	CompressedData string    `orm:"column(compressed_data);null" json:"compressed_data,string"`
	TerminalId     uint8     `orm:"column(terminal_id)" json:"terminal_id,string"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *BetRecords) TableName() string {
	return "bet_records"
}

func init() {
	orm.RegisterModel(new(BetRecords))
}

// AddBetRecords insert a new BetRecords into database and returns
// last inserted Id on success.
func AddBetRecords(m *BetRecords) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiBetRecords mutil insert a new BetRecords into database
func AddMultiBetRecords(mlist []*BetRecords) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetBetRecordsById retrieves BetRecords by Id. Returns error if
// Id doesn't exist
func GetBetRecordsById(id int) (v *BetRecords, err error) {
	o := orm.NewOrm()
	v = &BetRecords{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBetRecords retrieves all BetRecords matches certain condition. Returns empty list if
// no records exist
func GetAllBetRecords(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []BetRecords, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(BetRecords))
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

// UpdateBetRecords updates BetRecords by Id and returns error if
// the record to be updated doesn't exist
func UpdateBetRecordsById(m *BetRecords) (err error) {
	o := orm.NewOrm()
	v := BetRecords{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBetRecords deletes BetRecords by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBetRecords(id int) (err error) {
	o := orm.NewOrm()
	v := BetRecords{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BetRecords{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
