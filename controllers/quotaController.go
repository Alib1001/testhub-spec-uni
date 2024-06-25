package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// QuotaController обрабатывает запросы для работы с квотами.
type QuotaController struct {
	beego.Controller
}

// Create adds a new quota to the database.
// @Title Create
// @Description Создание новой квоты.
// @Param	body	body	models.Quota	true	"JSON с данными о квоте"
// Success 200 {object} map[string]int64	"ID созданной квоты"
// Failure 400 {string} string "400 ошибка разбора JSON или другая ошибка"
// @router / [post]
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

// Get returns information about a quota by its ID.
// @Title Get
// @Description Получение информации о квоте по ID.
// @Param	id		path	int	true	"ID квоты для получения информации"
// @Success 200 {object} models.Quota	"Информация о квоте"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:id [get]
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

// GetAll returns a list of all quotas.
// @Title GetAll
// @Description Получение списка всех квот.
// @Success 200 {array} models.Quota	"Список квот"
// @Failure 400 {string} string "400 ошибка получения списка или другая ошибка"
// @router / [get]
func (c *QuotaController) GetAll() {
	quotas, err := models.GetAllQuotas()
	if err == nil {
		c.Data["json"] = quotas
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update updates information about a quota by its ID.
// @Title Update
// @Description Обновление информации о квоте по ID.
// @Param	id		path	int	true	"ID квоты для обновления информации"
// @Param	body	body	models.Quota	true	"JSON с обновленными данными о квоте"
// Success 200 string "Обновление успешно выполнено"
// Failure 400 {string} string "400 некорректный ID, ошибка разбора JSON или другая ошибка"
// @router /:id [put]
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

// Delete deletes a quota by its ID.
// @Title Delete
// @Description Удаление квоты по ID.
// @Param	id		path	int	true	"ID квоты для удаления"
// Success 200 string "Удаление успешно выполнено"
// Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:id [delete]
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

// AddSpecialtyToQuota adds a specialty to a quota.
// @Title AddSpecialtyToQuota
// @Description Добавление специальности к квоте.
// @Param	specialty_id	body	int	true	"ID специальности"
// @Param	quota_id	body	int	true	"ID квоты"
// Success 200 string	"Специальность успешно добавлена к квоте"
// Failure 400 {string} string "400 ошибка разбора JSON или другая ошибка"
// @router /add-specialty [post]
func (c *QuotaController) AddSpecialtyToQuota() {
	var input struct {
		SpecialtyId int `json:"specialty_id"`
		QuotaId     int `json:"quota_id"`
	}

	// Получение тела запроса с помощью CopyBody()
	requestBody := c.Ctx.Input.CopyBody(1024)

	// Распаковка JSON из тела запроса
	err := json.Unmarshal(requestBody, &input)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Вызов метода добавления специальности к квоте
	err = models.AddSpecialtyToQuota(input.SpecialtyId, input.QuotaId)
	if err == nil {
		c.Data["json"] = "Specialty added to quota successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
