package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type FunctionalityRelations struct {
	Id               int       `orm:"column(id);auto"`
	FunctionalityId  uint      `orm:"column(functionality_id)"`
	RFunctionalityId uint      `orm:"column(r_functionality_id)"`
	Realm            uint8     `orm:"column(realm);null"`
	Position         uint8     `orm:"column(position);null"`
	ButtonOnclick    string    `orm:"column(button_onclick);size(64);null"`
	ConfirmMsgKey    string    `orm:"column(confirm_msg_key);size(64);null"`
	ForPage          int8      `orm:"column(for_page)"`
	ForPageBatch     int8      `orm:"column(for_page_batch)"`
	ForItem          int8      `orm:"column(for_item)"`
	Label            string    `orm:"column(label);size(50);null"`
	Precondition     string    `orm:"column(precondition);size(200)"`
	Params           string    `orm:"column(params);size(200)"`
	NewWindow        int8      `orm:"column(new_window)"`
	UseRedirector    int8      `orm:"column(use_redirector);null"`
	Disabled         int8      `orm:"column(disabled)"`
	Sequence         uint      `orm:"column(sequence)" description:"排序"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *FunctionalityRelations) TableName() string {
	return "functionality_relations"
}

func init() {
	orm.RegisterModel(new(FunctionalityRelations))
}

// AddFunctionalityRelations insert a new FunctionalityRelations into database and returns
// last inserted Id on success.
func AddFunctionalityRelations(m *FunctionalityRelations) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiFunctionalityRelations mutil insert a new FunctionalityRelations into database
func AddMultiFunctionalityRelations(mlist []*FunctionalityRelations) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetFunctionalityRelationsById retrieves FunctionalityRelations by Id. Returns error if
// Id doesn't exist
func GetFunctionalityRelationsById(id int) (v *FunctionalityRelations, err error) {
	o := orm.NewOrm()
	v = &FunctionalityRelations{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFunctionalityRelations retrieves all FunctionalityRelations matches certain condition. Returns empty list if
// no records exist
func GetAllFunctionalityRelations(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []FunctionalityRelations, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(FunctionalityRelations))
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

// UpdateFunctionalityRelations updates FunctionalityRelations by Id and returns error if
// the record to be updated doesn't exist
func UpdateFunctionalityRelationsById(m *FunctionalityRelations) (err error) {
	o := orm.NewOrm()
	v := FunctionalityRelations{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFunctionalityRelations deletes FunctionalityRelations by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFunctionalityRelations(id int) (err error) {
	o := orm.NewOrm()
	v := FunctionalityRelations{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&FunctionalityRelations{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
