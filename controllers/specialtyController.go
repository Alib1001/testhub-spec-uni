package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// SpecialityController обрабатывает запросы для работы со специальностями.
type SpecialityController struct {
	beego.Controller
}

// SearchSpecialitiesByName searches for specialities by name prefix using Elasticsearch.
// @Title SearchSpecialitiesByName
// @Description Search for specialities by name prefix.
// @Param	prefix	query	string	true	"Prefix of the speciality name to search for"
// @Success 200 {array} models.Speciality	"List of specialities"
// @Failure 400 error searching or other error
// @router /search_by_name [get]
func (c *SpecialityController) SearchSpecialitiesByName() {
	prefix := c.GetString("prefix")
	if specialities, err := models.SearchSpecialitiesByName(prefix); err == nil {
		c.Data["json"] = specialities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
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
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &speciality); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	if id, err := models.AddSpeciality(&speciality); err == nil {
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
	if speciality, err := models.GetSpecialityById(id); err == nil {
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
	if specialities, err := models.GetAllSpecialities(); err == nil {
		c.Data["json"] = specialities
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
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &speciality); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	speciality.Id = id
	if err := models.UpdateSpeciality(&speciality); err == nil {
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
	if err := models.DeleteSpeciality(id); err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetSpecialitiesInUni retrieves all specialities associated with a university by its ID.
// @Title GetSpecialitiesInUni
// @Description Получение списка специальностей, связанных с университетом.
// @Param	universityId		path	int	true	"ID университета"
// @Success 200 {array} models.Speciality	"Список специальностей университета"
// @Failure 400 некорректный ID или другая ошибка
// @router /byuni/:universityId [get]
func (c *SpecialityController) GetByUniversity() {
	universityId, err := c.GetInt(":universityId")
	if err != nil {
		c.CustomAbort(400, "Invalid university ID")
		return
	}

	specialities, err := models.GetSpecialitiesInUniversity(universityId)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	if len(specialities) == 0 {
		c.Data["json"] = []models.Speciality{}
	} else {
		c.Data["json"] = specialities
	}
	c.ServeJSON()
}

// @Title AssociateSpecialityWithSubjectPair
// @Description связывает специальность с парой предметов
// @Param	speciality_id		path 	int	true		"ID специальности"
// @Param	subject_pair_id		path 	int	true		"ID пары предметов"
// @Success 200 {object} models.SubjectPair
// @Failure 400 "Invalid input"
// @Failure 404 "SubjectPair not found"
// @router /associatePair/:speciality_id/:subject_pair_id [put]
func (c *SpecialityController) AssociateSpecialityWithSubjectPair() {
	specialityIdStr := c.Ctx.Input.Param(":speciality_id")
	subjectPairIdStr := c.Ctx.Input.Param(":subject_pair_id")

	specialityId, err := strconv.Atoi(specialityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid speciality_id")
	}

	subjectPairId, err := strconv.Atoi(subjectPairIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid subject_pair_id")
	}

	err = models.AssociateSpecialityWithSubjectPair(specialityId, subjectPairId)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]interface{}{"message": "Speciality associated with SubjectPair successfully"}
	c.ServeJSON()
}

// @Title GetSubjectPairsBySpecialityId
// @Description получает все пары предметов для заданной специальности
// @Param	speciality_id		path 	int	true		"ID специальности"
// @Success 200 {array} models.SubjectPair
// @Failure 400 "Invalid input"
// @Failure 404 "SubjectPairs not found"
// @router /byspec/:speciality_id [get]
func (c *SpecialityController) GetSubjectPairsBySpecialityId() {
	specialityIdStr := c.Ctx.Input.Param(":speciality_id")
	specialityId, err := strconv.Atoi(specialityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid speciality_id")
	}

	subjectPairs, err := models.GetSubjectPairsBySpecialityId(specialityId)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = subjectPairs
	c.ServeJSON()
}

// @Title GetSpecialityBySubjectPair
// @Description получает специальность по ID первого и второго предметов
// @Param	subject1_id		path 	int	true		"ID первого предмета"
// @Param	subject2_id		path 	int	true		"ID второго предмета"
// @Success 200 {object} models.Speciality
// @Failure 400 "Invalid input"
// @Failure 404 "Speciality not found"
// @router /bysubjects/:subject1_id/:subject2_id [get]
func (c *SpecialityController) GetSpecialitiesBySubjectPair() {
	subject1IdStr := c.Ctx.Input.Param(":subject1_id")
	subject2IdStr := c.Ctx.Input.Param(":subject2_id")

	subject1Id, err := strconv.Atoi(subject1IdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid subject1_id")
	}
	subject2Id, err := strconv.Atoi(subject2IdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid subject2_id")
	}

	speciality, err := models.GetSpecialitiesBySubjectPair(subject1Id, subject2Id)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	if speciality == nil {
		c.CustomAbort(http.StatusNotFound, "Speciality not found")
	}

	c.Data["json"] = speciality
	c.ServeJSON()
}

// AddPointStat добавляет статистику по баллам для специальности и университета.
// @Title AddPointStat
// @Description Добавление статистики по баллам для специальности и университета.
// @Param	universityId		path	int	true	 "ID университета"
// @Param	specialityId		path	int	true	"ID специальности"
// @Param	body	body	models.PointStat	true	"JSON с данными о статистике по баллам"
// @Success 200 {object} map[string]int64	"ID добавленной статистики"
// @Failure 400 ошибка разбора JSON или другая ошибка
// @router /addPointStat/:universityId/:specialityId [post]
func (c *SpecialityController) AddPointStat() {
	_ = c.Ctx.Input.CopyBody(1024)

	universityId, err := c.GetInt(":universityId")
	if err != nil {
		c.CustomAbort(400, "Invalid university ID")
		return
	}

	specialityId, err := c.GetInt(":specialityId")
	if err != nil {
		c.CustomAbort(400, "Invalid speciality ID")
		return
	}

	var pointStat models.PointStat
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &pointStat); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Устанавливаем ID университета и специальности из параметров запроса
	pointStat.University = &models.University{Id: universityId}
	pointStat.Speciality = &models.Speciality{Id: specialityId}

	id, err := models.AddPointStat(universityId, specialityId, &pointStat)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]int64{"id": id}
	}
	c.ServeJSON()
}

// GetPointStatsByUniversityAndSpeciality возвращает статистику по баллам для специальности и университета.
// @Title GetPointStatsByUniversityAndSpeciality
// @Description Получение статистики по баллам для специальности и университета.
// @Param	universityId		path	int	true	"ID университета"
// @Param	specialityId		path	int	true	"ID специальности"
// @Success 200 {array} models.PointStat	"Список статистики по баллам"
// @Failure 400 некорректный ID или другая ошибка
// @router /pointstatsbyparams/:universityId/:specialityId [get]
func (c *SpecialityController) GetPointStatsByUniversityAndSpeciality() {
	universityId, err := c.GetInt(":universityId")
	if err != nil {
		c.CustomAbort(400, "Invalid university ID")
		return
	}

	specialityId, err := c.GetInt(":specialityId")
	if err != nil {
		c.CustomAbort(400, "Invalid speciality ID")
		return
	}

	pointStats, err := models.GetPointStatsByUniversityAndSpeciality(universityId, specialityId)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	if len(pointStats) == 0 {
		c.Data["json"] = []models.PointStat{} // Возвращаем пустой массив, если статистика не найдена
	} else {
		c.Data["json"] = pointStats
	}
	c.ServeJSON()
}
