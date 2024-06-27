package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Subject struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	Specialities []*Speciality `orm:"rel(m2m);rel_table(subject_speciality)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Subject))
}

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
