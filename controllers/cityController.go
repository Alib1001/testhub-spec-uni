package controllers

import (
	"encoding/json"
	"log"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

type CityController struct {
	beego.Controller
}

// Create adds a new city to the database.
// @Title Create
// @Description Создание нового города.
// @Param	body	body	models.City	true	"JSON с данными о городе"
// @Success 200 {object} map[string]int64 {"id": 1} "ID созданного города"
// @Failure 400 {string} string "400 ошибка разбора JSON или другая ошибка"
// @router / [post]
func (c *CityController) Create() {
	var city models.City

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &city)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	id, err := models.AddCity(&city)
	if err == nil {
		city.Id = int(id)
		if err := models.IndexCity(&city); err != nil {
			log.Printf("Failed to index city: %v", err)
		}

		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Get returns information about a city by its ID.
// @Title Get
// @Description Получение информации о городе по ID.
// @Param	id		path	int	true	"ID города для получения информации"
// @Success 200 {object} models.City "Информация о городе"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:id [get]
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

// GetAll returns a list of all cities.
// @Title GetAll
// @Description Получение списка всех городов.
// @Success 200 {array} models.City "Список городов"
// @Failure 400 {string} string "400 ошибка получения списка или другая ошибка"
// @router / [get]
func (c *CityController) GetAll() {
	cities, err := models.GetAllCities()
	if err == nil {
		c.Data["json"] = cities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update updates information about a city by its ID.
// @Title Update
// @Description Обновление информации о городе по ID.
// @Param	id		path	int	true	"ID города для обновления информации"
// @Param	body	body	models.City	true	"JSON с обновленными данными о городе"
// @Success 200 string "Обновление успешно выполнено"
// @Failure 400 {string} string "400 некорректный ID, ошибка разбора JSON или другая ошибка"
// @router /:id [put]
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

// Delete deletes a city by its ID.
// @Summary Delete a city
// @Description Удаление города по ID
// @ID delete-city-by-id
// @Param id path int true "ID города для удаления"
// @Success 200 {string} string "Успешное удаление"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @Router /:id [delete]
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

// GetWithUniversities возвращает информацию о городе вместе с университетами по его ID.
// @Title GetWithUniversities
// @Description Получение информации о городе с университетами по ID.
// @Param	id		path	int	true	"ID города для получения информации"
// @Success 200 {object} models.City "Информация о городе с университетами"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /info/:id [get]
func (c *CityController) GetWithUniversities() {
	id, _ := c.GetInt(":id")
	city, err := models.GetCityWithUniversities(id)
	if err == nil {
		c.Data["json"] = city
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// SearchCities ищет города по имени.
// @Title SearchCities
// @Description Поиск городов по имени.
// @Param	name		query	string	true	"Имя города для поиска"
// @Success 200 {array} models.City "Список найденных городов"
// @Failure 400 {string} string "400 ошибка поиска или другая ошибка"
// @router /search [get]
func (c *CityController) SearchCities() {
	name := c.GetString("name")
	cities, err := models.SearchCitiesByName(name)
	if err == nil {
		c.Data["json"] = cities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
