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
// @router /assign_city/:universityId/:cityId [put]
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

// SearchUniversities ищет университеты по имени, абревиатуре или коду.
// @Title SearchUniversities
// @Description Поиск университетов по имени.
// @Param	name		query	string	true	"Имя университета для поиска"
// @Success 200 {array} models.University "Список найденных университетов"
// @Failure 400 {string} string "400 ошибка поиска или другая ошибка"
// @router /search [get]
func (c *UniversityController) SearchUniversitiesByName() {
	name := c.GetString("name")
	universities, err := models.SearchUniversitiesByName(name)
	if err == nil {
		c.Data["json"] = universities
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
// @Param	has_military_dept	query	bool	false	"Наличие военной кафедры"
// @Param	has_dormitory		query	bool	false	"Наличие общежития"
// @Param	city_id				query	int		false	"ID города"
// @Param	speciality_id		query	int		false	"Специальность"
// @Param   sort    			query   string  false  "Sort parameter (avg_fee_asc or avg_fee_desc)"
// @Success 200 {array} models.University "Список найденных университетов"
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

	sort := c.GetString("sort")
	if sort == "avg_fee_asc" || sort == "avg_fee_desc" {
		params["sort"] = sort
	}

	log.Printf("Received parameters map: %+v", params)

	universities, err := models.SearchUniversities(params)
	if err == nil {
		c.Data["json"] = universities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
