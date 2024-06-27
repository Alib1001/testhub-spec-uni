package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int       `orm:"auto;pk"`
	Username        string    `orm:"size(100);unique"`
	Password        string    `orm:"size(100)"`
	FirstName       string    `orm:"size(100)"`
	LastName        string    `orm:"size(100)"`
	Role            string    `orm:"size(50)"`
	ProfileImageURL string    `orm:"size(255)"`
	CreatedAt       time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt       time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AddUser(user User) error {
	o := orm.NewOrm()
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	_, err = o.Insert(&user)
	return err
}

func UpdateUser(user User) error {
	o := orm.NewOrm()
	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	_, err := o.Update(&user)
	return err
}

func GetUserByUsername(username string) (User, error) {
	o := orm.NewOrm()
	user := User{Username: username}
	err := o.Read(&user, "Username")
	return user, err
}

func DeleteUser(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&User{Id: id})
	return err
}

func GetAllUsers() ([]User, error) {
	o := orm.NewOrm()
	var users []User
	_, err := o.QueryTable("user").All(&users)
	return users, err
}

func GetUserById(id int) (User, error) {
	o := orm.NewOrm()
	user := User{Id: id}
	err := o.Read(&user)
	return user, err
}
