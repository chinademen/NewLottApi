package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Sources struct {
	Id          int    `orm:"column(id);auto" description:"主键"`
	SourceNo    string `orm:"column(source_no);size(10);null"`
	LotteryId   uint8  `orm:"column(lottery_id);null" description:"彩种ID"`
	Title       string `orm:"column(title);size(60);null" description:"仅用户后台管理页面显示的号源可读名称标识"`
	Type        uint8  `orm:"column(type);null" description:"抓取类型   1: CURL抓取  2:配置连接抓取"`
	OriginalUrl string `orm:"column(original_url);size(255);null" description:"仅用户后台管理页面显示的号源原URL可读名称标识"`
	ConfigUrl   string `orm:"column(config_url);null" description:"type0:抓取方式配置
type1:数据库连接配置"`
	ConfigUrlOld     string    `orm:"column(config_url_old);null"`
	ConfigContent    string    `orm:"column(config_content);null" description:"获取内容配置"`
	ConfigContentOld string    `orm:"column(config_content_old);null"`
	ConfigGrab       string    `orm:"column(config_grab);null" description:"抓取参数配置"`
	ConfigGrabOld    string    `orm:"column(config_grab_old);null"`
	ConfigCode       string    `orm:"column(config_code);null" description:"开奖号码配置"`
	ConfigCodeOld    string    `orm:"column(config_code_old);null"`
	Rank             uint      `orm:"column(rank);null" description:"权重/评分"`
	IsTest           int8      `orm:"column(is_test)"`
	UserId           uint      `orm:"column(user_id);null"`
	IsEnabled        int8      `orm:"column(is_enabled);null"`
	IsDebug          int8      `orm:"column(is_debug);null"`
	Remark           string    `orm:"column(remark);null" description:"仅用户后台管理页面显示的备注信息"`
	StartTime1       int       `orm:"column(start_time1);null" description:"第一阶段开始抓取的时间距离开奖时间的秒数"`
	StartTime1Freq   int       `orm:"column(start_time1_freq);null" description:"第一阶段开始抓取的频率"`
	StartTime2       int       `orm:"column(start_time2);null" description:"第二阶段开始抓取的时间距离开奖时间的秒数"`
	StartTime2Freq   int       `orm:"column(start_time2_freq);null" description:"第二阶段开始抓取的频率"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *Sources) TableName() string {
	return "sources"
}

func init() {
	orm.RegisterModel(new(Sources))
}

// AddSources insert a new Sources into database and returns
// last inserted Id on success.
func AddSources(m *Sources) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSources mutil insert a new Sources into database
func AddMultiSources(mlist []*Sources) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSourcesById retrieves Sources by Id. Returns error if
// Id doesn't exist
func GetSourcesById(id int) (v *Sources, err error) {
	o := orm.NewOrm()
	v = &Sources{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSources retrieves all Sources matches certain condition. Returns empty list if
// no records exist
func GetAllSources(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Sources, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Sources))
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

// UpdateSources updates Sources by Id and returns error if
// the record to be updated doesn't exist
func UpdateSourcesById(m *Sources) (err error) {
	o := orm.NewOrm()
	v := Sources{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSources deletes Sources by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSources(id int) (err error) {
	o := orm.NewOrm()
	v := Sources{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Sources{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
