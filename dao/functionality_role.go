package dao

import "time"

type FunctionalityRole struct {
	Id_RENAME       string    `orm:"column(id);size(16)"`
	FunctionalityId string    `orm:"column(functionality_id);size(16)"`
	RoleId          string    `orm:"column(role_id);size(16)"`
	CreatedAt       time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(datetime);null;auto_now_add"`
}
