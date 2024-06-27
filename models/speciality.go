package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Speciality struct {
	Id           int    `orm:"auto"`
	Name         string `orm:"size(128)"`
	Code         string `orm:"size(64)"`
	AnnualGrants int
	MinScore     int
	MaxScore     int
	VideoLink    string        `orm:"size(256)"`
	Description  string        `orm:"type(text)"`
	Universities []*University `orm:"rel(m2m);rel_table(speciality_university)"`
	Subjects     []*Subject    `orm:"reverse(many)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Speciality))
}

func AddSpeciality(speciality *Speciality) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(speciality)
	return id, err
}

func GetSpecialityById(id int) (*Speciality, error) {
	o := orm.NewOrm()
	speciality := &Speciality{Id: id}
	err := o.Read(speciality)
	return speciality, err
}

func GetAllSpecialities() ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	_, err := o.QueryTable("speciality").All(&specialities)
	return specialities, err
}

func UpdateSpeciality(speciality *Speciality) error {
	o := orm.NewOrm()
	_, err := o.Update(speciality)
	return err
}

func DeleteSpeciality(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Speciality{Id: id})
	return err
}

func AddSubjectToSpeciality(subjectId, specialityId int) error {
	o := orm.NewOrm()

	subject := &Subject{Id: subjectId}
	speciality := &Speciality{Id: specialityId}

	m2m := o.QueryM2M(speciality, "Subjects")
	_, err := m2m.Add(subject)
	return err
}
func GetSpecialitiesInUniversity(universityId int) ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	_, err := o.QueryTable("speciality").
		Filter("Universities__University__Id", universityId).
		All(&specialities)
	return specialities, err
}
