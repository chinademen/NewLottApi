package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserPrizeSets struct {
	Id           int       `orm:"column(id);auto" json:"id,string"`
	MerchantId   uint      `orm:"column(merchant_id)" json:"merchant_id,string" description:"商户id"`
	UserId       uint64    `orm:"column(user_id)" json:"user_id,string" description:"用户ID"`
	Username     string    `orm:"column(username);size(16)" json:"username,string" description:"用户名"`
	SeriesId     uint8     `orm:"column(series_id);null" json:"series_id,string" description:"系列id"`
	LotteryId    uint8     `orm:"column(lottery_id)" json:"lottery_id,string" description:"彩种"`
	GroupId      uint      `orm:"column(group_id)" json:"group_id,string" description:"组"`
	PrizeGroup   string    `orm:"column(prize_group);size(20);null" json:"prize_group,string" description:"奖金组"`
	ClassicPrize uint32    `orm:"column(classic_prize)" json:"classic_prize,string" description:"经典奖金"`
	Valid        uint8     `orm:"column(valid)" json:"valid,string" description:"有效"`
	IsAgent      uint8     `orm:"column(is_agent)" json:"is_agent,string" description:"是否代理0: 普通用户, 1: 代理"`
	CreatedAt    time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(datetime);null"`
}

func (t *UserPrizeSets) TableName() string {
	return "user_prize_sets"
}

func init() {
	orm.RegisterModel(new(UserPrizeSets))
}

// AddUserPrizeSets insert a new UserPrizeSets into database and returns
// last inserted Id on success.
func AddUserPrizeSets(m *UserPrizeSets) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiUserPrizeSets mutil insert a new UserPrizeSets into database
func AddMultiUserPrizeSets(mlist []*UserPrizeSets) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetUserPrizeSetsById retrieves UserPrizeSets by Id. Returns error if
// Id doesn't exist
func GetUserPrizeSetsById(id int) (v *UserPrizeSets, err error) {
	o := orm.NewOrm()
	v = &UserPrizeSets{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserPrizeSets retrieves all UserPrizeSets matches certain condition. Returns empty list if
// no records exist
func GetAllUserPrizeSets(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []UserPrizeSets, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserPrizeSets))
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

// UpdateUserPrizeSets updates UserPrizeSets by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserPrizeSetsById(m UserPrizeSets, o orm.Ormer) (err error) {
	v := UserPrizeSets{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserPrizeSets deletes UserPrizeSets by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserPrizeSets(id int) (err error) {
	o := orm.NewOrm()
	v := UserPrizeSets{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserPrizeSets{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
