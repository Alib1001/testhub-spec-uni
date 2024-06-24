package models

import (
	"github.com/astaxie/beego/orm"
)

type SpecialtyUniversity struct {
	Id           int `orm:"auto"`
	SpecialtyId  int
	UniversityId int
}

func init() {
	orm.RegisterModel(new(SpecialtyUniversity))
}

// CRUD methods

func AddSpecialtyUniversity(specialtyUniversity *SpecialtyUniversity) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(specialtyUniversity)
	return id, err
}

func GetSpecialtyUniversityById(id int) (*SpecialtyUniversity, error) {
	o := orm.NewOrm()
	specialtyUniversity := &SpecialtyUniversity{Id: id}
	err := o.Read(specialtyUniversity)
	return specialtyUniversity, err
}

func GetAllSpecialtyUniversities() ([]*SpecialtyUniversity, error) {
	o := orm.NewOrm()
	var specialtyUniversities []*SpecialtyUniversity
	_, err := o.QueryTable("specialty_university").All(&specialtyUniversities)
	return specialtyUniversities, err
}

func UpdateSpecialtyUniversity(specialtyUniversity *SpecialtyUniversity) error {
	o := orm.NewOrm()
	_, err := o.Update(specialtyUniversity)
	return err
}

func DeleteSpecialtyUniversity(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&SpecialtyUniversity{Id: id})
	return err
}
