package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"testhub-spec-uni/models"
	"time"
)

type SpecialityUniversityController struct {
	beego.Controller
}

// @Title CreateUniversitySpecialityDetail
// @Description create SpecialityUniversity
// @Param   university_id   path      int   true   "University ID"
// @Param   speciality_id   path      int   true   "Speciality ID"
// @Param   term            formData  int   true   "Term"
// @Param   edu_lang        formData  string   true   "EduLang"
// @Success 201 {int} models.SpecialityUniversity
// @Failure 400 {string} invalid input
// @router /:uid/:sid [post]
func (c *SpecialityUniversityController) CreateUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")
	term, _ := c.GetInt("term")
	eduLang := c.GetString("edu_lang")

	detail := models.SpecialityUniversity{
		University: &models.University{Id: universityID},
		Speciality: &models.Speciality{Id: specialityID},
		Term:       term,
		EduLang:    eduLang,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := detail.Create(); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = detail
	}
	c.ServeJSON()
}

// @Title UpdateUniversitySpecialityDetail
// @Description update SpecialityUniversity by id
// @Param   university_id   path      int   true   "University ID"
// @Param   speciality_id   path      int   true   "Speciality ID"
// @Param   term            formData  int   false  "Term"
// @Param   edu_lang        formData  string   false  "EduLang"
// @Success 200 {string} update success!
// @Failure 400 {string} invalid input
// @router /:uid/:sid [put]
func (c *SpecialityUniversityController) UpdateUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")

	// Получаем текущее значение записи
	detail, err := models.GetByUniversityAndSpeciality(universityID, specialityID)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Record not found"}
		c.ServeJSON()
		return
	}

	// Проверяем и обновляем только те поля, которые были переданы в запросе
	if term, err := c.GetInt("term"); err == nil {
		detail.Term = term
	}

	if eduLang := c.GetString("edu_lang"); eduLang != "" {
		detail.EduLang = eduLang
	}

	detail.UpdatedAt = time.Now()

	if err := detail.Update(); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = "update success!"
	}

	c.ServeJSON()
}

// @Title GetUniversitySpecialityDetail
// @Description get SpecialityUniversity by university_id and speciality_id
// @Param   university_id   path    int  true   "University ID"
// @Param   speciality_id   path    int  true   "Speciality ID"
// @Success 200 {object} models.SpecialityUniversity
// @Failure 404 {string} record not found
// @router /:uid/:sid [get]
func (c *SpecialityUniversityController) GetUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")

	detail, err := models.GetByUniversityAndSpeciality(universityID, specialityID)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "Record not found"}
	} else {
		c.Data["json"] = detail
	}
	c.ServeJSON()
}

// @Title DeleteUniversitySpecialityDetail
// @Description delete SpecialityUniversity by university_id and speciality_id
// @Param   university_id   path    int  true   "University ID"
// @Param   speciality_id   path    int  true   "Speciality ID"
// @Success 200 {string} delete success!
// @Failure 404 {string} record not found
// @router /:uid/:sid [delete]
func (c *SpecialityUniversityController) DeleteUniversitySpecialityDetail() {
	universityID, _ := c.GetInt(":uid")
	specialityID, _ := c.GetInt(":sid")

	if err := models.DeleteByUniversityAndSpeciality(universityID, specialityID); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "Record not found"}
	} else {
		c.Data["json"] = "delete success!"
	}
	c.ServeJSON()
}
