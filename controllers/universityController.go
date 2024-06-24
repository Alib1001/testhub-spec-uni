package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type UniversityController struct {
	beego.Controller
}

func (c *UniversityController) Create() {
	var university models.University
	json.Unmarshal(c.Ctx.Input.RequestBody, &university)
	id, err := models.AddUniversity(&university)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *UniversityController) Get() {
	id, _ := c.GetInt(":id")
	university, err := models.GetUniversityById(id)
	if err == nil {
		c.Data["json"] = university
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *UniversityController) GetAll() {
	universities, err := models.GetAllUniversities()
	if err == nil {
		c.Data["json"] = universities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *UniversityController) Update() {
	id, _ := c.GetInt(":id")
	var university models.University
	json.Unmarshal(c.Ctx.Input.RequestBody, &university)
	university.Id = id
	err := models.UpdateUniversity(&university)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *UniversityController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteUniversity(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
