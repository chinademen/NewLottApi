package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PrjPrizeSets struct {
	Id                int       `orm:"column(id);auto" json:"id,string"`
	MerchantId        uint      `orm:"column(merchant_id);null" json:"merchant_id,string"`
	TraceId           uint64    `orm:"column(trace_id);null" json:"trace_id,string"`
	ProjectId         uint64    `orm:"column(project_id)" json:"project_id,string"`
	UserId            uint      `orm:"column(user_id)" json:"user_id,string"`
	AccountId         uint      `orm:"column(account_id)" json:"account_id,string"`
	Username          string    `orm:"column(username);size(16)" json:"username,string"`
	IsTester          int8      `orm:"column(is_tester);null" json:"is_tester,string"`
	UserForefatherIds string    `orm:"column(user_forefather_ids);size(1024)" json:"user_forefather_ids,string"`
	ProjectNo         string    `orm:"column(project_no);size(32)" json:"project_no,string"`
	LotteryId         uint8     `orm:"column(lottery_id)" json:"lottery_id,string"`
	Issue             string    `orm:"column(issue);size(15)" json:"issue,string"`
	WayId             uint      `orm:"column(way_id)" json:"way_id,string"`
	BasicMethodId     uint32    `orm:"column(basic_method_id)" json:"basic_method_id,string"`
	Level             uint8     `orm:"column(level)" json:"level,string"`
	Coefficient       float64   `orm:"column(coefficient);digits(4);decimals(3)" json:"coefficient,string"`
	AgentSets         string    `orm:"column(agent_sets);size(1024);null" json:"agent_sets,string"`
	PrizeSet          float64   `orm:"column(prize_set);digits(14);decimals(4)" json:"prize_set,string"`
	SingleWonCount    uint      `orm:"column(single_won_count);null" json:"single_won_count,string"`
	Multiple          uint      `orm:"column(multiple)" json:"multiple,string"`
	Prize             float64   `orm:"column(prize);digits(16);decimals(6)" json:"prize,string"`
	WonCount          uint      `orm:"column(won_count)" json:"won_count,string"`
	Status            uint8     `orm:"column(status)" json:"status,string"`
	Locked            uint64    `orm:"column(locked)" json:"locked,string"`
	SentAt            string    `orm:"column(sent_at);type(datetime);null" json:"sent_at,string"`
	CreatedAt         time.Time `orm:"column(created_at);type(datetime)"`
	UpdatedAt         time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *PrjPrizeSets) TableName() string {
	return "prj_prize_sets"
}

func init() {
	orm.RegisterModel(new(PrjPrizeSets))
}

// AddPrjPrizeSets insert a new PrjPrizeSets into database and returns
// last inserted Id on success.
func AddPrjPrizeSets(m *PrjPrizeSets) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiPrjPrizeSets mutil insert a new PrjPrizeSets into database
func AddMultiPrjPrizeSets(mlist []*PrjPrizeSets) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetPrjPrizeSetsById retrieves PrjPrizeSets by Id. Returns error if
// Id doesn't exist
func GetPrjPrizeSetsById(id int) (v *PrjPrizeSets, err error) {
	o := orm.NewOrm()
	v = &PrjPrizeSets{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPrjPrizeSets retrieves all PrjPrizeSets matches certain condition. Returns empty list if
// no records exist
func GetAllPrjPrizeSets(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []PrjPrizeSets, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PrjPrizeSets))
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

// UpdatePrjPrizeSets updates PrjPrizeSets by Id and returns error if
// the record to be updated doesn't exist
func UpdatePrjPrizeSetsById(m *PrjPrizeSets) (err error) {
	o := orm.NewOrm()
	v := PrjPrizeSets{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePrjPrizeSets deletes PrjPrizeSets by Id and returns error if
// the record to be deleted doesn't exist
func DeletePrjPrizeSets(id int) (err error) {
	o := orm.NewOrm()
	v := PrjPrizeSets{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PrjPrizeSets{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
