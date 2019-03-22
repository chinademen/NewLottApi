package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Users struct {
	Id             int       `orm:"column(id);auto" json:"id,string"`
	Username       string    `orm:"column(username);size(16)" json:"username,string" description:"用户名"`
	Password       string    `orm:"column(password);size(60)" json:"password,string" description:"密码"`
	FundPassword   string    `orm:"column(fund_password);size(60)" json:"fund_password,string" description:"资金密码"`
	AccountId      uint      `orm:"column(account_id);null" json:"account_id,string" description:"账户id"`
	MerchantId     int       `orm:"column(merchant_id)" json:"merchant_id,string" description:"商户ID"`
	PrizeGroup     string    `orm:"column(prize_group);size(20);null" json:"prize_group,string" description:"奖金组"`
	Blocked        uint8     `orm:"column(blocked)" json:"blocked,string" description:"1=冻结"`
	Realname       string    `orm:"column(realname);size(30);null" json:"realname,string" description:"真实姓名"`
	Nickname       string    `orm:"column(nickname);size(16)" json:"nickname,string" description:"昵称"`
	Email          string    `orm:"column(email);size(50);null" json:"email,string" description:"邮件"`
	Mobile         string    `orm:"column(mobile);size(20);null" json:"mobile,string" description:"电话号码"`
	IsTester       int8      `orm:"column(is_tester);null" json:"is_tester,string" description:"1=测试"`
	BetMultiple    uint      `orm:"column(bet_multiple);null" json:"bet_multiple,string" description:"投注倍数"`
	BetCoefficient float64   `orm:"column(bet_coefficient);null;digits(4);decimals(3)" json:"bet_coefficient,string" description:"投注模式"`
	LoginIp        string    `orm:"column(login_ip);size(15);null" json:"login_ip,string" description:"最后登录ip"`
	RegisterIp     string    `orm:"column(register_ip);size(15);null" json:"register_ip,string" description:"注册ip"`
	Token          string    `orm:"column(token);size(200);null" json:"token,string" description:"忘记密码的token"`
	SigninAt       time.Time `orm:"column(signin_at);type(datetime);null" json:"signin_at,string" description:"登录时间"`
	ActivatedAt    time.Time `orm:"column(activated_at);type(datetime);null" json:"id,string" description:"活跃时间，投注"`
	RegisterAt     time.Time `orm:"column(register_at);type(datetime);null" json:"id,string" description:"注册时间"`
	DeletedAt      time.Time `orm:"column(deleted_at);type(datetime);null" json:"id,string" description:"删除时间"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime);null;auto_now_add" description:"数据库创建时间"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime);null" description:"更新时间"`
}

func (t *Users) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(Users))
}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiUsers mutil insert a new Users into database
func AddMultiUsers(mlist []*Users) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetUsersById retrieves Users by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []*Users, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
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

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users, o orm.Ormer) (err error) {
	v := Users{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsers deletes Users by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int) (err error) {
	o := orm.NewOrm()
	v := Users{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Users{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
