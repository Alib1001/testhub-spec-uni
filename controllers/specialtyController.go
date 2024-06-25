package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// SpecialtyController обрабатывает запросы для работы со специальностями.
type SpecialtyController struct {
	beego.Controller
}

// Create добавляет новую специальность в базу данных.
// @Title Create
// @Description Создание новой специальности.
// @Param	body	body	models.Specialty	true	"JSON с данными о специальности"
// Success 200 {object} map[string]int64	"ID созданной специальности"
// Failure 400 ошибка разбора JSON или другая ошибка
// @router / [post]
func (c *SpecialtyController) Create() {
	var specialty models.Specialty

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &specialty)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	id, err := models.AddSpecialty(&specialty)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Get возвращает информацию о специальности по ее ID.
// @Title Get
// @Description Получение информации о специальности по ID.
// @Param	id		path	int	true	"ID специальности для получения информации"
// Success 200 {object} models.Specialty	"Информация о специальности"
// Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
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

// GetAll возвращает список всех специальностей.
// @Title GetAll
// @Description Получение списка всех специальностей.
// Success 200 {array} models.Specialty	"Список специальностей"
// Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *SpecialtyController) GetAll() {
	specialties, err := models.GetAllSpecialties()
	if err == nil {
		c.Data["json"] = specialties
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update обновляет информацию о специальности по ее ID.
// @Title Update
// @Description Обновление информации о специальности по ID.
// @Param	id		path	int	true	"ID специальности для обновления информации"
// @Param	body	body	models.Specialty	true	"JSON с обновленными данными о специальности"
// Success 200 string	"Обновление успешно выполнено"
// Failure 400 некорректный ID, ошибка разбора JSON или другая ошибка
// @router /:id [put]
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

// Delete удаляет специальность по ее ID.
// @Title Delete
// @Description Удаление специальности по ID.
// @Param	id		path	int	true	"ID специальности для удаления"
// Success 200 string	"Удаление успешно выполнено"
// Failure 400 некорректный ID или другая ошибка
// @router /:id [delete]
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

// AddSubject добавляет предмет к специальности.
// @Title AddSubject
// @Description Добавление предмета к специальности.
// @Param	specialtyId		path	int	true	"ID специальности"
// @Param	subjectId		path	int	true	"ID предмета"
// Success 200 string	"Предмет успешно добавлен к специальности"
// Failure 400 некорректный ID или другая ошибка
// @router /:specialtyId/add-subject/:subjectId [post]
func (c *SpecialtyController) AddSubject() {
	specialtyId, _ := c.GetInt(":specialtyId")
	subjectId, _ := c.GetInt(":subjectId")

	err := models.AddSubjectToSpecialty(subjectId, specialtyId)
	if err == nil {
		c.Data["json"] = "Subject added to specialty successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
