package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Functionalities struct {
	Id             int       `orm:"column(id);auto" description:"菜单ID"`
	Title          string    `orm:"column(title);size(64);null"`
	ParentId       uint      `orm:"column(parent_id);null"`
	Parent         string    `orm:"column(parent);size(50)"`
	ForefatherIds  string    `orm:"column(forefather_ids);size(100)"`
	Forefathers    string    `orm:"column(forefathers);size(10240);null"`
	Description    string    `orm:"column(description);size(255)"`
	Controller     string    `orm:"column(controller);size(40)"`
	Action         string    `orm:"column(action);size(40)"`
	ButtonType     int8      `orm:"column(button_type)"`
	PopupId        uint      `orm:"column(popup_id);null"`
	PopupTitle     string    `orm:"column(popup_title);size(64);null"`
	ButtonOnclick  string    `orm:"column(button_onclick);size(64);null"`
	ConfirmMsgKey  string    `orm:"column(confirm_msg_key);size(200);null"`
	RefreshCycle   uint16    `orm:"column(refresh_cycle);null"`
	Menu           uint8     `orm:"column(menu)" description:"是否菜单项"`
	NeedCurd       int8      `orm:"column(need_curd)" description:"是否需要CURD权限"`
	NeedSearch     int8      `orm:"column(need_search)" description:"是否需要搜索表单"`
	SearchConfigId uint      `orm:"column(search_config_id);null"`
	Realm          uint8     `orm:"column(realm);null" description:"领域:1为管理，2为用户，3为全部"`
	NeedLog        uint8     `orm:"column(need_log)"`
	Disabled       uint8     `orm:"column(disabled)" description:"菜单是否启用（0 正常 1关闭）"`
	Sequence       uint      `orm:"column(sequence)" description:"菜单排序"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *Functionalities) TableName() string {
	return "functionalities"
}

func init() {
	orm.RegisterModel(new(Functionalities))
}

// AddFunctionalities insert a new Functionalities into database and returns
// last inserted Id on success.
func AddFunctionalities(m *Functionalities) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiFunctionalities mutil insert a new Functionalities into database
func AddMultiFunctionalities(mlist []*Functionalities) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetFunctionalitiesById retrieves Functionalities by Id. Returns error if
// Id doesn't exist
func GetFunctionalitiesById(id int) (v *Functionalities, err error) {
	o := orm.NewOrm()
	v = &Functionalities{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFunctionalities retrieves all Functionalities matches certain condition. Returns empty list if
// no records exist
func GetAllFunctionalities(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Functionalities, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Functionalities))
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

// UpdateFunctionalities updates Functionalities by Id and returns error if
// the record to be updated doesn't exist
func UpdateFunctionalitiesById(m *Functionalities) (err error) {
	o := orm.NewOrm()
	v := Functionalities{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFunctionalities deletes Functionalities by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFunctionalities(id int) (err error) {
	o := orm.NewOrm()
	v := Functionalities{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Functionalities{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
