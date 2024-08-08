package controllers

import (
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	"testhub-spec-uni/models"
	"time"
)

type UniversitySpecialityDetailController struct {
	beego.Controller
}

// @Title CreateUniversitySpecialityDetail
// @Description create UniversitySpecialityDetail
// @Param   university_id   path      int   true   "University ID"
// @Param   speciality_id   path      int   true   "Speciality ID"
// @Param   term            formData  int   true   "Term"
// @Param   edu_lang        formData  string   true   "EduLang"
// @Success 201 {int} models.UniversitySpecialityDetail
// @Failure 400 {string} invalid input
// @router /:uid/:sid [post]
func (c *UniversitySpecialityDetailController) CreateUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")
	term, _ := c.GetInt("term")
	eduLang := c.GetString("edu_lang")

	university := &models.University{Id: universityID}
	speciality := &models.Speciality{Id: specialityID}

	detail := models.UniversitySpecialityDetail{
		University: university,
		Speciality: speciality,
		Term:       term,
		EduLang:    eduLang,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	o := orm.NewOrm()
	_, err := o.Insert(&detail)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = detail
	}
	c.ServeJSON()
}

// @Title UpdateUniversitySpecialityDetail
// @Description update UniversitySpecialityDetail by id
// @Param   university_id   path      int   true   "University ID"
// @Param   speciality_id   path      int   true   "Speciality ID"
// @Param   term            formData  int   true   "Term"
// @Param   edu_lang        formData  string   true   "EduLang"
// @Success 200 {string} update success!
// @Failure 400 {string} invalid input
// @router /:uid/:sid [put]
func (c *UniversitySpecialityDetailController) UpdateUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")
	term, _ := c.GetInt("term")
	eduLang := c.GetString("edu_lang")

	o := orm.NewOrm()
	detail := models.UniversitySpecialityDetail{
		University: &models.University{Id: universityID},
		Speciality: &models.Speciality{Id: specialityID},
	}

	if err := o.QueryTable("university_speciality_detail").Filter("university_id", universityID).Filter("speciality_id", specialityID).One(&detail); err == nil {
		detail.Term = term
		detail.EduLang = eduLang
		detail.UpdatedAt = time.Now()

		if _, err := o.Update(&detail); err == nil {
			c.Data["json"] = "update success!"
		} else {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]string{"error": err.Error()}
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Record not found"}
	}
	c.ServeJSON()
}

// @Title GetUniversitySpecialityDetail
// @Description get UniversitySpecialityDetail by university_id and speciality_id
// @Param   university_id   path    int  true   "University ID"
// @Param   speciality_id   path    int  true   "Speciality ID"
// @Success 200 {object} models.UniversitySpecialityDetail
// @Failure 404 {string} record not found
// @router /:uid/:sid [get]
func (c *UniversitySpecialityDetailController) GetUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")

	o := orm.NewOrm()
	detail := models.UniversitySpecialityDetail{
		University: &models.University{Id: universityID},
		Speciality: &models.Speciality{Id: specialityID},
	}

	if err := o.QueryTable("university_speciality_detail").Filter("university_id", universityID).Filter("speciality_id", specialityID).One(&detail); err == nil {
		c.Data["json"] = detail
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "Record not found"}
	}
	c.ServeJSON()
}

// @Title DeleteUniversitySpecialityDetail
// @Description delete UniversitySpecialityDetail by university_id and speciality_id
// @Param   university_id   path    int  true   "University ID"
// @Param   speciality_id   path    int  true   "Speciality ID"
// @Success 200 {string} delete success!
// @Failure 404 {string} record not found
// @router /:uid/:sid [delete]
func (c *UniversitySpecialityDetailController) DeleteUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")

	o := orm.NewOrm()
	if _, err := o.QueryTable("university_speciality_detail").Filter("university_id", universityID).Filter("speciality_id", specialityID).Delete(); err == nil {
		c.Data["json"] = "delete success!"
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "Record not found"}
	}
	c.ServeJSON()
}
