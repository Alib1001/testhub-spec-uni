package models

import (
	"errors"
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

	// Check if the subject with given ID is already associated with the speciality
	existingSubjects, err := GetSubjectsBySpecialityID(specialityId)
	if err != nil {
		return err
	}
	for _, subj := range existingSubjects {
		if subj.Id == subjectId {
			return errors.New("subject with this ID is already associated with the speciality")
		}
	}

	// Check if the speciality already has two subjects associated
	speciality := &Speciality{Id: specialityId}
	err = o.Read(speciality)
	if err != nil {
		return err
	}

	// Load related subjects to ensure we have the latest data
	o.LoadRelated(speciality, "Subjects")

	if len(speciality.Subjects) >= 2 {
		return errors.New("cannot add more than two subjects to the speciality")
	}

	// Proceed to add the subject to the speciality
	subject := &Subject{Id: subjectId}
	m2m := o.QueryM2M(speciality, "Subjects")
	_, err = m2m.Add(subject)
	if err != nil {
		return err
	}

	// Reload speciality to reflect changes
	err = o.Read(speciality)
	if err != nil {
		return err
	}

	return nil
}

func GetSpecialitiesInUniversity(universityId int) ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	_, err := o.QueryTable("speciality").
		Filter("Universities__University__Id", universityId).
		All(&specialities)
	return specialities, err
}
func GetSubjectsBySpecialityID(specialityId int) ([]*Subject, error) {
	o := orm.NewOrm()

	var subjects []*Subject
	_, err := o.QueryTable("subject").
		Filter("Specialities__Speciality__Id", specialityId).
		All(&subjects)

	return subjects, err
}
