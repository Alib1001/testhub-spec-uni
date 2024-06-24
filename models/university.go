package models

import (
	"github.com/astaxie/beego/orm"
)

type University struct {
	Id               int    `orm:"auto"`
	UniversityCode   string `orm:"size(64)"`
	Name             string `orm:"size(128)"`
	Abbreviation     string `orm:"size(64)"`
	UniversityStatus string `orm:"size(64)"`
	CityId           int
	Address          string   `orm:"size(256)"`
	Website          string   `orm:"size(128)"`
	SocialMediaList  []string `orm:"-"`
	ContactList      []string `orm:"-"`
	AverageFee       int
	HasMilitaryDept  bool
	HasDormitory     bool
	ProfileImageUrl  string `orm:"size(256)"`
	MinEntryScore    int
	PhotosUrlList    []string     `orm:"-"`
	Description      string       `orm:"type(text)"`
	Specialties      []*Specialty `orm:"rel(m2m);rel_table(specialty_university)"`
}

func init() {
	orm.RegisterModel(new(University))
}

// CRUD methods

func AddUniversity(university *University) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(university)
	return id, err
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
		Filter("city_id", cityId).
		All(&universities)
	return universities, err
}
