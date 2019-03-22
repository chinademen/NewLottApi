package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Traces struct {
	Id               int       `orm:"column(id);auto" json:"id,string"`
	MerchantId       int       `orm:"column(merchant_id)" json:"merchant_id,string" description:"商户id"`
	TerminalId       uint8     `orm:"column(terminal_id)" json:"terminal_id,string" description:"终端id"`
	SerialNumber     string    `orm:"column(serial_number);size(32)" json:"serial_number,string" description:"追号编码"`
	UserId           int       `orm:"column(user_id)" json:"user_id,string" description:"用户id"`
	Username         string    `orm:"column(username);size(16)" json:"username,string" description:"用户名"`
	IsTester         int8      `orm:"column(is_tester);null" json:"is_tester,string" description:"是否测试，1=测试"`
	AccountId        int       `orm:"column(account_id);size(16)" json:"account_id,string" description:"账户ID"`
	PrizeGroup       string    `orm:"column(prize_group);size(20);null" json:"prize_group,string" description:"奖金组"`
	PrizeSet         string    `orm:"column(prize_set);size(1024);null" json:"prize_set,string" description:"奖金设置"`
	TotalIssues      uint32    `orm:"column(total_issues)" json:"total_issues,string" description:"总期数"`
	FinishedIssues   uint32    `orm:"column(finished_issues)" json:"finished_issues,string" description:"完成期数"`
	CanceledIssues   uint32    `orm:"column(canceled_issues)" json:"canceled_issues,string" description:"已取消期"`
	StopOnWon        int8      `orm:"column(stop_on_won)" json:"stop_on_won,string" description:"中奖即停"`
	LotteryId        uint8     `orm:"column(lottery_id)" json:"lottery_id,string" description:"彩种ID"`
	Title            string    `orm:"column(title);size(100)" json:"title,string" description:"玩法"`
	Position         string    `orm:"column(position);size(10);null" json:"position,string" description:"位置"`
	WayId            int       `orm:"column(way_id)" json:"way_id,string" description:"投注方式id"`
	BetNumber        string    `orm:"column(bet_number)" json:"bet_number,string" description:"投注号码"`
	WayTotalCount    int       `orm:"column(way_total_count);null" json:"way_total_count,string" description:"总注数"`
	SingleCount      int       `orm:"column(single_count);null" json:"single_count,string" description:"投注注数"`
	BetRate          float64   `orm:"column(bet_rate);null" json:"bet_rate,string" description:"投注比例"`
	DisplayBetNumber string    `orm:"column(display_bet_number);null" json:"display_bet_number,string" description:"显示的投注号码"`
	StartIssue       string    `orm:"column(start_issue);size(15)" json:"start_issue,string" description:"开始奖期"`
	WonIssue         string    `orm:"column(won_issue);size(15);null" json:"won_issue,string" description:"中奖奖期"`
	WonCount         uint32    `orm:"column(won_count)" json:"won_count,string" description:"中奖期数"`
	Prize            float64   `orm:"column(prize);digits(16);decimals(6)" json:"prize,string" description:"奖金"`
	Coefficient      float64   `orm:"column(coefficient);digits(4);decimals(3)" json:"coefficient,string" description:"奖金模式"`
	SingleAmount     float64   `orm:"column(single_amount);digits(14);decimals(4)" json:"single_amount,string" description:"单倍金额"`
	Amount           float64   `orm:"column(amount);digits(14);decimals(4)" json:"amount,string" description:"投注金额"`
	FinishedAmount   float64   `orm:"column(finished_amount);digits(14);decimals(4)" json:"finished_amount,string" description:"已完成金额"`
	CanceledAmount   float64   `orm:"column(canceled_amount);digits(14);decimals(4)" json:"canceled_amount,string" description:"已取消金额"`
	Status           int8      `orm:"column(status)" json:"status,string" description:"状态"`
	Ip               string    `orm:"column(ip);size(15);null" json:"ip,string" description:"用户ip"`
	ProxyIp          string    `orm:"column(proxy_ip);size(15);null" json:"proxy_ip,string" description:"代理服务器ip"`
	BetRecordId      uint64    `orm:"column(bet_record_id);null" json:"bet_record_id,string" description:"原始记录id"`
	BoughtAt         string    `orm:"column(bought_at);type(datetime);null" json:"bought_at,string" description:"投注时间"`
	CanceledAt       time.Time `orm:"column(canceled_at);type(datetime);null" description:"取消时间"`
	StopedAt         time.Time `orm:"column(stoped_at);type(datetime);null" description:"中止时间"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *Traces) TableName() string {
	return "traces"
}

func init() {
	orm.RegisterModel(new(Traces))
}

// AddTraces insert a new Traces into database and returns
// last inserted Id on success.
func AddTraces(o orm.Ormer, m *Traces) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// AddMultiTraces mutil insert a new Traces into database
func AddMultiTraces(mlist []*Traces) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetTracesById retrieves Traces by Id. Returns error if
// Id doesn't exist
func GetTracesById(id int) (v *Traces, err error) {
	o := orm.NewOrm()
	v = &Traces{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTraces retrieves all Traces matches certain condition. Returns empty list if
// no records exist
func GetAllTraces(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Traces, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Traces))
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

// UpdateTraces updates Traces by Id and returns error if
// the record to be updated doesn't exist
func UpdateTracesById(o orm.Ormer, m *Traces) (err error) {
	v := Traces{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTraces deletes Traces by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTraces(id int) (err error) {
	o := orm.NewOrm()
	v := Traces{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Traces{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
