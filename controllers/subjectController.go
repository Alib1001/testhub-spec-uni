package controllers

import (
	"encoding/json"
	"net/http"
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
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} models.SubjectResponse	"Информация о предмете"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
func (c *SubjectController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid subject ID")
		return
	}

	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	subject, err := models.GetSubjectById(id, language)
	if err != nil {
		c.CustomAbort(http.StatusNotFound, err.Error())
		return
	}

	c.Data["json"] = subject
	c.ServeJSON()
}

// GetAll возвращает список всех предметов.
// @Title GetAll
// @Description Получение списка всех предметов.
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.SubjectResponse	"Список предметов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *SubjectController) GetAll() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	subjects, err := models.GetAllSubjects(language)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = subjects
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
	idStr := c.Ctx.Input.Param(":id")
	_ = c.Ctx.Input.CopyBody(1024)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid subject ID")
		return
	}

	var subject models.Subject
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &subject); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	subject.Id = id

	if err := models.UpdateSubject(&subject); err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = "Update successful"
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
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid subject ID")
		return
	}

	if err := models.DeleteSubject(id); err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = "Delete successful"
	c.ServeJSON()
}

// SearchSubjectsByName ищет предметы по имени.
// @Title SearchSubjectsByName
// @Description Поиск предметов по имени.
// @Param	name	query	string	true	"Имя предмета для поиска"
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.SubjectResponse	"Список найденных предметов"
// @Failure 400 ошибка поиска или другая ошибка
// @router /search [get]
func (c *SubjectController) SearchSubjectsByName() {
	name := c.GetString("name")
	if name == "" {
		c.CustomAbort(http.StatusBadRequest, "Search name is required")
		return
	}

	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	subjects, err := models.SearchSubjectsByName(name, language)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = subjects
	c.ServeJSON()
}

// GetAllowedSecondSubjects возвращает список предметов, соответствующих первому предмету.
// @Title GetAllowedSecondSubjects
// @Description Получение списка предметов, соответствующих первому предмету.
// @Param	firstSubjectId	path	int	true	"ID первого предмета"
// @Param	lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.SubjectResponse	"Список предметов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router /secubjects/:firstSubjectId [get]
func (c *SubjectController) GetAllowedSecondSubjects() {
	subject1Id, err := c.GetInt(":firstSubjectId")
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	subjects, err := models.GetAllowedSecondSubjects(subject1Id)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	var subjectResponses []*models.SubjectResponse
	for _, subject := range subjects {
		name := subject.Name
		switch language {
		case "ru":
			name = subject.NameRu
		case "kz":
			name = subject.NameKz
		}
		subjectResponses = append(subjectResponses, &models.SubjectResponse{
			Id:   subject.Id,
			Name: name,
		})
	}

	c.Data["json"] = subjectResponses
	c.ServeJSON()
}
