package controllers

import (
	"encoding/json"
	"net/http"
	"testhub-spec-uni/models"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// QuotaController handles operations for Quota.
type QuotaController struct {
	beego.Controller
}

type QuotaResponse struct {
	Id           int    `json:"id"`
	QuotaType    string `json:"quotaType"`
	Count        int    `json:"count"`
	MinScore     int    `json:"minScore"`
	MaxScore     int    `json:"maxScore"`
	Specialities []int  `json:"specialities"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func ConvertTimeToString(t time.Time) string {
	return t.Format(time.RFC3339) // или другой формат, который вам нужен
}

// ConvertSpecialities converts []*models.Speciality to []int.
func ConvertSpecialities(specialities []*models.Speciality) []int {
	result := make([]int, len(specialities))
	for i, s := range specialities {
		result[i] = s.Id
	}
	return result
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
// @Param	lang	header	string	true	"Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} QuotaResponse	"Информация о квоте"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:id [get]
// Get returns information about a quota by its ID.
// @Title Get
// @Description Получение информации о квоте по ID.
// @Param	id		path	int	true	"ID квоты для получения информации"
// @Param	lang	header	string	true	"Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} QuotaResponse	"Информация о квоте"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:id [get]
func (c *QuotaController) Get() {
	id, _ := c.GetInt(":id")
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	quota, err := models.GetQuotaById(id, language)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		response := QuotaResponse{
			Id:           quota.Id,
			QuotaType:    quota.QuotaType,
			Count:        quota.Count,
			MinScore:     quota.MinScore,
			MaxScore:     quota.MaxScore,
			Specialities: ConvertSpecialities(quota.Specialities),
			CreatedAt:    ConvertTimeToString(quota.CreatedAt),
			UpdatedAt:    ConvertTimeToString(quota.UpdatedAt),
		}
		c.Data["json"] = response
	}
	c.ServeJSON()
}

// GetAll returns a list of all quotas.
// @Title GetAll
// @Description Получение списка всех квот.
// @Param	lang	header	string	true	"Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} QuotaResponse	"Список квот"
// @Failure 400 {string} string "400 ошибка получения списка или другая ошибка"
// @router / [get]
func (c *QuotaController) GetAll() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	quotas, err := models.GetAllQuotas(language)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		response := make([]QuotaResponse, len(quotas))
		for i, quota := range quotas {
			response[i] = QuotaResponse{
				Id:           quota.Id,
				QuotaType:    quota.QuotaType,
				Count:        quota.Count,
				MinScore:     quota.MinScore,
				MaxScore:     quota.MaxScore,
				Specialities: ConvertSpecialities(quota.Specialities),
				CreatedAt:    ConvertTimeToString(quota.CreatedAt),
				UpdatedAt:    ConvertTimeToString(quota.UpdatedAt),
			}
		}
		c.Data["json"] = response
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

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &quota); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	quota.Id = id

	var updatedFields map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &updatedFields); err != nil {
		c.Data["json"] = "Invalid input: " + err.Error()
		c.ServeJSON()
		return
	}

	fields := make([]string, 0, len(updatedFields))
	for field := range updatedFields {
		fields = append(fields, field)
	}

	if err := models.UpdateQuota(&quota, fields...); err == nil {
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

// AddSpecialityToQuota adds a speciality to a quota by their IDs.
// @Title AddSpecialityToQuota
// @Description Добавление специальности к квоте по их ID.
// @Param	quotaId		path	int	true	"ID квоты"
// @Param	specialityId	path	int	true	"ID специальности"
// @Success 200 string "Добавление успешно выполнено"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:quotaId/specialities/:specialityId [post]
func (c *QuotaController) AddSpecialityToQuota() {
	quotaId, _ := c.GetInt(":quotaId")
	specialityId, _ := c.GetInt(":specialityId")

	err := models.AddSpecialityToQuota(specialityId, quotaId)
	if err == nil {
		c.Data["json"] = "Add speciality to quota successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetQuotaWithSpecialities returns a list of all quotas with their associated specialities.
// @Title GetAllWithSpecialities
// @Description Получение списка всех квот со специальностями.
// @Param	lang	header	string	true	"Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} QuotaResponse	"Список квот со специальностями"
// @Failure 400 {string} string "400 ошибка получения списка или другая ошибка"
// @router all/:id [get]
func (c *QuotaController) GetQuotaWithSpecialities() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	quotas, err := models.GetAllQuotasWithSpecialities(language)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		response := make([]QuotaResponse, len(quotas))
		for i, quota := range quotas {
			response[i] = QuotaResponse{
				Id:           quota.Id,
				QuotaType:    quota.QuotaType,
				Count:        quota.Count,
				MinScore:     quota.MinScore,
				MaxScore:     quota.MaxScore,
				Specialities: ConvertSpecialities(quota.Specialities),
				CreatedAt:    ConvertTimeToString(quota.CreatedAt),
				UpdatedAt:    ConvertTimeToString(quota.UpdatedAt),
			}
		}
		c.Data["json"] = response
	}
	c.ServeJSON()
}

// GetWithSpecialitiesById returns information about a quota with its associated specialities by its ID.
// @Title GetWithSpecialitiesById
// @Description Получение информации о квоте со специальностями по ID.
// @Param	id		path	int	true	"ID квоты для получения информации"
// @Param	lang	header	string	true	"Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} QuotaResponse	"Информация о квоте со специальностями"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /with-specialities/:id [get]
func (c *QuotaController) GetWithSpecialitiesById() {
	id, _ := c.GetInt(":id")
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	quota, err := models.GetQuotaWithSpecialitiesById(id, language)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		response := QuotaResponse{
			Id:           quota.Id,
			QuotaType:    quota.QuotaType,
			Count:        quota.Count,
			MinScore:     quota.MinScore,
			MaxScore:     quota.MaxScore,
			Specialities: ConvertSpecialities(quota.Specialities),
			CreatedAt:    ConvertTimeToString(quota.CreatedAt),
			UpdatedAt:    ConvertTimeToString(quota.UpdatedAt),
		}
		c.Data["json"] = response
	}
	c.ServeJSON()
}
