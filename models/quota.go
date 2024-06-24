package models

import (
	"github.com/astaxie/beego/orm"
)

type Quota struct {
	Id          int    `orm:"auto"`
	QuotaType   string `orm:"size(64)"`
	SpecialtyId int
	MinScore    int
	MaxScore    int
}

func init() {
	orm.RegisterModel(new(Quota))
}

// CRUD methods

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
