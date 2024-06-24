package models

import (
	"github.com/astaxie/beego/orm"
)

type SpecialtySubject struct {
	Id         int `orm:"auto"`
	Subject1Id int
	Subject2Id int
	Specialty  *Specialty `orm:"rel(fk)"`
}

func init() {
	orm.RegisterModel(new(SpecialtySubject))
}

// CRUD methods

func AddSpecialtySubject(specialtySubject *SpecialtySubject) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(specialtySubject)
	return id, err
}

func GetSpecialtySubjectById(id int) (*SpecialtySubject, error) {
	o := orm.NewOrm()
	specialtySubject := &SpecialtySubject{Id: id}
	err := o.Read(specialtySubject)
	return specialtySubject, err
}

func GetAllSpecialtySubjects() ([]*SpecialtySubject, error) {
	o := orm.NewOrm()
	var specialtySubjects []*SpecialtySubject
	_, err := o.QueryTable("specialty_subject").All(&specialtySubjects)
	return specialtySubjects, err
}

func UpdateSpecialtySubject(specialtySubject *SpecialtySubject) error {
	o := orm.NewOrm()
	_, err := o.Update(specialtySubject)
	return err
}

func DeleteSpecialtySubject(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&SpecialtySubject{Id: id})
	return err
}
