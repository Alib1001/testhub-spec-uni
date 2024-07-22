package controllers

import (
	"encoding/json"
	"log"
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
// @Failure 400 оsaшибка разбора JSON или другая ошибка
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

// AssignCityToUniversity назначает город университету по их ID.
// @Title AssignCityToUniversity
// @Description Назначение города университету.
// @Param	universityId		path	int	true	"ID университета"
// @Param	cityId		        path	int	true	"ID города"
// @Success 200 string	"Город успешно назначен"
// @Failure 400 некорректные ID или другая ошибка
// @router /assigncity/:universityId/:cityId [put]
func (c *UniversityController) AssignCityToUniversity() {
	universityId, _ := c.GetInt(":universityId")
	cityId, _ := c.GetInt(":cityId")
	err := models.AssignCityToUniversity(universityId, cityId)
	if err == nil {
		c.Data["json"] = "City successfully assigned"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// AddSpecialityToUniversity adds a speciality to a university by their IDs.
// @Title AddSpecialityToUniversity
// @Description Добавление специальности к университету.
// @Param	universityId		path	int	true	"ID университета"
// @Param	specialityId		path	int	true	"ID специальности"
// @Success 200 string	"Специальность успешно добавлена к университету"
// @Failure 400 некорректные ID или другая ошибка
// @router /assign_speciality/:universityId/:specialityId [post]
func (c *UniversityController) AddSpecialityToUniversity() {
	universityId, _ := c.GetInt(":universityId")
	specialityId, _ := c.GetInt(":specialityId")

	err := models.AddSpecialityToUniversity(specialityId, universityId)
	if err == nil {
		c.Data["json"] = "Speciality added to university successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// AddSpecialitiesToUniversity добавляет несколько специальностей к университету по их ID.Передавать в тело с обычный массив [x,y]
// @Title AddSpecialitiesToUniversity
// @Description Добавление нескольких специальностей к университету.
// @Param	universityId		path	int				true	"ID университета"
// @Param	body				body	[]int			true	"Массив ID специальностей"
// @Success 200 string	"Специальности успешно добавлены к университету"
// @Failure 400 некорректные ID или другая ошибка
// @router /assign_specialities/:universityId [post]
func (c *UniversityController) AddSpecialitiesToUniversity() {
	universityId, _ := c.GetInt(":universityId")
	_ = c.Ctx.Input.CopyBody(512)
	var specialityIds []int
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &specialityIds); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	err := models.AddSpecialitiesToUniversity(specialityIds, universityId)
	if err == nil {
		c.Data["json"] = "Specialities added to university successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// AddServicesToUniversity добавляет несколько сервисов к университету по их ID.Передавать в тело с обычный массив [x,y]
// @Title AddServicesToUniversity
// @Description Добавление нескольких сервисов к университету.
// @Param	universityId		path	int				true	"ID университета"
// @Param	body				body	[]int			true	"Массив ID сервисов"
// @Success 200 string	"Сервисы успешно добавлены к университету"
// @Failure 400 некорректные ID или другая ошибка
// @router /assign_services/:universityId [post]
func (c *UniversityController) AddServicesToUniversity() {
	universityId, _ := c.GetInt(":universityId")
	var serviceIds []int
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &serviceIds); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	err := models.AddServicesToUniversity(serviceIds, universityId)
	if err == nil {
		c.Data["json"] = "Services added to university successfully"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// SearchUniversities ищет университеты по различным параметрам.
// @Title SearchUniversities
// @Description Поиск университетов по параметрам.
// @Param	min_score			query	int		false	"Минимальный балл"
// @Param	avg_fee				query	int		false	"Средняя цена"
// @Param	city_id				query	int		false	"ID города"
// @Param	speciality_ids		query	string	false	"Список специальностей в JSON формате, должны передавать массив с id специальнотей"
// @Param	service_ids			query	string	false	"Список сервисов в JSON формате, должны передавать массив с id сервисов"
// @Param	first_subject_id	query	int		false	"ID первого предмета"
// @Param	second_subject_id	query	int		false	"ID второго предмета"
// @Param	sort    			query   string  false  "Sort parameter (avg_fee_asc or avg_fee_desc)"
// @Param  name                query   string  false  "Название университета или его часть"
// @Param  study_format        query   string  false  "Формат обучения (full_time, part_time, etc.)"
// @Param  page                query   int     false  "Номер страницы"
// @Param  per_page            query   int     false  "Количество элементов на одной странице"
// @Success 200 {object} models.SearchResult "Список найденных университетов с информацией о пагинации"
// @Failure 400 {string} string "400 ошибка поиска или другая ошибка"
// @router /search [get]
func (c *UniversityController) SearchUniversities() {
	params := make(map[string]interface{})

	if minScore, err := c.GetInt("min_score"); err == nil {
		params["min_score"] = minScore
	}
	if avgFee, err := c.GetInt("avg_fee"); err == nil {
		params["avg_fee"] = avgFee
	}
	if hasMilitaryDept, err := c.GetBool("has_military_dept"); err == nil {
		params["has_military_dept"] = hasMilitaryDept
	}
	if hasDormitory, err := c.GetBool("has_dormitory"); err == nil {
		params["has_dormitory"] = hasDormitory
	}
	if cityID, err := c.GetInt("city_id"); err == nil {
		params["city_id"] = cityID
	}
	if specialityIDsStr := c.GetString("speciality_ids"); specialityIDsStr != "" {
		var specialityIDs []int
		err := json.Unmarshal([]byte(specialityIDsStr), &specialityIDs)
		if err == nil {
			params["speciality_ids"] = specialityIDs
		}
	}
	if serviceIDsStr := c.GetString("service_ids"); serviceIDsStr != "" {
		var serviceIDs []int
		err := json.Unmarshal([]byte(serviceIDsStr), &serviceIDs)
		if err == nil {
			params["service_ids"] = serviceIDs
		} else {
			log.Printf("Error unmarshaling service_ids: %v", err)
		}
	}
	if firstSubjectID, err := c.GetInt("first_subject_id"); err == nil {
		params["first_subject_id"] = firstSubjectID
	}
	if secondSubjectID, err := c.GetInt("second_subject_id"); err == nil {
		params["second_subject_id"] = secondSubjectID
	}
	if sort := c.GetString("sort"); sort == "avg_fee_asc" || sort == "avg_fee_desc" {
		params["sort"] = sort
	}
	if name := c.GetString("name"); name != "" {
		params["name"] = name
	}
	if studyFormat := c.GetString("study_format"); studyFormat != "" {
		params["study_format"] = studyFormat
	}
	if page, err := c.GetInt("page"); err == nil {
		params["page"] = page
	}
	if perPage, err := c.GetInt("per_page"); err == nil {
		params["per_page"] = perPage
	}

	log.Printf("Received parameters map: %+v", params)

	result, err := models.SearchUniversities(params)
	if err == nil {
		c.Data["json"] = result
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
