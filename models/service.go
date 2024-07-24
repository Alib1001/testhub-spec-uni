package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
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

func SearchServicesByName(prefix string) ([]Service, error) {
	var results []Service

	o := orm.NewOrm()
	query := "SELECT * FROM service WHERE name LIKE ?"
	searchPattern := fmt.Sprintf("%s%%", prefix)

	_, err := o.Raw(query, searchPattern).QueryRows(&results)
	if err != nil {
		return results, err
	}

	return results, nil
}
