package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Quota struct {
	Id           int    `orm:"auto"`
	QuotaType    string `orm:"size(64)"`
	QuotaTypeRu  string `orm:"size(64)"`
	QuotaTypeKz  string `orm:"size(64)"`
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

func GetQuotaById(id int, language string) (*Quota, error) {
	o := orm.NewOrm()
	quota := &Quota{Id: id}
	err := o.Read(quota)
	if err != nil {
		return nil, err
	}

	switch language {
	case "ru":
		quota.QuotaType = quota.QuotaTypeRu
	case "kz":
		quota.QuotaType = quota.QuotaTypeKz
	}

	return quota, err
}

func GetAllQuotas(language string) ([]*Quota, error) {
	o := orm.NewOrm()
	var quotas []*Quota
	_, err := o.QueryTable("quota").All(&quotas)
	if err != nil {
		return nil, err
	}

	for _, quota := range quotas {
		switch language {
		case "ru":
			quota.QuotaType = quota.QuotaTypeRu
		case "kz":
			quota.QuotaType = quota.QuotaTypeKz
		}
	}

	return quotas, err
}

func UpdateQuota(quota *Quota, fields ...string) error {
	o := orm.NewOrm()
	_, err := o.Update(quota, fields...)
	if err != nil {
		return err
	}
	return nil
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

func GetAllQuotasWithSpecialities(language string) ([]*Quota, error) {
	o := orm.NewOrm()
	var quotas []*Quota
	_, err := o.QueryTable("quota").RelatedSel("Specialities").All(&quotas)
	if err != nil {
		return nil, err
	}

	for _, quota := range quotas {
		switch language {
		case "ru":
			quota.QuotaType = quota.QuotaTypeRu
		case "kz":
			quota.QuotaType = quota.QuotaTypeKz
		}
	}

	return quotas, err
}

func GetQuotaWithSpecialitiesById(id int, language string) (*Quota, error) {
	o := orm.NewOrm()
	quota := &Quota{Id: id}

	err := o.Read(quota)
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(quota, "Specialities")
	if err != nil {
		return nil, err
	}

	switch language {
	case "ru":
		quota.QuotaType = quota.QuotaTypeRu
	case "kz":
		quota.QuotaType = quota.QuotaTypeKz
	}

	return quota, nil
}
