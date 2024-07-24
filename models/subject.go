package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Subject struct {
	Id        int       `orm:"auto"`
	Name      string    `orm:"size(128)"`
	NameKz    string    `orm:"size(128)"`
	NameRu    string    `orm:"size(128)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Subject))
}

func AddSubject(subject *Subject) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(subject)

	if err != nil {
		return 0, err
	}

	subject.Id = int(id)

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

func SearchSubjectsByName(prefix string) ([]Subject, error) {
	var results []Subject

	o := orm.NewOrm()
	searchPattern := fmt.Sprintf("%%%s%%", prefix)

	_, err := o.Raw(`
		SELECT * 
		FROM subject 
		WHERE name LIKE ? 
		OR name_kz LIKE ? 
		OR name_ru LIKE ?
	`, searchPattern, searchPattern, searchPattern).QueryRows(&results)

	if err != nil {
		return results, err
	}

	return results, nil
}

func GetAllowedSecondSubjects(subject1Id int) ([]*Subject, error) {
	o := orm.NewOrm()
	var subjects []*Subject

	_, err := o.Raw(`
		SELECT DISTINCT s2.*
		FROM subject_pair sp
		JOIN subject s2 ON sp.subject2_id = s2.id
		WHERE sp.subject1_id = ?
	`, subject1Id).QueryRows(&subjects)

	if err != nil {
		return nil, err
	}

	return subjects, nil
}
