package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type MethodTypes struct {
	Id            int       `orm:"column(id);auto"`
	LotteryType   uint8     `orm:"column(lottery_type)" description:"彩票类型: 1-数字排列类型 2-乐透类型"`
	Name          string    `orm:"column(name);size(20)"`
	AttributeCode string    `orm:"column(attribute_code);size(20);null" description:"特征码：A：区间；S：大小单双；C：不定位，I：趣味；O：原始"`
	WnFunction    string    `orm:"column(wn_function);size(64)" description:"计奖方法"`
	Sequencing    int8      `orm:"column(sequencing)" description:"定位"`
	DigitalCount  int8      `orm:"column(digital_count);null" description:"星"`
	UniqueCount   int8      `orm:"column(unique_count);null" description:"去重后的数字个数"`
	MaxRepeatTime int8      `orm:"column(max_repeat_time);null" description:"重号的最大重复次"`
	MinRepeatTime int8      `orm:"column(min_repeat_time);null" description:"重号的最小重复次数"`
	Shaped        int8      `orm:"column(shaped)" description:"是否考察数字的属性"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}

func (t *MethodTypes) TableName() string {
	return "method_types"
}

func init() {
	orm.RegisterModel(new(MethodTypes))
}

// AddMethodTypes insert a new MethodTypes into database and returns
// last inserted Id on success.
func AddMethodTypes(m *MethodTypes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiMethodTypes mutil insert a new MethodTypes into database
func AddMultiMethodTypes(mlist []*MethodTypes) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetMethodTypesById retrieves MethodTypes by Id. Returns error if
// Id doesn't exist
func GetMethodTypesById(id int) (v *MethodTypes, err error) {
	o := orm.NewOrm()
	v = &MethodTypes{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMethodTypes retrieves all MethodTypes matches certain condition. Returns empty list if
// no records exist
func GetAllMethodTypes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []MethodTypes, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MethodTypes))
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

// UpdateMethodTypes updates MethodTypes by Id and returns error if
// the record to be updated doesn't exist
func UpdateMethodTypesById(m *MethodTypes) (err error) {
	o := orm.NewOrm()
	v := MethodTypes{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMethodTypes deletes MethodTypes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMethodTypes(id int) (err error) {
	o := orm.NewOrm()
	v := MethodTypes{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MethodTypes{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
