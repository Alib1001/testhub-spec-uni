package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type SpecialityUniversity struct {
	Id         int         `orm:"auto" json:"id"`
	University *University `orm:"rel(fk);on_delete(cascade)" json:"university"`
	Speciality *Speciality `orm:"rel(fk);on_delete(cascade)" json:"speciality"`
	Term       int         `json:"term"`
	EduLang    string      `json:"edu_lang"`
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(SpecialityUniversity))
}

func (u *SpecialityUniversity) Create() error {
	o := orm.NewOrm()
	_, err := o.Insert(u)
	return err
}

func (u *SpecialityUniversity) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(u)
	return err
}

func GetByUniversityAndSpeciality(universityID, specialityID int) (*SpecialityUniversity, error) {
	o := orm.NewOrm()
	detail := &SpecialityUniversity{
		University: &University{Id: universityID},
		Speciality: &Speciality{Id: specialityID},
	}
	err := o.QueryTable("speciality_university").Filter("university_id", universityID).Filter("speciality_id", specialityID).One(detail)
	return detail, err
}

func DeleteByUniversityAndSpeciality(universityID, specialityID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("speciality_university").Filter("university_id", universityID).Filter("speciality_id", specialityID).Delete()
	return err
}
