package models

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testhub-spec-uni/conf"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/elastic/go-elasticsearch/esapi"
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

type SpecialitySearchResponse struct {
	Hits struct {
		Hits []struct {
			Source Speciality `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func init() {
	orm.RegisterModel(new(Speciality))
}

func AddSpeciality(speciality *Speciality) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(speciality)

	if err != nil {
		return 0, err
	}

	speciality.Id = int(id)
	err = IndexSpeciality(speciality)
	if err != nil {
		return id, fmt.Errorf("Speciality added, but failed index in ElasticSearh")
	}

	return id, err
}

func IndexSpeciality(speciality *Speciality) error {
	// Преобразование специальности в JSON
	data, err := json.Marshal(speciality)
	if err != nil {
		return err
	}

	// Создание запроса на индексирование
	req := esapi.IndexRequest{
		Index:      "specialities",
		DocumentID: fmt.Sprintf("%d", speciality.Id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Выполнение запроса
	res, err := req.Do(context.Background(), conf.EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

func SearchSpecialitiesByName(prefix string) ([]Speciality, error) {
	var results []Speciality

	query := fmt.Sprintf(`{
        "query": {
            "query_string": {
                "query": "%s*",
                "fields": ["Name"]
            }
        }
    }`, prefix)

	res, err := conf.EsClient.Search(
		conf.EsClient.Search.WithContext(context.Background()),
		conf.EsClient.Search.WithIndex("specialities"),
		conf.EsClient.Search.WithBody(strings.NewReader(query)),
		conf.EsClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return results, err
		} else {
			return results, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var sr SpecialitySearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
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
