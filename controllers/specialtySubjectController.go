package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type SpecialtySubjectController struct {
	beego.Controller
}

func (c *SpecialtySubjectController) Create() {
	var specialtySubject models.SpecialtySubject
	json.Unmarshal(c.Ctx.Input.RequestBody, &specialtySubject)
	id, err := models.AddSpecialtySubject(&specialtySubject)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtySubjectController) Get() {
	id, _ := c.GetInt(":id")
	specialtySubject, err := models.GetSpecialtySubjectById(id)
	if err == nil {
		c.Data["json"] = specialtySubject
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtySubjectController) GetAll() {
	specialtySubjects, err := models.GetAllSpecialtySubjects()
	if err == nil {
		c.Data["json"] = specialtySubjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtySubjectController) Update() {
	id, _ := c.GetInt(":id")
	var specialtySubject models.SpecialtySubject
	json.Unmarshal(c.Ctx.Input.RequestBody, &specialtySubject)
	specialtySubject.Id = id
	err := models.UpdateSpecialtySubject(&specialtySubject)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtySubjectController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteSpecialtySubject(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
