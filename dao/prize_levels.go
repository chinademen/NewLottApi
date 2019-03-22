package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PrizeLevels struct {
	Id            int       `orm:"column(id);auto" json:"id,string"`
	LotteryTypeId uint8     `orm:"column(lottery_type_id)" json:"lottery_type_id,string" description:"彩票类型"`
	SeriesId      uint8     `orm:"column(series_id);null" json:"series_id,string" description:"系列ＩＤ"`
	BasicMethodId uint32    `orm:"column(basic_method_id)" json:"basic_method_id,string" description:"基础玩法ＩＤ"`
	BasicMethod   string    `orm:"column(basic_method);size(20);null" json:"basic_method,string" description:"基础玩法名"`
	Level         uint8     `orm:"column(level)" json:"level,string" description:"奖级"`
	OfficalPrize  float64   `orm:"column(offical_prize);null;digits(10);decimals(2)" json:"offical_prize,string"`
	MaxWinCount   uint      `orm:"column(max_win_count);null" json:"max_win_count,string"`
	Probability   float64   `orm:"column(probability);digits(16);decimals(16)" json:"probability,string" description:"中奖率"`
	FullPrize     float64   `orm:"column(full_prize);null;digits(16);decimals(6)" json:"full_prize,string" description:"全额奖金(返奖率为100%)"`
	MaxPrize      float64   `orm:"column(max_prize);digits(10);decimals(2)" json:"max_prize,string" description:"最高奖金"`
	MaxGroup      uint32    `orm:"column(max_group)" json:"max_group,string" description:"最高奖金组"`
	MinWater      float64   `orm:"column(min_water);digits(4);decimals(4)" json:"min_water,string" description:"最小利润率"`
	Rule          string    `orm:"column(rule);size(50)" json:"rule,string" description:"规则"`
	PrizeAllcount uint      `orm:"column(prize_allcount)" json:"prize_allcount,string"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *PrizeLevels) TableName() string {
	return "prize_levels"
}

func init() {
	orm.RegisterModel(new(PrizeLevels))
}

// AddPrizeLevels insert a new PrizeLevels into database and returns
// last inserted Id on success.
func AddPrizeLevels(m *PrizeLevels) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPrizeLevels mutil insert a new PrizeLevels into database
func AddMultiPrizeLevels(mlist []*PrizeLevels) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPrizeLevelsById retrieves PrizeLevels by Id. Returns error if
// Id doesn't exist
func GetPrizeLevelsById(id int) (v *PrizeLevels, err error) {
	o := orm.NewOrm()
	v = &PrizeLevels{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPrizeLevels retrieves all PrizeLevels matches certain condition. Returns empty list if
// no records exist
func GetAllPrizeLevels(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []PrizeLevels, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PrizeLevels))
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

// UpdatePrizeLevels updates PrizeLevels by Id and returns error if
// the record to be updated doesn't exist
func UpdatePrizeLevelsById(m *PrizeLevels) (err error) {
	o := orm.NewOrm()
	v := PrizeLevels{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePrizeLevels deletes PrizeLevels by Id and returns error if
// the record to be deleted doesn't exist
func DeletePrizeLevels(id int) (err error) {
	o := orm.NewOrm()
	v := PrizeLevels{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PrizeLevels{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
