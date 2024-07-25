package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Subject struct {
	Id        int       `orm:"auto"`
	Name      string    `orm:"size(128)"`
	NameRu    string    `orm:"size(128)"`
	NameKz    string    `orm:"size(128)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Subject))
}

type SubjectResponse struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

func AddSubject(subject *Subject) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(subject)
	if err != nil {
		return 0, err
	}
	subject.Id = int(id)
	return id, nil
}

func GetSubjectById(id int, language string) (*SubjectResponse, error) {
	o := orm.NewOrm()
	subject := &Subject{Id: id}
	err := o.Read(subject)
	if err != nil {
		return nil, err
	}

	name := subject.Name
	switch language {
	case "ru":
		name = subject.NameRu
	case "kz":
		name = subject.NameKz
	}

	return &SubjectResponse{
		Id:   subject.Id,
		Name: name,
	}, nil
}

func GetAllSubjects(language string) ([]*SubjectResponse, error) {
	o := orm.NewOrm()
	var subjects []*Subject
	_, err := o.QueryTable("subject").All(&subjects)
	if err != nil {
		return nil, err
	}

	var subjectResponses []*SubjectResponse
	for _, subject := range subjects {
		name := subject.Name
		switch language {
		case "ru":
			name = subject.NameRu
		case "kz":
			name = subject.NameKz
		}
		subjectResponses = append(subjectResponses, &SubjectResponse{
			Id:   subject.Id,
			Name: name,
		})
	}

	return subjectResponses, nil
}

func UpdateSubject(subject *Subject) error {
	o := orm.NewOrm()
	existingSubject := Subject{Id: subject.Id}
	if err := o.Read(&existingSubject); err != nil {
		return err
	}

	if subject.Name != "" {
		existingSubject.Name = subject.Name
	}
	if subject.NameRu != "" {
		existingSubject.NameRu = subject.NameRu
	}
	if subject.NameKz != "" {
		existingSubject.NameKz = subject.NameKz
	}

	_, err := o.Update(&existingSubject)
	return err
}

func DeleteSubject(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Subject{Id: id})
	return err
}

func SearchSubjectsByName(prefix, language string) ([]SubjectResponse, error) {
	var results []Subject
	var field string

	switch language {
	case "ru":
		field = "name_ru"
	case "kz":
		field = "name_kz"
	default:
		field = "name"
	}

	o := orm.NewOrm()
	query := fmt.Sprintf("SELECT * FROM subject WHERE %s LIKE ?", field)
	searchPattern := fmt.Sprintf("%%%s%%", prefix)

	_, err := o.Raw(query, searchPattern).QueryRows(&results)
	if err != nil {
		return nil, err
	}

	var subjectResponses []SubjectResponse
	for _, subject := range results {
		name := subject.Name
		switch language {
		case "ru":
			name = subject.NameRu
		case "kz":
			name = subject.NameKz
		}
		subjectResponses = append(subjectResponses, SubjectResponse{
			Id:   subject.Id,
			Name: name,
		})
	}

	return subjectResponses, nil
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
