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

type City struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"` // Для отображения на нужном языке
	NameRu       string        `orm:"size(128)"`
	NameKz       string        `orm:"size(128)"`
	Universities []*University `orm:"reverse(many)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

type CitySearchResponse struct {
	Hits struct {
		Hits []struct {
			Source City `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func init() {
	orm.RegisterModel(new(City))
}

func AddCity(city *City) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(city)
	if err != nil {
		return 0, err
	}
	city.Id = int(id)

	err = IndexCity(city)
	if err != nil {
		return id, fmt.Errorf("city added but failed to index in Elasticsearch: %v", err)
	}

	return id, nil
}

func GetCityById(id int, language string) (*City, error) {
	o := orm.NewOrm()
	city := &City{Id: id}
	err := o.Read(city)
	if err != nil {
		return nil, err
	}

	// Применяем фильтрацию по языку
	switch language {
	case "ru":
		city.Name = city.NameRu
	case "kz":
		city.Name = city.NameKz
	}

	return city, nil
}

func UpdateCity(city *City) error {
	o := orm.NewOrm()
	_, err := o.Update(city)
	return err
}

func DeleteCity(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&City{Id: id})
	return err
}

func GetCityWithUniversities(id int, language string) (*City, error) {
	o := orm.NewOrm()
	city := &City{Id: id}
	if err := o.Read(city); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(city, "Universities"); err != nil {
		return nil, err
	}

	// Применяем фильтрацию по языку
	switch language {
	case "ru":
		city.Name = city.NameRu
	case "kz":
		city.Name = city.NameKz
	}

	return city, nil
}

func IndexCity(city *City) error {
	data, err := json.Marshal(city)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "cities",
		DocumentID: fmt.Sprintf("%d", city.Id),
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

func GetAllCitiesByLanguage(language string) ([]*City, error) {
	o := orm.NewOrm()
	var cities []*City
	_, err := o.QueryTable("city").All(&cities)
	if err != nil {
		return nil, err
	}

	// Применяем фильтрацию по языку
	for _, city := range cities {
		switch language {
		case "ru":
			city.Name = city.NameRu
		case "kz":
			city.Name = city.NameKz
		}
	}

	return cities, nil
}

func SearchCitiesByName(name, language string) ([]City, error) {
	var results []City

	// Применяем фильтрацию по языку
	var field string
	switch language {
	case "ru":
		field = "NameRu"
	case "kz":
		field = "NameKz"
	default:
		return results, fmt.Errorf("unsupported language: %s", language)
	}

	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"%s": "%s*"
			}
		}
	}`, field, name)

	res, err := conf.EsClient.Search(
		conf.EsClient.Search.WithContext(context.Background()),
		conf.EsClient.Search.WithIndex("cities"),
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

	var sr CitySearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
