package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Service struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	NameRu       string        `orm:"size(128)"`
	NameKz       string        `orm:"size(128)"`
	ImageUrl     string        `orm:"size(256)"`
	Universities []*University `orm:"reverse(many)"`
}

type AddServiceForAdminResponse struct {
	Id       int    `orm:"auto"`
	NameRu   string `form:"NameRu" validate:"required"`
	NameKz   string `form:"NameKz" validate:"required"`
	ImageUrl string `form:"ImageUrl"`
}

type ServiceResponseForUser struct {
	Id       int    `json:"Id"`
	Name     string `json:"Name"`
	ImageUrl string `json:"ImageUrl"`
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
	return id, nil
}

func DeleteService(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Service{Id: id})
	return err
}

func UpdateService(service *Service, fields ...string) error {
	o := orm.NewOrm()
	_, err := o.Update(service, fields...)
	if err != nil {
		return err
	}
	return nil
}

func GetServiceById(id int, language string) (*ServiceResponseForUser, error) {
	o := orm.NewOrm()
	service := &Service{Id: id}
	err := o.Read(service)
	if err != nil {
		return nil, err
	}

	name := service.Name
	switch language {
	case "ru":
		name = service.NameRu
	case "kz":
		name = service.NameKz
	}

	return &ServiceResponseForUser{
		Id:       service.Id,
		Name:     name,
		ImageUrl: service.ImageUrl,
	}, nil
}
func GetServiceByID(serviceID int) (*Service, error) {
	o := orm.NewOrm()
	service := &Service{Id: serviceID}
	err := o.Read(service)
	if err == orm.ErrNoRows {
		return nil, errors.New("Service not found")
	} else if err != nil {
		return nil, err
	}
	return service, nil
}

func GetAllServices(language string) ([]*ServiceResponseForUser, error) {
	o := orm.NewOrm()
	var services []*Service
	_, err := o.QueryTable("service").All(&services)
	if err != nil {
		return nil, err
	}

	var serviceResponses []*ServiceResponseForUser
	for _, service := range services {
		name := service.Name
		switch language {
		case "ru":
			name = service.NameRu
		case "kz":
			name = service.NameKz
		}
		serviceResponses = append(serviceResponses, &ServiceResponseForUser{
			Id:       service.Id,
			Name:     name,
			ImageUrl: service.ImageUrl,
		})
	}

	return serviceResponses, nil
}

func SearchServicesByName(prefix, language string) ([]ServiceResponseForUser, error) {
	var results []Service
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
	query := fmt.Sprintf("SELECT * FROM service WHERE %s LIKE ?", field)
	searchPattern := fmt.Sprintf("%s%%", prefix)

	_, err := o.Raw(query, searchPattern).QueryRows(&results)
	if err != nil {
		return nil, err
	}

	var serviceResponses []ServiceResponseForUser
	for _, service := range results {
		name := service.Name
		switch language {
		case "ru":
			name = service.NameRu
		case "kz":
			name = service.NameKz
		}
		serviceResponses = append(serviceResponses, ServiceResponseForUser{
			Id:       service.Id,
			Name:     name,
			ImageUrl: service.ImageUrl,
		})
	}

	return serviceResponses, nil
}

func GetServicesByUniversityId(universityId int, language string) ([]*ServiceResponseForUser, error) {
	o := orm.NewOrm()
	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return nil, err
	}

	var services []*Service
	_, err := o.QueryTable("service").Filter("Universities__University__Id", universityId).All(&services)
	if err != nil {
		return nil, err
	}

	var serviceResponses []*ServiceResponseForUser
	for _, service := range services {
		name := service.Name
		switch language {
		case "ru":
			name = service.NameRu
		case "kz":
			name = service.NameKz
		}
		serviceResponses = append(serviceResponses, &ServiceResponseForUser{
			Id:       service.Id,
			Name:     name,
			ImageUrl: service.ImageUrl,
		})
	}

	return serviceResponses, nil
}

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
