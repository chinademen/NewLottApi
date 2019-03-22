package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type SeriesWays struct {
	Id                int    `orm:"column(id);auto" json:"id,string"`
	LotteryType       uint8  `orm:"column(lottery_type);null" json:"lottery_type,string" description:"彩票类型"`
	SeriesId          uint8  `orm:"column(series_id)" json:"series_id,string" description:"系列ID"`
	Name              string `orm:"column(name);size(30)" json:"name,string" description:"名称"`
	ShortName         string `orm:"column(short_name);size(30);null" json:"short_name,string" description:"短名称"`
	NeedSplit         uint8  `orm:"column(need_split)" json:"need_split,string" description:"任选投注方方式标记 1-任选 0-非任选"`
	WayMaps           string `orm:"column(way_maps);size(1024)" json:"way_maps,string"`
	SeriesWayMethodId uint   `orm:"column(series_way_method_id);null" json:"series_way_method_id,string" description:"系列投注方式与玩法关系"`
	BasicWayId        uint8  `orm:"column(basic_way_id)" json:"basic_way_id,string" description:"基础投注方式"`
	BasicMethods      string `orm:"column(basic_methods);size(200)" json:"basic_methods,string" description:"基础玩法"`
	SeriesMethods     string `orm:"column(series_methods);size(200)" json:"series_methods,string" description:"系列玩法"`
	WayFunction       string `orm:"column(way_function);size(64);null" json:"way_function,string" description:"任选投注方方式标记 下投注号码处理规则标记"`
	WnFunction        string `orm:"column(wn_function);size(64);null" json:"wn_function,string" description:"计奖方法标记"`
	DigitalCount      uint8  `orm:"column(digital_count);null" json:"digital_count,string" description:"星"`
	NeedPosition      uint8  `orm:"column(need_position)" json:"need_position,string"`
	Price             uint16 `orm:"column(price)" json:"price,string" description:"单价"`
	Offset            string `orm:"column(offset);size(100);null" json:"offset,string" description:"起始位"`
	Position          string `orm:"column(position);size(100);null" json:"position,string"`
	BuyLength         uint8  `orm:"column(buy_length);null" json:"buy_length,string" description:"投注码长度"`
	WnLength          uint8  `orm:"column(wn_length);null" json:"wn_length,string" description:"中奖号码长度"`
	WnCount           uint   `orm:"column(wn_count);null" json:"wn_count,string" description:"中奖号码数量"`
	AreaCount         int8   `orm:"column(area_count);null" json:"area_count,string" description:"投注码分区数"`
	AreaConfig        string `orm:"column(area_config);size(20);null" json:"area_config,string" description:"投注码分区配置"`
	ValidNums         string `orm:"column(valid_nums);size(50);null" json:"valid_nums,string" description:"合法数字"`
	Rule              string `orm:"column(rule);size(50);null" json:"rule,string" description:"规则"`
	AllCount          string `orm:"column(all_count);size(100)" json:"all_count,string" description:"投注码数量"`
	BonusNote         string `orm:"column(bonus_note);size(1024);null" json:"bonus_note,string" description:"中奖说明"`
	BetNote           string `orm:"column(bet_note);size(1024);null" json:"bet_note,string" description:"选号规则"`
	CreatedAt         string `orm:"column(created_at);type(datetime);null" json:"created_at,string"`
	UpdatedAt         string `orm:"column(updated_at);type(datetime);null" json:"updated_at,string"`
}

func (t *SeriesWays) TableName() string {
	return "series_ways"
}

func init() {
	orm.RegisterModel(new(SeriesWays))
}

// AddSeriesWays insert a new SeriesWays into database and returns
// last inserted Id on success.
func AddSeriesWays(m *SeriesWays) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSeriesWays mutil insert a new SeriesWays into database
func AddMultiSeriesWays(mlist []*SeriesWays) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSeriesWaysById retrieves SeriesWays by Id. Returns error if
// Id doesn't exist
func GetSeriesWaysById(id int) (v *SeriesWays, err error) {
	o := orm.NewOrm()
	v = &SeriesWays{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSeriesWays retrieves all SeriesWays matches certain condition. Returns empty list if
// no records exist
func GetAllSeriesWays(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []SeriesWays, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SeriesWays))
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

// UpdateSeriesWays updates SeriesWays by Id and returns error if
// the record to be updated doesn't exist
func UpdateSeriesWaysById(m *SeriesWays) (err error) {
	o := orm.NewOrm()
	v := SeriesWays{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSeriesWays deletes SeriesWays by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSeriesWays(id int) (err error) {
	o := orm.NewOrm()
	v := SeriesWays{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SeriesWays{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
