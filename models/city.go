package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type City struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	Universities []*University `orm:"reverse(many)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(City))
}

func AddCity(city *City) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(city)
	return id, err
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
