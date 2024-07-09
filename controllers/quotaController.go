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
// @Success 200 {object} map[string]int64	"ID созданной квоты"
// @Failure 400 {string} string "400 ошибка разбора JSON или другая ошибка"
// @router / [post]
func (c *QuotaController) Create() {
	var quota models.Quota

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &quota)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

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
// @Success 200 string "Обновление успешно выполнено"
// @Failure 400 {string} string "400 некорректный ID, ошибка разбора JSON или другая ошибка"
// @router /:id [put]
func (c *QuotaController) Update() {
	_ = c.Ctx.Input.CopyBody(1024)
	id, _ := c.GetInt(":id")
	var quota models.Quota

	// Проверка десериализации JSON
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &quota); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	quota.Id = id
	if err := models.UpdateQuota(&quota); err == nil {
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
// @Success 200 string "Удаление успешно выполнено"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
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

// AddSpecialityToQuota adds a speciality to a quota by their respective IDs.
// @Title AddSpecialityToQuota
// @Description Добавление специальности к квоте по их ID.
// @Param	quota_id	path	int	true	"ID квоты"
// @Param	speciality_id	path	int	true	"ID специальности"
// @Success 200 string	"Специальность успешно добавлена к квоте"
// @Failure 400 {string} string "400 ошибка разбора JSON или другая ошибка"
// @router /:quota_id/specialities/:speciality_id [post]
func (c *QuotaController) AddSpecialityToQuota() {
	quotaId, err1 := c.GetInt(":quota_id")
	specialityId, err2 := c.GetInt(":speciality_id")

	if err1 != nil || err2 != nil {
		c.Data["json"] = "quota_id and speciality_id must be integers"
		c.ServeJSON()
		return
	}

	err := models.AddSpecialityToQuota(specialityId, quotaId)
	if err == nil {
		c.Data["json"] = "Speciality added to quota successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetQuotaWithSpecialities retrieves information about a quota including its associated specialities by ID.
// @Title GetQuotaWithSpecialities
// @Description Получение информации о квоте вместе с ассоциированными специальностями по ID.
// @Param	id		path	int	true	"ID квоты для получения информации"
// @Success 200 {object} models.Quota "Информация о квоте вместе с специальностями"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /all/:id [get]
func (c *QuotaController) GetQuotaWithSpecialities() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.Data["json"] = "Invalid ID"
		c.ServeJSON()
		return
	}

	quota, err := models.GetQuotaWithSpecialitiesById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = quota
	}
	c.ServeJSON()
}
