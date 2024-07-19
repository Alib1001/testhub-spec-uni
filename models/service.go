package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testhub-spec-uni/conf"

	"github.com/astaxie/beego/orm"
	"github.com/elastic/go-elasticsearch/esapi"
)

type Service struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	ImageUrl     string        `orm:"size(256)"`
	Universities []*University `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(Service))
}

func AddService(service *Service) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(service)
	if err != nil {
		return 0, err
	}
	service.Id = int(id)

	err = IndexService(service)
	if err != nil {
		return id, fmt.Errorf("service added but failed to index in Elasticsearch: %v", err)
	}

	return id, nil
}

func DeleteService(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Service{Id: id})
	return err
}

// UpdateService updates the details of an existing service in the database
func UpdateService(service *Service) error {
	o := orm.NewOrm()
	_, err := o.Update(service)
	if err != nil {
		return err
	}

	err = IndexService(service)
	if err != nil {
		return fmt.Errorf("service updated but failed to index in Elasticsearch: %v", err)
	}

	return nil
}

// GetServiceById retrieves a service by its ID from the database
func GetServiceById(id int) (*Service, error) {
	o := orm.NewOrm()
	service := &Service{Id: id}
	err := o.Read(service)
	return service, err
}

// GetAllServices retrieves all services from the database
func GetAllServices() ([]*Service, error) {
	o := orm.NewOrm()
	var services []*Service
	_, err := o.QueryTable("service").All(&services)
	return services, err
}

// GetServicesByUniversityId retrieves services associated with a university by its ID
func GetServicesByUniversityId(universityId int) ([]*Service, error) {
	o := orm.NewOrm()

	// Create a university object to read by its ID
	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return nil, err
	}

	// Load related services for the university
	var services []*Service
	_, err := o.QueryTable("service").Filter("Universities__University__Id", universityId).All(&services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

// AddServiceToUniversity binds a service to a university by their IDs
func AddServiceToUniversity(serviceId, universityId int) error {
	o := orm.NewOrm()

	service := &Service{Id: serviceId}
	if err := o.Read(service); err != nil {
		return err
	}

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}

	exist := o.QueryM2M(university, "Services").Exist(service)
	if exist {
		return fmt.Errorf("service with ID %d is already assigned to university with ID %d", serviceId, universityId)
	}

	_, err := o.QueryM2M(university, "Services").Add(service)
	if err != nil {
		return err
	}

	o.LoadRelated(university, "Services")
	fmt.Printf("Services for university %d: %v\n", universityId, university.Services)

	return nil
}

// IndexService indexes a service in Elasticsearch
func IndexService(service *Service) error {
	// Convert service to JSON
	data, err := json.Marshal(service)
	if err != nil {
		return err
	}

	// Create Elasticsearch index request
	req := esapi.IndexRequest{
		Index:      "services",
		DocumentID: fmt.Sprintf("%d", service.Id),
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	// Execute the request
	res, err := req.Do(context.Background(), conf.EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Handle response
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		}
		return fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	return nil
}

// SearchServicesByName searches for services by name in Elasticsearch
func SearchServicesByName(prefix string) ([]Service, error) {
	var results []Service

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
		conf.EsClient.Search.WithIndex("services"),
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

	var sr struct {
		Hits struct {
			Hits []struct {
				Source Service `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
