package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testhub-spec-uni/conf"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/elastic/go-elasticsearch/esapi"
)

type Subject struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	Specialities []*Speciality `orm:"rel(m2m);rel_table(subject_speciality)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

type SubjectSearchResponse struct {
	Hits struct {
		Hits []struct {
			Source Subject `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
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
	err = IndexSubject(subject)
	if err != nil {
		return id, fmt.Errorf("Subject added, but failed index in ElasticSearh")
	}

	return id, err
}

func IndexSubject(subject *Subject) error {
	data, err := json.Marshal(subject)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "subjects",
		DocumentID: fmt.Sprintf("%d", subject.Id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

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
		conf.EsClient.Search.WithIndex("subjects"),
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

	var sr SubjectSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
