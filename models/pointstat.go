package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type PointStat struct {
	Id            int `orm:"auto"`
	GrantCount    int
	MinScore      int
	MinGrantScore int
	Year          int
	AvgSalary     int
	Price         int
	Speciality    *Speciality `orm:"rel(fk)"`
	University    *University `orm:"rel(fk);on_delete(cascade)"`
	CreatedAt     time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt     time.Time   `orm:"auto_now;type(datetime)"`
}

type GetPointStatResponse struct {
	Id            int `orm:"auto"`
	GrantCount    int `json:"grant_count"`
	MinScore      int `json:"min_score"`
	MinGrantScore int `json:"min_grant_score"`
	Year          int `json:"year"`
	AvgSalary     int `json:"avg_salary"`
	Price         int `json:"price"`
}

type AddPointStatResponse struct {
	Id            int `form:"Id"`
	AnnualGrants  int `form:"GrantCount" validate:"required"`
	MinScore      int `form:"MinScore" validate:"required"`
	MinGrantScore int `form:"MinGrantScore" validate:"required"`
	Year          int `form:"Year" validate:"required"`
	AvgSalary     int `form:"AvgSalary" validate:"required"`
	Price         int `form:"Price" validate:"required"`
}

type UpdatePointStatResponse struct {
	Id            int `form:"Id"`
	AnnualGrants  int `form:"GrantCount"`
	MinScore      int `form:"MinScore"`
	MinGrantScore int `form:"MinGrantScore"`
	Year          int `form:"Year"`
	AvgSalary     int `form:"AvgSalary"`
	Price         int `form:"Price"`
}

func init() {
	orm.RegisterModel(new(PointStat))
}
func AddPointStat(universityId, specialityId int, pointStat *PointStat) (int64, error) {
	o := orm.NewOrm()

	exists := o.QueryTable("point_stat").Filter("University__Id", universityId).Filter("Speciality__Id", specialityId).Filter("Year", pointStat.Year).Exist()
	if exists {
		return 0, fmt.Errorf("PointStat with year %d already exists for the given university and speciality", pointStat.Year)
	}

	pointStat.University.Id = universityId
	pointStat.Speciality.Id = specialityId

	id, err := o.Insert(pointStat)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func GetPointStatsByUniversityAndSpeciality(universityId, specialityId int) ([]*GetPointStatResponse, error) {
	o := orm.NewOrm()
	var pointStats []*PointStat
	_, err := o.QueryTable("point_stat").
		Filter("University__Id", universityId).
		Filter("Speciality__Id", specialityId).
		All(&pointStats)
	if err != nil {
		return nil, err
	}

	var response []*GetPointStatResponse
	for _, ps := range pointStats {
		response = append(response, &GetPointStatResponse{
			Id:            ps.Id,
			GrantCount:    ps.GrantCount,
			MinScore:      ps.MinScore,
			MinGrantScore: ps.MinGrantScore,
			Year:          ps.Year,
			AvgSalary:     ps.AvgSalary,
			Price:         ps.Price,
		})
	}

	return response, nil
}

func DeletePointStat(id int) error {
	o := orm.NewOrm()
	pointStat := PointStat{Id: id}
	if _, err := o.Delete(&pointStat); err != nil {
		return err
	}
	return nil
}

func UpdatePointStatById(id int, form *UpdatePointStatResponse) error {
	pointStat, err := GetPointStatById(id)
	if err != nil {
		return err
	}

	// Update fields only if they are set in the form
	if form.AnnualGrants != 0 {
		pointStat.GrantCount = form.AnnualGrants
	}
	if form.MinScore != 0 {
		pointStat.MinScore = form.MinScore
	}
	if form.MinGrantScore != 0 {
		pointStat.MinGrantScore = form.MinGrantScore
	}
	if form.Year != 0 {
		pointStat.Year = form.Year
	}
	if form.AvgSalary != 0 {
		pointStat.AvgSalary = form.AvgSalary
	}
	if form.Price != 0 {
		pointStat.Price = form.Price
	}

	return UpdatePointStat(pointStat)
}

func UpdatePointStat(pointStat *PointStat) error {
	o := orm.NewOrm()
	if _, err := o.Update(pointStat); err != nil {
		return err
	}
	return nil
}

func GetPointStatById(id int) (*PointStat, error) {
	o := orm.NewOrm()
	pointStat := PointStat{Id: id}
	if err := o.Read(&pointStat); err != nil {
		return nil, err
	}
	return &pointStat, nil
}
