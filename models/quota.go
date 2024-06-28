package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Quota struct {
	Id           int    `orm:"auto"`
	QuotaType    string `orm:"size(64)"`
	Count        int
	MinScore     int
	MaxScore     int
	Specialities []*Speciality `orm:"rel(m2m);rel_table(quota_specialities)"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Quota))
}

func AddQuota(quota *Quota) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(quota)
	return id, err
}

func GetQuotaById(id int) (*Quota, error) {
	o := orm.NewOrm()
	quota := &Quota{Id: id}
	err := o.Read(quota)
	return quota, err
}

func GetAllQuotas() ([]*Quota, error) {
	o := orm.NewOrm()
	var quotas []*Quota
	_, err := o.QueryTable("quota").All(&quotas)
	return quotas, err
}

func UpdateQuota(quota *Quota) error {
	o := orm.NewOrm()
	_, err := o.Update(quota)
	return err
}

func DeleteQuota(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Quota{Id: id})
	return err
}
func AddSpecialityToQuota(specialityId, quotaId int) error {
	o := orm.NewOrm()

	speciality := &Speciality{Id: specialityId}
	if err := o.Read(speciality); err != nil {
		return err
	}

	quota := &Quota{Id: quotaId}
	if err := o.Read(quota); err != nil {
		return err
	}

	m2m := o.QueryM2M(quota, "Specialities")
	exists := m2m.Exist(speciality)
	if exists {
		return fmt.Errorf("Speciality with ID %d already exists in quota with ID %d", specialityId, quotaId)
	}

	_, err := m2m.Add(speciality)
	return err
}
func GetAllQuotasWithSpecialities() ([]*Quota, error) {
	o := orm.NewOrm()
	var quotas []*Quota
	_, err := o.QueryTable("quota").RelatedSel("Specialities").All(&quotas)
	return quotas, err
}
func GetQuotaWithSpecialitiesById(id int) (*Quota, error) {
	o := orm.NewOrm()
	quota := &Quota{Id: id}

	// Загрузка квоты по ID
	err := o.Read(quota)
	if err != nil {
		return nil, err
	}

	// Предзагрузка списка специальностей для данной квоты
	_, err = o.LoadRelated(quota, "Specialities")
	if err != nil {
		return nil, err
	}

	return quota, nil
}
