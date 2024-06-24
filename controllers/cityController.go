package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type CityController struct {
	beego.Controller
}

func (c *CityController) Create() {
	var city models.City

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &city)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Добавление города в базу данных
	id, err := models.AddCity(&city)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *CityController) Get() {
	id, _ := c.GetInt(":id")
	city, err := models.GetCityById(id)
	if err == nil {
		c.Data["json"] = city
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *CityController) GetAll() {
	cities, err := models.GetAllCities()
	if err == nil {
		c.Data["json"] = cities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *CityController) Update() {
	id, _ := c.GetInt(":id")
	var city models.City
	json.Unmarshal(c.Ctx.Input.RequestBody, &city)
	city.Id = id
	err := models.UpdateCity(&city)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *CityController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteCity(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
