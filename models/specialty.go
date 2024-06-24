package models

import (
	"github.com/astaxie/beego/orm"
)

type Specialty struct {
	Id           int    `orm:"auto"`
	Name         string `orm:"size(128)"`
	Code         string `orm:"size(64)"`
	AnnualGrants int
	MinScore     int
	MaxScore     int
	VideoLink    string `orm:"size(256)"`
	Description  string `orm:"type(text)"`
}

func init() {
	orm.RegisterModel(new(Specialty))
}

// CRUD methods

func AddSpecialty(specialty *Specialty) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(specialty)
	return id, err
}

func GetSpecialtyById(id int) (*Specialty, error) {
	o := orm.NewOrm()
	specialty := &Specialty{Id: id}
	err := o.Read(specialty)
	return specialty, err
}

func GetAllSpecialties() ([]*Specialty, error) {
	o := orm.NewOrm()
	var specialties []*Specialty
	_, err := o.QueryTable("specialty").All(&specialties)
	return specialties, err
}

func UpdateSpecialty(specialty *Specialty) error {
	o := orm.NewOrm()
	_, err := o.Update(specialty)
	return err
}

func DeleteSpecialty(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Specialty{Id: id})
	return err
}
