package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type UniversitySpecialityDetail struct {
	Id         int         `orm:"auto" json:"id"`
	University *University `orm:"rel(fk);on_delete(cascade)" json:"university"`
	Speciality *Speciality `orm:"rel(fk);on_delete(cascade)" json:"speciality"`
	Term       int         `json:"term"`
	EduLang    string      `json:"edu_lang"`
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(UniversitySpecialityDetail))
}
