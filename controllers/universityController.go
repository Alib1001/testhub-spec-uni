package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// UniversityController обрабатывает запросы для работы с университетами.
type UniversityController struct {
	beego.Controller
}

// Create добавляет новый университет в базу данных.
// @Title Create
// @Description Создание нового университета.
// @Param	body	body	models.University	true	"JSON с данными о университете"
// @Success 200 {object} map[string]int64	"ID созданного университета"
// @Failure 400 ошибка разбора JSON или другая ошибка
// @router / [post]
func (c *UniversityController) Create() {
	var university models.University

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &university)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	id, err := models.AddUniversity(&university)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Get возвращает информацию о университете по его ID.
// @Title Get
// @Description Получение информации о университете по ID.
// @Param	id		path	int	true	"ID университета для получения информации"
// @Success 200 {object} models.University	"Информация о университете"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
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

// GetAll возвращает список всех университетов.
// @Title GetAll
// @Description Получение списка всех университетов.
// @Success 200 {array} models.University	"Список университетов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *UniversityController) GetAll() {
	universities, err := models.GetAllUniversities()
	if err == nil {
		c.Data["json"] = universities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update обновляет информацию о университете по его ID.
// @Title Update
// @Description Обновление информации о университете по ID.
// @Param	id		path	int	true	"ID университета для обновления информации"
// @Param	body	body	models.University	true	"JSON с обновленными данными о университете"
// @Success 200 string	"Обновление успешно выполнено"
// @Failure 400 некорректный ID, ошибка разбора JSON или другая ошибка
// @router /:id [put]
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

// Delete удаляет университет по его ID.
// @Title Delete
// @Description Удаление университета по ID.
// @Param	id		path	int	true	"ID университета для удаления"
// @Success 200 string	"Удаление успешно выполнено"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [delete]
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
