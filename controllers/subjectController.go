package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type SubjectController struct {
	beego.Controller
}

func (c *SubjectController) Create() {
	var subject models.Subject
	json.Unmarshal(c.Ctx.Input.RequestBody, &subject)
	id, err := models.AddSubject(&subject)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SubjectController) Get() {
	id, _ := c.GetInt(":id")
	subject, err := models.GetSubjectById(id)
	if err == nil {
		c.Data["json"] = subject
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SubjectController) GetAll() {
	subjects, err := models.GetAllSubjects()
	if err == nil {
		c.Data["json"] = subjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SubjectController) Update() {
	id, _ := c.GetInt(":id")
	var subject models.Subject
	json.Unmarshal(c.Ctx.Input.RequestBody, &subject)
	subject.Id = id
	err := models.UpdateSubject(&subject)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *SubjectController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteSubject(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
