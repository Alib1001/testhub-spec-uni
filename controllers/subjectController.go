package controllers

import (
	"encoding/json"
	"strconv"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// SubjectController обрабатывает запросы для работы с предметами.
type SubjectController struct {
	beego.Controller
}

// Create добавляет новый предмет в базу данных.
// @Title Create
// @Description Создание нового предмета.
// @Param	body	body	models.Subject	true	"JSON с данными о предмете"
// @Success 200 {object} map[string]int64	"ID созданного предмета"
// @Failure 400 ошибка разбора JSON или другая ошибка
// @router / [post]
func (c *SubjectController) Create() {
	var subject models.Subject

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &subject)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	id, err := models.AddSubject(&subject)
	if err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Get возвращает информацию о предмете по его ID.
// @Title Get
// @Description Получение информации о предмете по ID.
// @Param	id		path	int	true	"ID предмета для получения информации"
// @Success 200 {object} models.Subject	"Информация о предмете"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
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

// GetAll возвращает список всех предметов.
// @Title GetAll
// @Description Получение списка всех предметов.
// @Success 200 {array} models.Subject	"Список предметов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *SubjectController) GetAll() {
	subjects, err := models.GetAllSubjects()
	if err == nil {
		c.Data["json"] = subjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update обновляет информацию о предмете по его ID.
// @Title Update
// @Description Обновление информации о предмете по ID.
// @Param	id		path	int	true	"ID предмета для обновления информации"
// @Param	body	body	models.Subject	true	"JSON с обновленными данными о предмете"
// @Success 200 string	"Обновление успешно выполнено"
// @Failure 400 некорректный ID, ошибка разбора JSON или другая ошибка
// @router /:id [put]
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

// Delete удаляет предмет по его ID.
// @Title Delete
// @Description Удаление предмета по ID.
// @Param	id		path	int	true	"ID предмета для удаления"
// @Success 200 string	"Удаление успешно выполнено"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [delete]
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

func (c *SubjectController) SearchSubjectsByName() {
	name := c.GetString("name")
	subjects, err := models.SearchSubjectsByName(name)
	if err == nil {
		c.Data["json"] = subjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetAllowedSecondSubjects возвращает список допустимых вторых предметов.
// @Title GetAllowedSecondSubjects
// @Description Получение списка допустимых вторых предметов по ID первого предмета.
// @Param	firstSubjectId	query	int	true	"ID первого предмета"
// @Success 200 {array} models.Subject	"Список допустимых вторых предметов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router /second_subjects [get]
func (c *SubjectController) GetAllowedSecondSubjects() {
	firstSubjectIdStr := c.GetString(":firstSubjectId")
	if firstSubjectIdStr == "" {
		c.Data["json"] = "firstSubjectId parameter is required"
		c.ServeJSON()
		return
	}

	firstSubjectId, err := strconv.Atoi(firstSubjectIdStr)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	subjects, err := models.GetAllowedSecondSubjects(firstSubjectId)
	if err == nil {
		c.Data["json"] = subjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
