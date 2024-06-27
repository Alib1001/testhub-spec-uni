package models

import (
	"time"

	"github.com/astaxie/beego/orm"
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
	Specialities     []*Speciality `orm:"rel(m2m)"`
	City             *City         `orm:"rel(fk)"`
	CreatedAt        time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt        time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(University))
}

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

	// Add speciality to university
	_, err := o.QueryM2M(university, "Specialities").Add(speciality)
	return err
}
