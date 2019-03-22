package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PrizeDetails struct {
	Id           int       `orm:"column(id);auto" json:"id,string"`
	SeriesId     uint8     `orm:"column(series_id)" json:"series_id,string" description:"系列ＩＤ"`
	GroupId      uint      `orm:"column(group_id)" json:"group_id,string" description:"奖金组ID"`
	GroupName    string    `orm:"column(group_name);size(20)" json:"group_name,string" description:"奖金组"`
	ClassicPrize uint32    `orm:"column(classic_prize)" json:"classic_prize,string" description:"经典奖金"`
	MethodId     uint32    `orm:"column(method_id)" json:"method_id,string" description:"玩法ID"`
	MethodName   string    `orm:"column(method_name);size(20);null" json:"method_name,string" description:"玩法名"`
	Level        uint8     `orm:"column(level)" json:"level,string" description:"奖级"`
	Probability  float64   `orm:"column(probability);digits(11);decimals(11)" json:"probability,string" description:"中奖率"`
	Prize        float64   `orm:"column(prize);digits(10);decimals(2)" json:"prize,string" description:"奖金"`
	FullPrize    float64   `orm:"column(full_prize);digits(16);decimals(6)" json:"full_prize,string" description:"全额奖金(返奖率为100%)"`
	CreatedAt    time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *PrizeDetails) TableName() string {
	return "prize_details"
}

func init() {
	orm.RegisterModel(new(PrizeDetails))
}

// AddPrizeDetails insert a new PrizeDetails into database and returns
// last inserted Id on success.
func AddPrizeDetails(m *PrizeDetails) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPrizeDetails mutil insert a new PrizeDetails into database
func AddMultiPrizeDetails(mlist []*PrizeDetails) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPrizeDetailsById retrieves PrizeDetails by Id. Returns error if
// Id doesn't exist
func GetPrizeDetailsById(id int) (v *PrizeDetails, err error) {
	o := orm.NewOrm()
	v = &PrizeDetails{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPrizeDetails retrieves all PrizeDetails matches certain condition. Returns empty list if
// no records exist
func GetAllPrizeDetails(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []PrizeDetails, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PrizeDetails))
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

// UpdatePrizeDetails updates PrizeDetails by Id and returns error if
// the record to be updated doesn't exist
func UpdatePrizeDetailsById(m *PrizeDetails) (err error) {
	o := orm.NewOrm()
	v := PrizeDetails{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePrizeDetails deletes PrizeDetails by Id and returns error if
// the record to be deleted doesn't exist
func DeletePrizeDetails(id int) (err error) {
	o := orm.NewOrm()
	v := PrizeDetails{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PrizeDetails{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
