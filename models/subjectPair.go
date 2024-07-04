package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type SubjectPair struct {
	Id           int           `orm:"auto"`
	Subject1     *Subject      `orm:"rel(fk);on_delete(do_nothing)"`
	Subject2     *Subject      `orm:"rel(fk);on_delete(do_nothing)"`
	Specialities []*Speciality `orm:"reverse(many)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(SubjectPair))
}

func AddSubjectPair(subjectPair *SubjectPair) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(subjectPair)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetSubjectPairById(id int) (*SubjectPair, error) {
	o := orm.NewOrm()
	subjectPair := &SubjectPair{Id: id}
	err := o.Read(subjectPair)
	if err != nil {
		return nil, err
	}

	// Load related subjects
	if subjectPair.Subject1 != nil {
		err = o.Read(subjectPair.Subject1)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}
	}
	if subjectPair.Subject2 != nil {
		err = o.Read(subjectPair.Subject2)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}
	}

	return subjectPair, nil
}

func GetAllSubjectPairs() ([]*SubjectPair, error) {
	o := orm.NewOrm()
	var subjectPairs []*SubjectPair
	_, err := o.QueryTable("subject_pair").All(&subjectPairs)
	if err != nil {
		return nil, err
	}

	// Load related subjects for each pair
	for _, pair := range subjectPairs {
		if pair.Subject1 != nil {
			err := o.Read(pair.Subject1)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
		}
		if pair.Subject2 != nil {
			err := o.Read(pair.Subject2)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
		}
	}

	return subjectPairs, nil
}

func UpdateSubjectPair(subjectPair *SubjectPair) error {
	o := orm.NewOrm()
	_, err := o.Update(subjectPair)
	return err
}

func DeleteSubjectPair(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&SubjectPair{Id: id})
	return err
}

func GetSubjectPairBySubjectIds(subject1Id, subject2Id int) (*SubjectPair, error) {
	o := orm.NewOrm()
	subjectPair := &SubjectPair{}
	err := o.QueryTable("subject_pair").Filter("Subject1__Id", subject1Id).Filter("Subject2__Id", subject2Id).One(subjectPair)
	if err != nil {
		return nil, err
	}
	return subjectPair, nil
}
