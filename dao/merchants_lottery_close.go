package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type MerchantsLotteryClose struct {
	Id              int       `orm:"column(id);auto" json:"id,string"`
	MerchantId      int       `orm:"column(merchant_id);null" json:"merchant_id,string" description:"商户id"`
	LotteryIds      string    `orm:"column(lottery_ids);size(500);null" json:"lottery_ids,string" description:"已关闭的彩种"`
	WayGroupIds     string    `orm:"column(way_group_ids);size(500);null" json:"way_group_ids,string" description:"已关闭的玩法组"`
	SeriesMethodIds string    `orm:"column(series_method_ids);size(500);null" json:"series_method_ids,string" description:"已关闭的玩法"`
	SeriesWayIds    string    `orm:"column(series_way_ids);size(500);null" json:"series_way_ids,string" description:"已关闭的投注方式"`
	CreatedAt       time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *MerchantsLotteryClose) TableName() string {
	return "merchants_lottery_close"
}

func init() {
	orm.RegisterModel(new(MerchantsLotteryClose))
}

// AddMerchantsLotteryClose insert a new MerchantsLotteryClose into database and returns
// last inserted Id on success.
func AddMerchantsLotteryClose(m *MerchantsLotteryClose) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiMerchantsLotteryClose mutil insert a new MerchantsLotteryClose into database
func AddMultiMerchantsLotteryClose(mlist []*MerchantsLotteryClose) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetMerchantsLotteryCloseById retrieves MerchantsLotteryClose by Id. Returns error if
// Id doesn't exist
func GetMerchantsLotteryCloseById(id int) (v *MerchantsLotteryClose, err error) {
	o := orm.NewOrm()
	v = &MerchantsLotteryClose{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMerchantsLotteryClose retrieves all MerchantsLotteryClose matches certain condition. Returns empty list if
// no records exist
func GetAllMerchantsLotteryClose(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []MerchantsLotteryClose, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MerchantsLotteryClose))
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

// UpdateMerchantsLotteryClose updates MerchantsLotteryClose by Id and returns error if
// the record to be updated doesn't exist
func UpdateMerchantsLotteryCloseById(m *MerchantsLotteryClose) (err error) {
	o := orm.NewOrm()
	v := MerchantsLotteryClose{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMerchantsLotteryClose deletes MerchantsLotteryClose by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMerchantsLotteryClose(id int) (err error) {
	o := orm.NewOrm()
	v := MerchantsLotteryClose{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MerchantsLotteryClose{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
