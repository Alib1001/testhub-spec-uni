package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type QuotaController struct {
	beego.Controller
}

func (c *QuotaController) Create() {
	var quota models.Quota

	// Получение тела запроса с помощью CopyBody()
	requestBody := c.Ctx.Input.CopyBody(1024)

	// Распаковка JSON из тела запроса
	err := json.Unmarshal(requestBody, &quota)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Добавление квоты в базу данных
	id, err := models.AddQuota(&quota)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *QuotaController) Get() {
	id, _ := c.GetInt(":id")
	quota, err := models.GetQuotaById(id)
	if err == nil {
		c.Data["json"] = quota
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *QuotaController) GetAll() {
	quotas, err := models.GetAllQuotas()
	if err == nil {
		c.Data["json"] = quotas
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *QuotaController) Update() {
	id, _ := c.GetInt(":id")
	var quota models.Quota
	json.Unmarshal(c.Ctx.Input.RequestBody, &quota)
	quota.Id = id
	err := models.UpdateQuota(&quota)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *QuotaController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteQuota(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
