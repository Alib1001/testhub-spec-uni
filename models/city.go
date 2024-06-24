package models

import (
	"github.com/astaxie/beego/orm"
)

type City struct {
	Id   int    `orm:"auto"`
	Name string `orm:"size(128)"`
}

func init() {
	orm.RegisterModel(new(City))
}

// CRUD methods

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
