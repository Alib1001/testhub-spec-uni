package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// SpecialityController обрабатывает запросы для работы со специальностями.
type SpecialityController struct {
	beego.Controller
}

// Create добавляет новую специальность в базу данных.
// @Title Create
// @Description Создание новой специальности.
// @Param	body	body	models.Speciality	true	"JSON с данными о специальности"
// @Success 200 {object} map[string]int64	"ID созданной специальности"
// @Failure 400 ошибка разбора JSON или другая ошибка
// @router / [post]
func (c *SpecialityController) Create() {
	var speciality models.Speciality

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &speciality)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	id, err := models.AddSpeciality(&speciality)
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
// @Success 200 {object} models.Speciality	"Информация о специальности"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
func (c *SpecialityController) Get() {
	id, _ := c.GetInt(":id")
	speciality, err := models.GetSpecialityById(id)
	if err == nil {
		c.Data["json"] = speciality
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetAll возвращает список всех специальностей.
// @Title GetAll
// @Description Получение списка всех специальностей.
// @Success 200 {array} models.Speciality	"Список специальностей"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *SpecialityController) GetAll() {
	specialties, err := models.GetAllSpecialities()
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
// @Param	body	body	models.Speciality	true	"JSON с обновленными данными о специальности"
// @Success 200 string	"Обновление успешно выполнено"
// @Failure 400 некорректный ID, ошибка разбора JSON или другая ошибка
// @router /:id [put]
func (c *SpecialityController) Update() {
	id, _ := c.GetInt(":id")
	var speciality models.Speciality
	json.Unmarshal(c.Ctx.Input.RequestBody, &speciality)
	speciality.Id = id
	err := models.UpdateSpeciality(&speciality)
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
// @Success 200 string	"Удаление успешно выполнено"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [delete]
func (c *SpecialityController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteSpeciality(id)
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
// @Param	specialityId		path	int	true	"ID специальности"
// @Param	subjectId		path	int	true	"ID предмета"
// @Success 200 string	"Предмет успешно добавлен к специальности"
// @Failure 400 некорректный ID или другая ошибка
// @router /:specialityId/add-subject/:subjectId [post]
func (c *SpecialityController) AddSubject() {
	specialityId, _ := c.GetInt(":specialityId")
	subjectId, _ := c.GetInt(":subjectId")

	err := models.AddSubjectToSpeciality(subjectId, specialityId)
	if err == nil {
		c.Data["json"] = "Subject added to speciality successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
