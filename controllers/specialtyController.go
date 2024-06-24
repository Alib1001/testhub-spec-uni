package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type SpecialtyController struct {
	beego.Controller
}

func (c *SpecialtyController) Create() {
	var specialty models.Specialty
	json.Unmarshal(c.Ctx.Input.RequestBody, &specialty)
	id, err := models.AddSpecialty(&specialty)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyController) Get() {
	id, _ := c.GetInt(":id")
	specialty, err := models.GetSpecialtyById(id)
	if err == nil {
		c.Data["json"] = specialty
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyController) GetAll() {
	specialties, err := models.GetAllSpecialties()
	if err == nil {
		c.Data["json"] = specialties
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyController) Update() {
	id, _ := c.GetInt(":id")
	var specialty models.Specialty
	json.Unmarshal(c.Ctx.Input.RequestBody, &specialty)
	specialty.Id = id
	err := models.UpdateSpecialty(&specialty)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SpecialtyController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteSpecialty(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
