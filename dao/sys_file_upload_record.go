package dao

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SysFileUploadRecord struct {
	Id               int       `orm:"column(id);auto" description:"上传文件ID"`
	Size             int       `orm:"column(size)" description:"文件大小"`
	Ext              string    `orm:"column(ext);size(10)" description:"文件扩展名"`
	OriginalFileName string    `orm:"column(original_file_name);size(255)" description:"源文件名"`
	NewFileName      string    `orm:"column(new_file_name);size(255)" description:"上传后文件名"`
	CreatedUser      string    `orm:"column(created_user);size(128)" description:"上传操作人"`
	CreatedAt        time.Time `orm:"column(created_at);type(datetime)"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(datetime);auto_now"`
}

func (t *SysFileUploadRecord) TableName() string {
	return "sys_file_upload_record"
}

func init() {
	orm.RegisterModel(new(SysFileUploadRecord))
}

// AddSysFileUploadRecord insert a new SysFileUploadRecord into database and returns
// last inserted Id on success.
func AddSysFileUploadRecord(m *SysFileUploadRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddMultiSysFileUploadRecord mutil insert a new SysFileUploadRecord into database
func AddMultiSysFileUploadRecord(mlist []*SysFileUploadRecord) (err error) {
	o := orm.NewOrm()
	_, err = o.InsertMulti(20, mlist)
	return
}

// GetSysFileUploadRecordById retrieves SysFileUploadRecord by Id. Returns error if
// Id doesn't exist
func GetSysFileUploadRecordById(id int) (v *SysFileUploadRecord, err error) {
	o := orm.NewOrm()
	v = &SysFileUploadRecord{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSysFileUploadRecord retrieves all SysFileUploadRecord matches certain condition. Returns empty list if
// no records exist
func GetAllSysFileUploadRecord(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []SysFileUploadRecord, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SysFileUploadRecord))
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

// UpdateSysFileUploadRecord updates SysFileUploadRecord by Id and returns error if
// the record to be updated doesn't exist
func UpdateSysFileUploadRecordById(m *SysFileUploadRecord) (err error) {
	o := orm.NewOrm()
	v := SysFileUploadRecord{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSysFileUploadRecord deletes SysFileUploadRecord by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSysFileUploadRecord(id int) (err error) {
	o := orm.NewOrm()
	v := SysFileUploadRecord{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SysFileUploadRecord{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
