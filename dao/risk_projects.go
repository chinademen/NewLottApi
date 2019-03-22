package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type RiskProjects struct {
	Id               int       `orm:"column(id);auto"`
	MerchantId       uint      `orm:"column(merchant_id)" description:"商户ID"`
	ProjectId        uint64    `orm:"column(project_id)" description:"注单ID"`
	TerminalId       uint8     `orm:"column(terminal_id)" description:"终端ＩD"`
	SerialNumber     string    `orm:"column(serial_number);size(32);null" description:"序列号"`
	TraceId          uint      `orm:"column(trace_id);null" description:"追号任务ID"`
	UserId           uint      `orm:"column(user_id)" description:"用户ID"`
	Username         string    `orm:"column(username);size(32)" description:"用户名"`
	IsTester         int8      `orm:"column(is_tester);null" description:"是否测试"`
	AccountId        uint      `orm:"column(account_id)" description:"账户ＩＤ"`
	PrizeGroup       string    `orm:"column(prize_group);size(20);null" description:"奖金组"`
	LotteryId        uint8     `orm:"column(lottery_id)" description:"彩种ＩＤ"`
	Issue            string    `orm:"column(issue);size(15)" description:"奖期"`
	EndTime          uint      `orm:"column(end_time);null" description:"结束时间"`
	WayId            uint      `orm:"column(way_id)" description:"玩法"`
	Title            string    `orm:"column(title);size(100)"`
	Position         string    `orm:"column(position);size(10);null"`
	BetNumber        string    `orm:"column(bet_number)" description:"投注号码"`
	WayTotalCount    uint64    `orm:"column(way_total_count);null" description:"玩法总注数"`
	SingleCount      uint      `orm:"column(single_count);null" description:"注数"`
	BetRate          float32   `orm:"column(bet_rate);null" description:"胜率"`
	DisplayBetNumber string    `orm:"column(display_bet_number);null" description:"显示投注号码"`
	Multiple         uint      `orm:"column(multiple)" description:"倍数"`
	Coefficient      float64   `orm:"column(coefficient);digits(4);decimals(3)" description:"模式"`
	SingleAmount     float64   `orm:"column(single_amount);digits(14);decimals(4)" description:"单注金额"`
	Amount           float64   `orm:"column(amount);digits(14);decimals(4)" description:"总金额"`
	WinningNumber    string    `orm:"column(winning_number);size(60);null" description:"开奖号码"`
	Prize            float64   `orm:"column(prize);null;digits(14);decimals(4)" description:"奖金"`
	PrizeSaleRate    float64   `orm:"column(prize_sale_rate);null" description:"销售率"`
	Status           uint8     `orm:"column(status)" description:"0: 待审核；1：已审核；2：审核未通过"`
	Auditor          string    `orm:"column(auditor);size(16);null" description:"审核人"`
	AuditedAt        time.Time `orm:"column(audited_at);type(datetime);null" description:"审核时间"`
	RefuseReason     string    `orm:"column(refuse_reason);null" description:"拒绝原因"`
	PrizeSet         string    `orm:"column(prize_set);size(1024);null" description:"奖金设置"`
	SingleWonCount   uint      `orm:"column(single_won_count);null"`
	WonCount         uint      `orm:"column(won_count);null"`
	WonData          string    `orm:"column(won_data);size(10240);null"`
	Ip               string    `orm:"column(ip);size(15)" description:"ＩＰ"`
	ProxyIp          string    `orm:"column(proxy_ip);size(15)" description:"代理ＩＰ"`
	BetRecordId      uint64    `orm:"column(bet_record_id);null" description:"投注原始记录ＩＤ"`
	BoughtAt         time.Time `orm:"column(bought_at);type(datetime)" description:"投注时间"`
	CountedAt        time.Time `orm:"column(counted_at);type(datetime);null" description:"计奖时间"`
	BoughtTime       uint      `orm:"column(bought_time);null" description:"冗投注时间"`
	BetCommitTime    uint      `orm:"column(bet_commit_time);null" description:"投注提交到库时间"`
	CountedTime      uint      `orm:"column(counted_time);null" description:"冗余计奖时间"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *RiskProjects) TableName() string {
	return "risk_projects"
}

func init() {
	orm.RegisterModel(new(RiskProjects))
}

// AddRiskProjects insert a new RiskProjects into database and returns
// last inserted Id on success.
func AddRiskProjects(m *RiskProjects) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiRiskProjects mutil insert a new RiskProjects into database
func AddMultiRiskProjects(mlist []*RiskProjects) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetRiskProjectsById retrieves RiskProjects by Id. Returns error if
// Id doesn't exist
func GetRiskProjectsById(id int) (v *RiskProjects, err error) {
	o := orm.NewOrm()
	v = &RiskProjects{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRiskProjects retrieves all RiskProjects matches certain condition. Returns empty list if
// no records exist
func GetAllRiskProjects(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []RiskProjects, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RiskProjects))
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

// UpdateRiskProjects updates RiskProjects by Id and returns error if
// the record to be updated doesn't exist
func UpdateRiskProjectsById(m *RiskProjects) (err error) {
	o := orm.NewOrm()
	v := RiskProjects{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRiskProjects deletes RiskProjects by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRiskProjects(id int) (err error) {
	o := orm.NewOrm()
	v := RiskProjects{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RiskProjects{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
