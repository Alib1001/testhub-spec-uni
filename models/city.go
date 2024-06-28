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
	Name         string        `orm:"size(128)"`
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

func GetCityById(id int) (*City, error) {
	o := orm.NewOrm()
	city := &City{Id: id}
	err := o.Read(city)
	return city, err
}

func GetAllCities() ([]*City, error) {
	o := orm.NewOrm()
	var cities []*City
	_, err := o.QueryTable("city").All(&cities)
	return cities, err
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

func GetCityWithUniversities(id int) (*City, error) {
	o := orm.NewOrm()
	city := &City{Id: id}
	if err := o.Read(city); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(city, "Universities"); err != nil {
		return nil, err
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

func SearchCitiesByName(prefix string) ([]City, error) {
	var results []City

	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"Name": "%s*"
			}
		}
	}`, prefix)

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
