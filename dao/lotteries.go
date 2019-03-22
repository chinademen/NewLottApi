package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Lotteries struct {
	Id                int    `orm:"column(id);auto" json:"id,string"`
	SeriesId          uint8  `orm:"column(series_id)" json:"series_id,string" description:"系列类型ＩＤ"`
	Name              string `orm:"column(name);size(20)" json:"name,string" description:"彩种名"`
	Type              uint8  `orm:"column(type)" json:"type,string" description:"类型：1-数字排列类型，2-乐透类型"`
	LottoType         uint8  `orm:"column(lotto_type);null" json:"lotto_type,string" description:"乐透类型 1-单区乐透类型 2-双区乐透类型"`
	IsSelf            uint8  `orm:"column(is_self)" json:"is_self,string" description:"是否自主彩种 0-否 1-是"`
	IsInstant         int8   `orm:"column(is_instant)" json:"is_instant,string" description:"是否即时彩 0-否 1-是"`
	HighFrequency     uint8  `orm:"column(high_frequency)" json:"high_frequency,string" description:"是否高频 0-低频 1-高频 "`
	SortWinningNumber uint8  `orm:"column(sort_winning_number);null" json:"sort_winning_number,string" description:"中奖号码排序 1-排序  0-不排序"`
	ValidNums         string `orm:"column(valid_nums);size(300);null" json:"valid_nums,string" description:"可选的号码范围，多区时用+分开"`
	BuyLength         string `orm:"column(buy_length);size(10)" json:"buy_length,string" description:"投注码长度"`
	WnLength          string `orm:"column(wn_length);size(10)" json:"wn_length,string" description:"中奖号码长度"`
	Identifier        string `orm:"column(identifier);size(10)" json:"identifier,string" description:"彩种标识符"`
	Days              uint8  `orm:"column(days)" json:"days,string" description:"开奖日"`
	IssueOverMidnight uint8  `orm:"column(issue_over_midnight)" json:"issue_over_midnight,string" description:"奖期跨越零点"`
	IssueFormat       string `orm:"column(issue_format);size(16)" json:"issue_format,string" description:"奖期格式"`
	BetTemplate       string `orm:"column(bet_template);size(20);null" json:"bet_template,string" description:"投注模板"`
	BeginTime         string `orm:"column(begin_time);type(time);null" json:"begin_time,string" description:"开售时间"`
	EndTime           string `orm:"column(end_time);type(time);null" json:"end_time,string" description:"截止时间"`
	Sequence          uint16 `orm:"column(sequence)" json:"sequence,string" description:"排序值"`
	Status            uint8  `orm:"column(status)" json:"status,string" description:"状态"`
	NeedDraw          uint8  `orm:"column(need_draw)" json:"need_draw,string" description:"是否抓号"`
	DailyIssueCount   uint   `orm:"column(daily_issue_count);null" json:"daily_issue_count,string" description:"每日奖期数"`
	TraceIssueCount   uint   `orm:"column(trace_issue_count);null" json:"trace_issue_count,string" description:"追号奖期数"`
	MaxBetGroup       uint   `orm:"column(max_bet_group)" json:"max_bet_group,string" description:"最高投注奖金组"`
	SeriesWays        string `orm:"column(series_ways);size(10140);null" json:"series_ways,string" description:"系列投注方式"`
	CreatedAt         string `orm:"column(created_at);type(datetime);null" json:"created_at,string"`
	UpdatedAt         string `orm:"column(updated_at);type(datetime);null" json:"updated_at,string"`
}

func (t *Lotteries) TableName() string {
	return "lotteries"
}

func init() {
	orm.RegisterModel(new(Lotteries))
}

// AddLotteries insert a new Lotteries into database and returns
// last inserted Id on success.
func AddLotteries(m *Lotteries) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiLotteries mutil insert a new Lotteries into database
func AddMultiLotteries(mlist []*Lotteries) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetLotteriesById retrieves Lotteries by Id. Returns error if
// Id doesn't exist
func GetLotteriesById(id int) (v *Lotteries, err error) {
	o := orm.NewOrm()
	v = &Lotteries{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, nil
}

// GetLotteriesByName retrieves Lotteries by Name. Returns error if
// Name doesn't exist
func GetLotteriesByName(sName string) (v *Lotteries, err error) {
	o := orm.NewOrm()
	v = &Lotteries{Name: sName}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLotteries retrieves all Lotteries matches certain condition. Returns empty list if
// no records exist
func GetAllLotteries(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*Lotteries, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lotteries))
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

// UpdateLotteries updates Lotteries by Id and returns error if
// the record to be updated doesn't exist
func UpdateLotteriesById(m *Lotteries) (err error) {
	o := orm.NewOrm()
	v := Lotteries{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLotteries deletes Lotteries by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLotteries(id int) (err error) {
	o := orm.NewOrm()
	v := Lotteries{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Lotteries{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
