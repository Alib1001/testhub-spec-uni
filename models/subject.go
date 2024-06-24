package models

import (
	"github.com/astaxie/beego/orm"
)

type Subject struct {
	Id          int          `orm:"auto"`
	Name        string       `orm:"size(128)"`
	Specialties []*Specialty `orm:"rel(m2m);rel_table(subject_specialty)"`
}

func init() {
	orm.RegisterModel(new(Subject))
}

// CRUD methods

func AddSubject(subject *Subject) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(subject)
	return id, err
}

func GetSubjectById(id int) (*Subject, error) {
	o := orm.NewOrm()
	subject := &Subject{Id: id}
	err := o.Read(subject)
	return subject, err
}

func GetAllSubjects() ([]*Subject, error) {
	o := orm.NewOrm()
	var subjects []*Subject
	_, err := o.QueryTable("subject").All(&subjects)
	return subjects, err
}

func UpdateSubject(subject *Subject) error {
	o := orm.NewOrm()
	_, err := o.Update(subject)
	return err
}

func DeleteSubject(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Subject{Id: id})
	return err
}
