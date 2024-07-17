package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type PointStat struct {
	Id            int `orm:"auto"`
	AnnualGrants  int
	MinScore      int
	MinGrantScore int
	Year          int
	Speciality    *Speciality `orm:"rel(fk)"`
	University    *University `orm:"rel(fk)"`
	CreatedAt     time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt     time.Time   `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(PointStat))
}
func AddPointStat(universityId, specialityId int, pointStat *PointStat) (int64, error) {
	pointStat.University.Id = universityId
	pointStat.Speciality.Id = specialityId

	o := orm.NewOrm()
	id, err := o.Insert(pointStat)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetPointStatsByUniversityAndSpeciality(universityId, specialityId int) ([]*PointStat, error) {
	o := orm.NewOrm()
	var pointStats []*PointStat
	_, err := o.QueryTable("point_stat").
		Filter("University__Id", universityId).
		Filter("Speciality__Id", specialityId).
		All(&pointStats)
	return pointStats, err
}
