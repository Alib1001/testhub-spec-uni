package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// SpecialityController обрабатывает запросы для работы со специальностями.
type SpecialityController struct {
	beego.Controller
}

// Create adds a new speciality to the database.
// @Title Create
// @Description Create a new speciality.
// @Param   NameRu         formData string true "Name of the speciality in Russian"
// @Param   NameKz         formData string true "Name of the speciality in Kazakh"
// @Param   AbbreviationRu formData string true "Abbreviation of the speciality in Russian"
// @Param   AbbreviationKz formData string true "Abbreviation of the speciality in Kazakh"
// @Param   Subject_1      formData int    true "ID of the first subject"
// @Param   Subject_2      formData int    true "ID of the second subject"
// @Param   Degree         formData string true "Degree of the speciality"
// @Param   Code           formData string true "Code of the speciality"
// @Param   Term           formData int    true "Term of the speciality"
// @Param   DescriptionRu  formData string false "Description of the speciality in Russian"
// @Param   DescriptionKz  formData string false "Description of the speciality in Kazakh"
// @Param   Scholarship    formData bool   false "Whether the speciality has a scholarship"
// @Success 200 {object} map[string]int64 "ID of the created speciality"
// @Failure 400 {string} string "JSON parsing error or other error"
// @router / [post]
func (c *SpecialityController) Create() {
	var data models.AddSpecialityResponse

	if err := c.ParseForm(&data); err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Validate form data
	validate := validator.New()
	if err := validate.Struct(&data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMap := make(map[string]string)
		for _, err := range validationErrors {
			errMap[err.Field()] = fmt.Sprintf("Validation failed on the '%s' tag", err.Tag())
		}
		c.Data["json"] = errMap
		c.ServeJSON()
		return
	}

	// Create speciality
	if id, err := models.AddSpecialityFromFormData(&data); err == nil {
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
	lang := c.Ctx.Input.Header("lang")

	if speciality, err := models.GetSpecialityById(id, lang); err == nil {
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
	lang := c.Ctx.Input.Header("lang")

	if specialities, err := models.GetAllSpecialities(lang); err == nil {
		c.Data["json"] = specialities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update updates the information of a speciality by its ID.
// @Title Update Speciality
// @Description Update the information of a speciality by its ID.
// @Param   id             path     int    true  "ID of the speciality to update"
// @Param   NameRu         formData string false "Updated name of the speciality in Russian"
// @Param   NameKz         formData string false "Updated name of the speciality in Kazakh"
// @Param   AbbreviationRu formData string false "Updated abbreviation of the speciality in Russian"
// @Param   AbbreviationKz formData string false "Updated abbreviation of the speciality in Kazakh"
// @Param   Degree         formData string false "Updated degree of the speciality"
// @Param   Code           formData string false "Updated code of the speciality"
// @Param   Term           formData int    false "Updated term of the speciality"
// @Param   DescriptionRu  formData string false "Updated description of the speciality in Russian"
// @Param   DescriptionKz  formData string false "Updated description of the speciality in Kazakh"
// @Param   Scholarship    formData bool   false "Updated scholarship status of the speciality"
// @Param   Subject_1      formData int    false "Updated ID of the first subject"
// @Param   Subject_2      formData int    false "Updated ID of the second subject"
// @Success 200 {string} string "Update successful"
// @Failure 400 {string} string "Invalid input or other error"
// @router /:id [put]
func (c *SpecialityController) Update() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid speciality ID")
		return
	}

	// Log raw form data
	log.Printf("Raw form data: %+v\n", c.Ctx.Input.RequestBody)

	var data models.UpdateSpecialityResponse
	if err := c.ParseForm(&data); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	data.Id = id

	// Log the parsed form data
	log.Printf("Parsed form data: %+v\n", data)

	if err := models.UpdateSpecialityFromFormData(&data); err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = "Update successful"
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

	lang := c.Ctx.Input.Header("lang")
	if lang == "" {
		lang = "ru"
	}

	specialities, err := models.GetSpecialitiesInUniversityForUser(universityId, lang)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	if len(specialities) == 0 {
		c.Data["json"] = []models.GetSpecialityForAdmResponse{}
	} else {
		c.Data["json"] = specialities
	}
	c.ServeJSON()
}

// GetByUniversityForAdmin retrieves all specialities associated with a university by its ID.
// @Title GetSpecialitiesInUni
// @Description Получение списка специальностей, связанных с университетом.
// @Param	universityId		path	int	true	"ID университета"
// @Success 200 {array} models.GetByUniResponseForAdm	"Список специальностей университета"
// @Failure 400 некорректный ID или другая ошибка
// @router /byuni/:universityId [get]
func (c *SpecialityController) GetByUniversityForAdmin() {
	universityId, err := c.GetInt(":universityId")
	if err != nil {
		c.CustomAbort(400, "Invalid university ID")
		return
	}

	specialities, err := models.GetSpecialitiesInUniversityForAdmin(universityId)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = specialities
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

// GetSubjectPairsBySpecialityId получает все пары предметов для заданной специальности.
// @Title GetSubjectPairsBySpecialityId
// @Description Получение всех пар предметов для заданной специальности.
// @Param	speciality_id		path 	int	true		"ID специальности"
// @Param	lang	header	string	false	"Язык для фильтрации"
// @Success 200 {array} models.SubjectPair	"Список пар предметов"
// @Failure 400 "Некорректный ввод"
// @Failure 404 "Пары предметов не найдены"
// @router /byspec/:speciality_id [get]
func (c *SpecialityController) GetSubjectPairsBySpecialityId() {
	specialityIdStr := c.Ctx.Input.Param(":speciality_id")
	specialityId, err := strconv.Atoi(specialityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid speciality_id")
	}

	lang := c.Ctx.Input.Header("lang")

	subjectPairs, err := models.GetSubjectPairsBySpecialityId(specialityId, lang)
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
	lang := c.Ctx.Input.Header("lang")

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

	speciality, err := models.GetSpecialitiesBySubjectPair(subject1Id, subject2Id, lang)
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

	var form models.AddPointStatResponse
	if err := c.ParseForm(&form); err != nil {
		c.CustomAbort(400, "Invalid form data")
		return
	}

	pointStat := models.PointStat{
		GrantCount:    form.AnnualGrants,
		MinScore:      form.MinScore,
		MinGrantScore: form.MinGrantScore,
		Year:          form.Year,
		AvgSalary:     form.AvgSalary,
		Price:         form.Price,
		University:    &models.University{Id: universityId},
		Speciality:    &models.Speciality{Id: specialityId},
	}

	fmt.Printf("PointStat before insert: %+v\n", pointStat)

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
		c.Data["json"] = []models.PointStat{}
	} else {
		c.Data["json"] = pointStats
	}
	c.ServeJSON()
}

// SearchSpecialities выполняет поиск с различными фильтрами и пагинацией.
// @Title SearchSpecialities
// @Description Выполнение поиска специальностей с использованием фильтров и пагинации.
// @Param	name	query	string	false	"Префикс имени специальности для поиска"
// @Param	subject1_id	query	int	false	"ID первого предмета для фильтрации"
// @Param	subject2_id	query	int	false	"ID второго предмета для фильтрации"
// @Param	university_id	query	int	false	"ID университета для фильтрации"
// @Param	page	query	int	false	"Номер страницы для пагинации"
// @Param	per_page	query	int	false	"Количество элементов на странице"
// @Param	lang	header	string	false	"Язык для фильтрации"
// @Success 200 {object} models.SpecialitySearchResult	"Результаты поиска со специальностями"
// @Failure 400 "Ошибка поиска или другая ошибка"
// @router /search [get]
func (c *SpecialityController) SearchSpecialities() {
	params := make(map[string]interface{})

	if name := c.GetString("name"); name != "" {
		params["name"] = name
	}

	if subject1Id, err := c.GetInt("subject1_id"); err == nil {
		params["subject1_id"] = subject1Id
	}

	if subject2Id, err := c.GetInt("subject2_id"); err == nil {
		params["subject2_id"] = subject2Id
	}

	if universityId, err := c.GetInt("university_id"); err == nil {
		params["university_id"] = universityId
	}

	if page, err := c.GetInt("page"); err == nil {
		params["page"] = page
	}

	if perPage, err := c.GetInt("per_page"); err == nil {
		params["per_page"] = perPage
	}

	lang := c.Ctx.Input.Header("lang")
	if lang != "" {
		params["lang"] = lang
	}

	result, err := models.SearchSpecialities(params, lang)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = result
	c.ServeJSON()
}
