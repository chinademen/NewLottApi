package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type BetNumbers struct {
	Id        int       `orm:"column(id);auto"`
	LotteryId uint      `orm:"column(lottery_id)"`
	Issue     string    `orm:"column(issue);size(15)"`
	Offset    uint8     `orm:"column(offset)"`
	Len       uint8     `orm:"column(len);null"`
	Number    string    `orm:"column(number);size(100)"`
	Multiple  float64   `orm:"column(multiple)"`
	Prize     float64   `orm:"column(prize)"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *BetNumbers) TableName() string {
	return "bet_numbers"
}

func init() {
	orm.RegisterModel(new(BetNumbers))
}

// AddBetNumbers insert a new BetNumbers into database and returns
// last inserted Id on success.
func AddBetNumbers(m *BetNumbers) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiBetNumbers mutil insert a new BetNumbers into database
func AddMultiBetNumbers(mlist []*BetNumbers) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetBetNumbersById retrieves BetNumbers by Id. Returns error if
// Id doesn't exist
func GetBetNumbersById(id int) (v *BetNumbers, err error) {
	o := orm.NewOrm()
	v = &BetNumbers{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBetNumbers retrieves all BetNumbers matches certain condition. Returns empty list if
// no records exist
func GetAllBetNumbers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []BetNumbers, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(BetNumbers))
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

// UpdateBetNumbers updates BetNumbers by Id and returns error if
// the record to be updated doesn't exist
func UpdateBetNumbersById(m *BetNumbers) (err error) {
	o := orm.NewOrm()
	v := BetNumbers{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBetNumbers deletes BetNumbers by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBetNumbers(id int) (err error) {
	o := orm.NewOrm()
	v := BetNumbers{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BetNumbers{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
