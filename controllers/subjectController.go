package controllers

import (
	"encoding/json"
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
	_ = c.Ctx.Input.CopyBody(1024)
	var subject models.Subject
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &subject); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	if id, err := models.AddSubject(&subject); err == nil {
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
	if subject, err := models.GetSubjectById(id); err == nil {
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
	if subjects, err := models.GetAllSubjects(); err == nil {
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
	_ = c.Ctx.Input.CopyBody(1024)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &subject); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	subject.Id = id
	if err := models.UpdateSubject(&subject); err == nil {
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
	if err := models.DeleteSubject(id); err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// SearchSubjectsByName ищет предметы по имени.
// @Title SearchSubjectsByName
// @Description Поиск предметов по имени.
// @Param	name	query	string	true	"Имя предмета для поиска"
// @Success 200 {array} models.Subject	"Список найденных предметов"
// @Failure 400 ошибка выполнения поиска
// @router /search [get]
func (c *SubjectController) SearchSubjectsByName() {
	name := c.GetString("name")
	if subjects, err := models.SearchSubjectsByName(name); err == nil {
		c.Data["json"] = subjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetAllowedSecondSubjects возвращает список предметов, соответствующих первому предмету.
// @Title GetAllowedSecondSubjects
// @Description Получение списка предметов, соответствующих первому предмету.
// @Param	firstSubjectId	path	int	true	"ID первого предмета"
// @Success 200 {array} models.Subject	"Список предметов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router /secubjects/:firstSubjectId [get]
func (c *SubjectController) GetAllowedSecondSubjects() {
	subject1Id, err := c.GetInt(":firstSubjectId")
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	if subjects, err := models.GetAllowedSecondSubjects(subject1Id); err == nil {
		c.Data["json"] = subjects
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
