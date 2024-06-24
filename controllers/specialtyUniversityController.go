package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type SpecialtyUniversityController struct {
	beego.Controller
}

func (c *SpecialtyUniversityController) Create() {
	var specialtyUniversity models.SpecialtyUniversity
	json.Unmarshal(c.Ctx.Input.RequestBody, &specialtyUniversity)
	id, err := models.AddSpecialtyUniversity(&specialtyUniversity)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyUniversityController) Get() {
	id, _ := c.GetInt(":id")
	specialtyUniversity, err := models.GetSpecialtyUniversityById(id)
	if err == nil {
		c.Data["json"] = specialtyUniversity
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyUniversityController) GetAll() {
	specialtyUniversities, err := models.GetAllSpecialtyUniversities()
	if err == nil {
		c.Data["json"] = specialtyUniversities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyUniversityController) Update() {
	id, _ := c.GetInt(":id")
	var specialtyUniversity models.SpecialtyUniversity
	json.Unmarshal(c.Ctx.Input.RequestBody, &specialtyUniversity)
	specialtyUniversity.Id = id
	err := models.UpdateSpecialtyUniversity(&specialtyUniversity)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyUniversityController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteSpecialtyUniversity(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
