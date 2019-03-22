package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Issues struct {
	Id              int       `orm:"column(id);auto" json:"lottery_id,string"`
	LotteryId       uint8     `orm:"column(lottery_id)" json:"lottery_id,string" description:"彩种id"`
	Issue           string    `orm:"column(issue);size(15)" json:"issue,string" description:"奖期"`
	IssueRuleId     uint16    `orm:"column(issue_rule_id);null" json:"issue_rule_id,string" description:"对应的奖期规则id"`
	BeginTime       string    `orm:"column(begin_time);size(16);null" json:"begin_time,string" description:"开始销售时间"`
	EndTime         string    `orm:"column(end_time);size(16);null" json:"end_time,string" description:"销售截止时间"`
	EndTime2        string    `orm:"column(end_time2);type(datetime);null" json:"end_time2,encoded_at,string" description:"冗余销售截止时间格式"`
	OfficalTime     string    `orm:"column(offical_time);size(16);null" json:"offical_time,string" description:"官方开奖时间"`
	Cycle           string    `orm:"column(cycle);size(16);null" json:"cycle,string" description:"周期"`
	WnNumber        string    `orm:"column(wn_number);size(60)" json:"wn_number,string" description:"中奖号码"`
	AllowEncodeTime string    `orm:"column(allow_encode_time);size(16);null" json:"allow_encode_time,string" description:"允许录号时间"`
	EncoderId       string    `orm:"column(encoder_id);size(16);null" json:"encoder_id,string" description:"录号者id"`
	Encoder         string    `orm:"column(encoder);size(32);null" json:"encoder,string" description:"录号者名称"`
	EncodedAt       string    `orm:"column(encoded_at);type(datetime);null" json:"encoded_at,string" description:"录号时间"`
	Status          uint8     `orm:"column(status)" json:"status,string" description:"状态(3:已开奖1：等待开奖2：已经输入号码，等待审核4：号码已审核8：号码已取消开奖32：提前开奖A，获取到开奖号码的时间早于官方理论开奖时间64：提前开奖B，获取到开奖号码的时间早于销售截止时间)"`
	Locker          string    `orm:"column(locker);size(16);null" json:"status,string" json:"locker,string" description:"锁定者ID"`
	LockerFund      string    `orm:"column(locker_fund);size(16)" json:"status,string" json:"locker_fund,string" description:"锁定资金"`
	Tag             string    `orm:"column(tag);size(50);null" json:"status,string" json:"tag,string" description:"附加信息"`
	CalculatedAt    string    `orm:"column(calculated_at);type(datetime);null" json:"status,string" description:"计奖时间"`
	PrizeSentAt     string    `orm:"column(prize_sent_at);type(datetime);null" json:"datetime,string" description:"派奖时间"`
	CreatedAt       time.Time `orm:"column(created_at);type(datetime);null;auto_now_add" json:"created_at,string"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add" json:"updated_at,string"`
}

func (t *Issues) TableName() string {
	return "issues"
}

func init() {
	orm.RegisterModel(new(Issues))
}

// AddIssues insert a new Issues into database and returns
// last inserted Id on success.
func AddIssues(m *Issues) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiIssues mutil insert a new Issues into database
func AddMultiIssues(mlist []*Issues) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetIssuesById retrieves Issues by Id. Returns error if
// Id doesn't exist
func GetIssuesById(id int) (v *Issues, err error) {
	o := orm.NewOrm()
	v = &Issues{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllIssues retrieves all Issues matches certain condition. Returns empty list if
// no records exist
func GetAllIssues(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Issues, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Issues))
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

// UpdateIssues updates Issues by Id and returns error if
// the record to be updated doesn't exist
func UpdateIssuesById(m *Issues) (err error) {
	o := orm.NewOrm()
	v := Issues{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteIssues deletes Issues by Id and returns error if
// the record to be deleted doesn't exist
func DeleteIssues(id int) (err error) {
	o := orm.NewOrm()
	v := Issues{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Issues{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
