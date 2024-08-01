package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Gallery struct {
	Id         int         `orm:"auto"`
	University *University `orm:"rel(fk);;on_delete(cascade)"`
	PhotoUrl   string      `orm:"size(256)"`
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)"`
}

type GalleryResponse struct {
	Id       int    `json:"Id"`
	PhotoUrl string `json:"PhotoUrl"`
}

func init() {
	orm.RegisterModel(new(Gallery))
}
