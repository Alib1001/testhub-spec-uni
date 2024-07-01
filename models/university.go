package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testhub-spec-uni/conf"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/elastic/go-elasticsearch/esapi"
)

type University struct {
	Id               int      `orm:"auto"`
	UniversityCode   string   `orm:"size(64)"`
	Name             string   `orm:"size(128)"`
	Abbreviation     string   `orm:"size(64)"`
	UniversityStatus string   `orm:"size(64)"`
	Address          string   `orm:"size(256)"`
	Website          string   `orm:"size(128)"`
	SocialMediaList  []string `orm:"-"`
	ContactList      []string `orm:"-"`
	AverageFee       int
	HasMilitaryDept  bool
	HasDormitory     bool
	ProfileImageUrl  string `orm:"size(256)"`
	MinEntryScore    int
	PhotosUrlList    []string      `orm:"-"`
	Description      string        `orm:"type(text)"`
	Specialities     []*Speciality `orm:"rel(m2m);rel_table(speciality_university)"`
	City             *City         `orm:"rel(fk)"`
	CreatedAt        time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt        time.Time     `orm:"auto_now;type(datetime)"`
}

type UniversitySearchResponse struct {
	Hits struct {
		Hits []struct {
			Source University `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func init() {
	orm.RegisterModel(new(University))
}

func AddUniversity(university *University) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(university)
	if err != nil {
		return 0, err
	}
	university.Id = int(id)

	err = IndexUniversity(university)
	if err != nil {
		return id, fmt.Errorf("university added but failed to index in Elasticsearch: %v", err)
	}

	return id, nil
}

func GetUniversityById(id int) (*University, error) {
	o := orm.NewOrm()
	university := &University{Id: id}
	err := o.Read(university)
	return university, err
}

func GetAllUniversities() ([]*University, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").All(&universities)
	return universities, err
}

func UpdateUniversity(university *University) error {
	o := orm.NewOrm()
	_, err := o.Update(university)
	return err
}

func DeleteUniversity(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&University{Id: id})
	return err
}

func GetUniversitiesInCity(cityId int) ([]*University, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").
		Filter("City__Id", cityId).
		All(&universities)
	return universities, err
}

func AssignCityToUniversity(universityId int, cityId int) error {
	o := orm.NewOrm()

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}
	city := &City{Id: cityId}
	if err := o.Read(city); err != nil {
		return err
	}

	university.City = city

	if _, err := o.Update(university); err != nil {
		return err
	}

	return nil
}
func AddSpecialityToUniversity(specialityId, universityId int) error {
	o := orm.NewOrm()

	speciality := &Speciality{Id: specialityId}
	if err := o.Read(speciality); err != nil {
		return err
	}

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}

	exist := o.QueryM2M(university, "Specialities").Exist(speciality)
	if exist {
		return fmt.Errorf("speciality with ID %d is already assigned to university with ID %d", specialityId, universityId)
	}

	_, err := o.QueryM2M(university, "Specialities").Add(speciality)
	if err != nil {
		return err
	}

	o.LoadRelated(university, "Specialities")
	fmt.Printf("Specialities for university %d: %v\n", universityId, university.Specialities)

	return nil
}

func IndexUniversity(university *University) error {
	data, err := json.Marshal(university)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "universities",
		DocumentID: fmt.Sprintf("%d", university.Id),
		Body:       strings.NewReader(string(data)),
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

func SearchUniversitiesByName(prefix string) ([]University, error) {
	var results []University

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
		conf.EsClient.Search.WithIndex("universities"),
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

	var sr UniversitySearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
func SearchUniversities(params map[string]interface{}) ([]University, error) {
	var results []University

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
	}

	if minScore, ok := params["min_score"].(int); ok {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"range": map[string]interface{}{
					"MinEntryScore": map[string]interface{}{
						"gte": minScore,
					},
				},
			},
		)
	}

	if hasMilitaryDept, ok := params["has_military_dept"].(bool); ok {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					"HasMilitaryDept": hasMilitaryDept,
				},
			},
		)
	}

	if hasDormitory, ok := params["has_dormitory"].(bool); ok {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					"HasDormitory": hasDormitory,
				},
			},
		)
	}

	if cityID, ok := params["city_id"].(int); ok {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					"City.Id": cityID,
				},
			},
		)
	}

	if speciality, ok := params["speciality"].(string); ok {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					"Specialities.Name": speciality,
				},
			},
		)
	}

	// Добавьте другие фильтры (если есть) подобным образом

	queryBody, err := json.Marshal(query)
	if err != nil {
		return results, err
	}

	res, err := conf.EsClient.Search(
		conf.EsClient.Search.WithContext(context.Background()),
		conf.EsClient.Search.WithIndex("universities"),
		conf.EsClient.Search.WithBody(strings.NewReader(string(queryBody))),
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

	var sr UniversitySearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
