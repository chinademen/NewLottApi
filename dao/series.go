package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Series struct {
	Id                int     `orm:"column(id);auto"`
	Type              uint8   `orm:"column(type);null" description:"类型 1-数字排列类型 2-乐透类型 "`
	LottoType         uint8   `orm:"column(lotto_type);null" description:"乐透类型 1-单区乐透 2-双区乐透"`
	Name              string  `orm:"column(name);size(20)" description:"名称"`
	Identifier        string  `orm:"column(identifier);size(10);null" description:"系列标识符"`
	SortWinningNumber int8    `orm:"column(sort_winning_number);null" description:"中奖号码排序 0-不排序  1-排序"`
	BuyLength         string  `orm:"column(buy_length);size(10)" description:"投注码长度"`
	WnLength          string  `orm:"column(wn_length);size(10)" description:"中奖号码长度"`
	DigitalCount      string  `orm:"column(digital_count);size(10)" description:"星"`
	ClassicAmount     uint    `orm:"column(classic_amount);null" description:"直选投注码数量"`
	GroupType         int8    `orm:"column(group_type);null" description:"返点类型 1-奖金组返点 2-百分比返点"`
	MaxPercentGroup   int     `orm:"column(max_percent_group);null"`
	MaxPrizeGroup     uint    `orm:"column(max_prize_group);null" description:"最高奖金组"`
	MaxRealGroup      uint    `orm:"column(max_real_group)" description:"最高实际奖金组"`
	MaxBetGroup       int     `orm:"column(max_bet_group);null" description:"最高投注奖金组"`
	ValidNums         string  `orm:"column(valid_nums);size(300);null" description:"合法数字"`
	OfficalPrizeRate  float64 `orm:"column(offical_prize_rate);null;digits(2);decimals(2)" description:"官方中奖率"`
	DefaultWayId      uint    `orm:"column(default_way_id);null" description:"默认投注方式"`
	LinkTo            uint8   `orm:"column(link_to);null" description:"关联至系列"`
	Lotteries         string  `orm:"column(lotteries);size(200);null" description:"彩种"`
	BonusEnabled      int8    `orm:"column(bonus_enabled)" description:"是否派奖"`
}

func (t *Series) TableName() string {
	return "series"
}

func init() {
	orm.RegisterModel(new(Series))
}

// AddSeries insert a new Series into database and returns
// last inserted Id on success.
func AddSeries(m *Series) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSeries mutil insert a new Series into database
func AddMultiSeries(mlist []*Series) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSeriesById retrieves Series by Id. Returns error if
// Id doesn't exist
func GetSeriesById(id int) (v *Series, err error) {
	o := orm.NewOrm()
	v = &Series{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//SQL
func Query(sSql string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	o.Raw(sSql, 2).Values(&maps)
	return maps
}

// GetAllSeries retrieves all Series matches certain condition. Returns empty list if
// no records exist
func GetAllSeries(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*Series, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Series))
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

// UpdateSeries updates Series by Id and returns error if
// the record to be updated doesn't exist
func UpdateSeriesById(m *Series) (err error) {
	o := orm.NewOrm()
	v := Series{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSeries deletes Series by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSeries(id int) (err error) {
	o := orm.NewOrm()
	v := Series{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Series{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
