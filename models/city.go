package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type City struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	NameRu       string        `orm:"size(128)"`
	NameKz       string        `orm:"size(128)"`
	Universities []*University `orm:"reverse(many)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

type CityResponse struct {
	Id        int       `json:"Id"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type CityResponseById struct {
	Id     int    `json:"Id"`
	NameRu string `json:"NameRu"`
	NameKz string `json:"NameKz"`
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
	return id, nil
}

func GetCityById(id int, language string) (*City, error) {
	o := orm.NewOrm()
	city := &City{Id: id}
	err := o.Read(city)
	if err != nil {
		return nil, err
	}

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

	switch language {
	case "ru":
		city.Name = city.NameRu
	case "kz":
		city.Name = city.NameKz
	}

	return city, nil
}

func GetAllCitiesByLanguage(language string) ([]*City, error) {
	o := orm.NewOrm()
	var cities []*City
	_, err := o.QueryTable("city").All(&cities)
	if err != nil {
		return nil, err
	}

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
	var field string
	switch language {
	case "ru":
		field = "name_ru"
	case "kz":
		field = "name_kz"
	default:
		return results, fmt.Errorf("unsupported language: %s", language)
	}

	o := orm.NewOrm()
	query := fmt.Sprintf("SELECT * FROM city WHERE %s LIKE ?", field)
	searchPattern := fmt.Sprintf("%s%%", name)

	_, err := o.Raw(query, searchPattern).QueryRows(&results)
	if err != nil {
		return results, err
	}

	return results, nil
}
