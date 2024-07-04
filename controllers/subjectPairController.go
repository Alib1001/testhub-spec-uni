package controllers

import (
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// SubjectPairController обрабатывает запросы для работы с парами предметов.
type SubjectPairController struct {
	beego.Controller
}

// Add добавляет новую пару предметов в базу данных.
// @Title Add
// @Description Создание новой пары предметов.
// @Param	firstSubjectId	query	int	true	"ID первого предмета"
// @Param	secondSubjectId	query	int	true	"ID второго предмета"
// @Success 200 {object} map[string]int64	"ID созданной пары предметов"
// @Failure 400 некорректные ID или другая ошибка
// @router /add/:firstSubjectId/:secondSubjectId [post]
func (c *SubjectPairController) Add() {
	firstSubjectId, err1 := c.GetInt(":firstSubjectId")
	secondSubjectId, err2 := c.GetInt(":secondSubjectId")
	if err1 != nil || err2 != nil {
		c.Data["json"] = "Invalid subject IDs"
		c.ServeJSON()
		return
	}

	var subjectPair models.SubjectPair
	subjectPair.Subject1 = &models.Subject{Id: firstSubjectId}
	subjectPair.Subject2 = &models.Subject{Id: secondSubjectId}

	if id, err := models.AddSubjectPair(&subjectPair); err == nil {
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Get возвращает информацию о паре предметов по ее ID.
// @Title Get
// @Description Получение информации о паре предметов по ID.
// @Param	id		path	int	true	"ID пары предметов для получения информации"
// @Success 200 {object} models.SubjectPair	"Информация о паре предметов"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
func (c *SubjectPairController) Get() {
	id, _ := c.GetInt(":id")
	if subjectPair, err := models.GetSubjectPairById(id); err == nil {
		c.Data["json"] = subjectPair
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetAll возвращает список всех пар предметов.
// @Title GetAll
// @Description Получение списка всех пар предметов.
// @Success 200 {array} models.SubjectPair	"Список пар предметов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *SubjectPairController) GetAll() {
	if subjectPairs, err := models.GetAllSubjectPairs(); err == nil {
		c.Data["json"] = subjectPairs
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update обновляет информацию о паре предметов по ее ID.
// @Title Update
// @Description Обновление информации о паре предметов по ID.
// @Param	id		path	int	true	"ID пары предметов для обновления информации"
// @Param	firstSubjectId	query	int	true	"ID первого предмета"
// @Param	secondSubjectId	query	int	true	"ID второго предмета"
// @Success 200 string	"Обновление успешно выполнено"
// @Failure 400 некорректные ID или другая ошибка
// @router /:id/:firstSubjectId/:secondSubjectId [put]
func (c *SubjectPairController) Update() {
	id, _ := c.GetInt(":id")
	firstSubjectId, err1 := c.GetInt(":firstSubjectId")
	secondSubjectId, err2 := c.GetInt(":secondSubjectId")
	if err1 != nil || err2 != nil {
		c.Data["json"] = "Invalid subject IDs"
		c.ServeJSON()
		return
	}

	var subjectPair models.SubjectPair
	subjectPair.Id = id
	subjectPair.Subject1 = &models.Subject{Id: firstSubjectId}
	subjectPair.Subject2 = &models.Subject{Id: secondSubjectId}

	if err := models.UpdateSubjectPair(&subjectPair); err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete удаляет пару предметов по ее ID.
// @Title Delete
// @Description Удаление пары предметов по ID.
// @Param	id		path	int	true	"ID пары предметов для удаления"
// @Success 200 string	"Удаление успешно выполнено"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [delete]
func (c *SubjectPairController) Delete() {
	id, _ := c.GetInt(":id")
	if err := models.DeleteSubjectPair(id); err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetBySubjectIds возвращает ID пары предметов по ID первого и второго предмета.
// @Title GetBySubjectIds
// @Description Получение ID пары предметов по ID первого и второго предмета.
// @Param	firstSubjectId	path	int	true	"ID первого предмета"
// @Param	secondSubjectId	path	int	true	"ID второго предмета"
// @Success 200 {object} map[string]int	"ID пары предметов"
// @Failure 400 некорректные ID или другая ошибка
// @router /get/:firstSubjectId/:secondSubjectId [get]
func (c *SubjectPairController) GetBySubjectIds() {
	firstSubjectId, err1 := c.GetInt(":firstSubjectId")
	secondSubjectId, err2 := c.GetInt(":secondSubjectId")
	if err1 != nil || err2 != nil {
		c.Data["json"] = "Invalid subject IDs"
		c.ServeJSON()
		return
	}

	subjectPair, err := models.GetSubjectPairBySubjectIds(firstSubjectId, secondSubjectId)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]int{"id": subjectPair.Id}
	}
	c.ServeJSON()
}
