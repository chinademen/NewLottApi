package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Vocabularies struct {
	Id           int       `orm:"column(id);auto"`
	DictionaryId uint      `orm:"column(dictionary_id)"`
	Dictionary   string    `orm:"column(dictionary);size(64);null"`
	Keyword      string    `orm:"column(keyword);size(100)"`
	En           string    `orm:"column(en);size(512);null"`
	ZhCn         string    `orm:"column(zh_cn);size(512);null"`
	CreatedAt    time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt    time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *Vocabularies) TableName() string {
	return "vocabularies"
}

func init() {
	orm.RegisterModel(new(Vocabularies))
}

// AddVocabularies insert a new Vocabularies into database and returns
// last inserted Id on success.
func AddVocabularies(m *Vocabularies) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiVocabularies mutil insert a new Vocabularies into database
func AddMultiVocabularies(mlist []*Vocabularies) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetVocabulariesById retrieves Vocabularies by Id. Returns error if
// Id doesn't exist
func GetVocabulariesById(id int) (v *Vocabularies, err error) {
	o := orm.NewOrm()
	v = &Vocabularies{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllVocabularies retrieves all Vocabularies matches certain condition. Returns empty list if
// no records exist
func GetAllVocabularies(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Vocabularies, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Vocabularies))
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

// UpdateVocabularies updates Vocabularies by Id and returns error if
// the record to be updated doesn't exist
func UpdateVocabulariesById(m *Vocabularies) (err error) {
	o := orm.NewOrm()
	v := Vocabularies{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteVocabularies deletes Vocabularies by Id and returns error if
// the record to be deleted doesn't exist
func DeleteVocabularies(id int) (err error) {
	o := orm.NewOrm()
	v := Vocabularies{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Vocabularies{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
