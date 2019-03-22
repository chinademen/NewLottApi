package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Projects struct {
	Id               int       `orm:"column(id);auto"`
	MerchantId       int       `orm:"column(merchant_id)" description:"追号任务ID"`
	TerminalId       int       `orm:"column(terminal_id)" description:"终端id"`
	SerialNumber     string    `orm:"column(serial_number);size(32);null" description:"注单编号"`
	TraceId          int64     `orm:"column(trace_id);null" description:"追号任务ID"`
	UserId           int64     `orm:"column(user_id)" description:"用户ID"`
	Username         string    `orm:"column(username);size(32)" description:"用户名"`
	IsTester         int8      `orm:"column(is_tester);null" description:"1=测试"`
	AccountId        int64     `orm:"column(account_id)" description:"账户id"`
	PrizeGroup       string    `orm:"column(prize_group);size(20);null" description:"投注时的奖金组"`
	LotteryId        uint8     `orm:"column(lottery_id)" description:"彩种id"`
	Issue            string    `orm:"column(issue);size(15)" description:"奖期"`
	EndTime          int       `orm:"column(end_time);size(16);null" description:"奖期截至时间"`
	WayId            int       `orm:"column(way_id)" description:"投注方式id"`
	Title            string    `orm:"column(title);size(100)" description:"标题"`
	Position         string    `orm:"column(position);size(10);null" description:"位置"`
	BetNumber        string    `orm:"column(bet_number)" description:"投注号码"`
	WayTotalCount    uint64    `orm:"column(way_total_count);null" description:"总注单数"`
	SingleCount      int       `orm:"column(single_count);null" description:"投注方式总注注数"`
	BetRate          float32   `orm:"column(bet_rate);null" description:"投注比例"`
	DisplayBetNumber string    `orm:"column(display_bet_number);null" description:"显示出来的投注号码"`
	Multiple         int       `orm:"column(multiple)" description:"倍数"`
	SingleAmount     float64   `orm:"column(single_amount);digits(14);decimals(4)" description:"单注金额"`
	Amount           float64   `orm:"column(amount);digits(14);decimals(4)" description:"总金额"`
	WinningNumber    string    `orm:"column(winning_number);size(60);null" description:"开奖号码"`
	Prize            float64   `orm:"column(prize);null;digits(14);decimals(4)" description:"奖金"`
	Status           int8      `orm:"column(status)" description:"0: 正常；1：已撤销；2：未中奖；3：已中奖；4：已派奖；5：系统撤销（通过redis和帐变表去重复）"`
	StatusPrize      int8      `orm:"column(status_prize)" description:"0: 正常；1：已撤销；2：未中奖；3：已中奖；4：已派奖；5：系统撤销（通过redis和帐变表去重复）"`
	StatusSync       int8      `orm:"column(status_sync);null" description:"注单同步，1=同步"`
	PrizeSet         string    `orm:"column(prize_set);size(1024);null" description:"奖金设置"`
	SingleWonCount   int       `orm:"column(single_won_count);null" description:"单倍中奖注数"`
	WonCount         int       `orm:"column(won_count);null" description:"中奖注数"`
	WonData          string    `orm:"column(won_data);size(10240);null" description:"中奖详情"`
	Ip               string    `orm:"column(ip);size(15)" description:"ip"`
	ProxyIp          string    `orm:"column(proxy_ip);size(15)" description:"代理ip"`
	BetRecordId      uint64    `orm:"column(bet_record_id);null" description:"投注单id"`
	CanceledBy       string    `orm:"column(canceled_by);size(16);null" description:"执行撤销的管理员"`
	BoughtAt         string    `orm:"column(bought_at);type(datetime)"`
	CanceledAt       time.Time `orm:"column(canceled_at);type(datetime);null" description:"撤单时间"`
	PrizeSentAt      time.Time `orm:"column(prize_sent_at);type(datetime);null" description:"派奖时间"`
	BoughtTime       int       `orm:"column(bought_time);null" description:"投注时间"`
	BetCommitTime    int       `orm:"column(bet_commit_time);null" description:"投注提交至数据库时间"`
	Coefficient      float64   `orm:"column(coefficient);null;digits(4);decimals(3)" description:"投注金额模式(1=2元;0.5=1元;0.1=2角0.05=1角;0.01=2分;0.001=2厘)"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *Projects) TableName() string {
	return "projects"
}

func init() {
	orm.RegisterModel(new(Projects))
}

// AddProjects insert a new Projects into database and returns
// last inserted Id on success.
func AddProjects(o orm.Ormer, m *Projects) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// AddMultiProjects mutil insert a new Projects into database
func AddMultiProjects(mlist []*Projects) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetProjectsById retrieves Projects by Id. Returns error if
// Id doesn't exist
func GetProjectsById(id int) (v *Projects, err error) {
	o := orm.NewOrm()
	v = &Projects{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllProjects retrieves all Projects matches certain condition. Returns empty list if
// no records exist
func GetAllProjects(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Projects, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Projects))
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

// UpdateProjects updates Projects by Id and returns error if
// the record to be updated doesn't exist
func UpdateProjectsById(m *Projects) (err error) {
	o := orm.NewOrm()
	v := Projects{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProjects deletes Projects by Id and returns error if
// the record to be deleted doesn't exist
func DeleteProjects(id int) (err error) {
	o := orm.NewOrm()
	v := Projects{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Projects{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
