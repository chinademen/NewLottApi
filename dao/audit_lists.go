package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AuditLists struct {
	Id          int       `orm:"column(id);auto"`
	TypeId      uint      `orm:"column(type_id)" description:"审核类型id"`
	UserId      uint      `orm:"column(user_id)" description:"被审核用户的id"`
	AuthorId    uint      `orm:"column(author_id)" description:"提交人的id"`
	AuditorId   uint      `orm:"column(auditor_id);null" description:"审核人的id"`
	Username    string    `orm:"column(username);size(16)" description:"被审核用户名称"`
	Author      string    `orm:"column(author);size(16)" description:"提交人名称"`
	Auditor     string    `orm:"column(auditor);size(16);null" description:"审核人名称"`
	AuditData   string    `orm:"column(audit_data)" description:"待审核的数据"`
	Description string    `orm:"column(description);size(255);null" description:"备注,描述"`
	Status      int8      `orm:"column(status)" description:"0:审核中, 1:审核通过, 2:审核拒绝"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *AuditLists) TableName() string {
	return "audit_lists"
}

func init() {
	orm.RegisterModel(new(AuditLists))
}

// AddAuditLists insert a new AuditLists into database and returns
// last inserted Id on success.
func AddAuditLists(m *AuditLists) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiAuditLists mutil insert a new AuditLists into database
func AddMultiAuditLists(mlist []*AuditLists) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetAuditListsById retrieves AuditLists by Id. Returns error if
// Id doesn't exist
func GetAuditListsById(id int) (v *AuditLists, err error) {
	o := orm.NewOrm()
	v = &AuditLists{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAuditLists retrieves all AuditLists matches certain condition. Returns empty list if
// no records exist
func GetAllAuditLists(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []AuditLists, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AuditLists))
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

// UpdateAuditLists updates AuditLists by Id and returns error if
// the record to be updated doesn't exist
func UpdateAuditListsById(m *AuditLists) (err error) {
	o := orm.NewOrm()
	v := AuditLists{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAuditLists deletes AuditLists by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAuditLists(id int) (err error) {
	o := orm.NewOrm()
	v := AuditLists{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AuditLists{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
