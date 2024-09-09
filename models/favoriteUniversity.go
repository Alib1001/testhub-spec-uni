package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type UserInfo struct {
	ID        int    `orm:"column(id);auto"`
	UserID    int    `orm:"column(user_id)"`
	UUID      string `orm:"column(uuid)"`
	Role      string `orm:"column(role)"`
	FirstName string `orm:"column(first_name)"`
	LastName  string `orm:"column(last_name)"`
	Balance   int    `orm:"column(balance)"`
}

type FavoriteUniversity struct {
	Id         int         `orm:"auto"`
	UserId     int         `orm:"column(user_id)"`
	University *University `orm:"rel(fk);on_delete(cascade)"`
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(FavoriteUniversity))
}

func AddFavoriteUniversity(userId int, universityId int) error {
	o := orm.NewOrm()

	favorite := &FavoriteUniversity{
		UserId:     userId,
		University: &University{Id: universityId},
	}

	_, err := o.Insert(favorite)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFavoriteUniversity(userId int, universityId int) error {
	o := orm.NewOrm()

	_, err := o.QueryTable("favorite_university").Filter("user_id", userId).Filter("University__Id", universityId).Delete()
	if err != nil {
		return err
	}
	return nil
}

func ListFavoriteUniversities(userId int) ([]*University, error) {
	o := orm.NewOrm()
	var favorites []*FavoriteUniversity

	_, err := o.QueryTable("favorite_university").Filter("user_id", userId).RelatedSel().All(&favorites)
	if err != nil {
		return nil, err
	}

	var universities []*University
	for _, favorite := range favorites {
		universities = append(universities, favorite.University)
	}
	return universities, nil
}
