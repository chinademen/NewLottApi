package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type BasicMethods struct {
	Id             int     `orm:"column(id);auto" json:"id,string"`
	LotteryType    uint8   `orm:"column(lottery_type)" json:"lottery_type,string" description:"彩票类型: 1-数字排列类型 2-乐透类型"`
	SeriesId       uint8   `orm:"column(series_id);null" json:"series_id,string" description:"系列ID"`
	Type           uint    `orm:"column(type);null" json:"type,string" description:"类型 来自于method_types"`
	Name           string  `orm:"column(name);size(20)" json:"name,string"`
	WnFunction     string  `orm:"column(wn_function);size(64);null" description:"计奖方法标记" json:"wn_function,string"`
	Sequencing     int8    `orm:"column(sequencing)" json:"sequencing,string" description:"定位"`
	DigitalCount   uint8   `orm:"column(digital_count)" json:"digital_count,string" description:"星"`
	UniqueCount    int8    `orm:"column(unique_count);null" json:"unique_count,string" description:"去重后的数字个数"`
	MaxRepeatTime  int8    `orm:"column(max_repeat_time);null" json:"max_repeat_time,string" description:"重号的最大重复次数"`
	MinRepeatTime  int8    `orm:"column(min_repeat_time);null" json:"min_repeat_time,string" description:"最小重复次数"`
	Span           uint8   `orm:"column(span);null" json:"span,string" description:"跨度"`
	MinSpan        uint8   `orm:"column(min_span);null" json:"min_span,string" description:"最小跨度"`
	ChooseCount    uint8   `orm:"column(choose_count);null" json:"choose_count,string" description:"组合数字个数"`
	MinChooseCount uint8   `orm:"column(min_choose_count);null" json:"min_choose_count,string" description:"最小组合数"`
	SpecialCount   uint8   `orm:"column(special_count);null" json:"special_count,string" description:"特号个数"`
	FixedNumber    int8    `orm:"column(fixed_number);null" json:"fixed_number,string" description:"固定号码"`
	Price          uint16  `orm:"column(price)" json:"price,string" description:"单价"`
	BuyLength      uint8   `orm:"column(buy_length)" json:"buy_length,string" description:"投注码长度"`
	WnLength       uint8   `orm:"column(wn_length)" json:"wn_length,string" description:"中奖号码长度"`
	WnCount        uint    `orm:"column(wn_count)" json:"wn_count,string" description:"中奖号码数量"`
	ValidNums      string  `orm:"column(valid_nums);size(50)" json:"valid_nums,string" description:"合法数字"`
	Rule           string  `orm:"column(rule);size(50)" json:"rule,string" description:"规则"`
	AllCount       uint64  `orm:"column(all_count)" json:"all_count,string" description:"投注码数量"`
	FullPrize      float64 `orm:"column(full_prize);null;digits(10);decimals(2)" json:"full_prize,string" description:"全额奖金(返奖率为100%)"`
	BetRule        string  `orm:"column(bet_rule);size(1024);null" json:"bet_rule,string" description:"玩法规则说明"`
	BonusNote      string  `orm:"column(bonus_note);size(1024);null" json:"bonus_note,string" description:"中奖说明"`
	Status         uint8   `orm:"column(status)" json:"status,string" json:"status,string"`
	CreatedAt      string  `orm:"column(created_at);type(datetime);null" json:"created_at,string"`
	UpdatedAt      string  `orm:"column(updated_at);type(datetime);null" json:"updated_at,string"`
}

func (t *BasicMethods) TableName() string {
	return "basic_methods"
}

func init() {
	orm.RegisterModel(new(BasicMethods))
}

// AddBasicMethods insert a new BasicMethods into database and returns
// last inserted Id on success.
func AddBasicMethods(m *BasicMethods) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiBasicMethods mutil insert a new BasicMethods into database
func AddMultiBasicMethods(mlist []*BasicMethods) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetBasicMethodsById retrieves BasicMethods by Id. Returns error if
// Id doesn't exist
func GetBasicMethodsById(id int) (v *BasicMethods, err error) {
	o := orm.NewOrm()
	v = &BasicMethods{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBasicMethods retrieves all BasicMethods matches certain condition. Returns empty list if
// no records exist
func GetAllBasicMethods(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []BasicMethods, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(BasicMethods))
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

// UpdateBasicMethods updates BasicMethods by Id and returns error if
// the record to be updated doesn't exist
func UpdateBasicMethodsById(m *BasicMethods) (err error) {
	o := orm.NewOrm()
	v := BasicMethods{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBasicMethods deletes BasicMethods by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBasicMethods(id int) (err error) {
	o := orm.NewOrm()
	v := BasicMethods{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BasicMethods{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
