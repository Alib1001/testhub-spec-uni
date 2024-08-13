package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testhub-spec-uni/models"
)

// UniversityController обрабатывает запросы для работы с университетами.
type UniversityController struct {
	beego.Controller
}

// Create
// @Title CreateUniversity
// @Description Создает новый университет с привязкой к сервисам и галерее изображений
// @Param   NameRu            formData    string  true  "Название на русском языке"
// @Param   NameKz            formData    string  true  "Название на казахском языке"
// @Param   UniversityStatusRu formData   string  true  "Статус университета на русском языке"
// @Param   UniversityStatusKz formData   string  true  "Статус университета на казахском языке"
// @Param   Website           formData    string  true  "Вебсайт"
// @Param   CallCenterNumber  formData    string  true  "Номер колл-центра"
// @Param   WhatsAppNumber    formData    string  true  "Номер WhatsApp"
// @Param   Address           formData    string  true  "Адрес"
// @Param   UniversityCode    formData    string  true  "Код университета"
// @Param   StudyFormatRu     formData    string  true  "Формат обучения на русском языке"
// @Param   StudyFormatKz     formData    string  true  "Формат обучения на казахском языке"
// @Param   AbbreviationRu    formData    string  true  "Аббревиатура на русском языке"
// @Param   AbbreviationKz    formData    string  true  "Аббревиатура на казахском языке"
// @Param   MainImageUrl      formData    file    true  "Главное изображение"
// @Param   AddressLink       formData    string  true  "Ссылка на адрес"
// @Param   DescriptionRu     formData    string  true  "Описание на русском языке"
// @Param   DescriptionKz     formData    string  true  "Описание на казахском языке"
// @Param   Rating            formData    string  true  "Рейтинг"
// @Param   Gallery           formData    file    true  "Галерея изображений"
// @Param   CityId            formData    int     true  "ID города"
// @Param   ServiceIds        formData    string  true  "Список ID сервисов в формате [1,2,3]"
// @Param   MinScore           formData    int     true  "Минимальный балл для университета"
// @Success 200 {int64} id "ID созданного университета"
// @Failure 400 {object} map[string]string "Error message"
// @Failure 500 {object} map[string]string "Error message"
// @router /universities [post]
func (c *UniversityController) Create() {
	err := c.Ctx.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse form data: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	var partialResponse models.AddUniversityPartial

	if err := c.ParseForm(&partialResponse); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse form data: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	var universityResponse models.AddUUniversityResponse
	serviceIdsStr := c.GetString("ServiceIds")
	if serviceIdsStr != "" {
		serviceIdsStr = strings.Trim(serviceIdsStr, "[]")
		serviceIds := strings.Split(serviceIdsStr, ",")
		for _, serviceID := range serviceIds {
			id, err := strconv.Atoi(strings.TrimSpace(serviceID))
			if err != nil {
				c.Data["json"] = map[string]string{"error": "Invalid service ID: " + serviceID}
				c.Ctx.Output.SetStatus(400)
				c.ServeJSON()
				return
			}
			universityResponse.ServiceIds = append(universityResponse.ServiceIds, id)
		}
	}

	universityResponse.NameRu = partialResponse.NameRu
	universityResponse.NameKz = partialResponse.NameKz
	universityResponse.UniversityStatusRu = partialResponse.UniversityStatusRu
	universityResponse.UniversityStatusKz = partialResponse.UniversityStatusKz
	universityResponse.Website = partialResponse.Website
	universityResponse.Email = partialResponse.Email
	universityResponse.CallCenterNumber = partialResponse.CallCenterNumber
	universityResponse.WhatsAppNumber = partialResponse.WhatsAppNumber
	universityResponse.Address = partialResponse.Address
	universityResponse.UniversityCode = partialResponse.UniversityCode
	universityResponse.StudyFormatRu = partialResponse.StudyFormatRu
	universityResponse.StudyFormatKz = partialResponse.StudyFormatKz
	universityResponse.AbbreviationRu = partialResponse.AbbreviationRu
	universityResponse.AbbreviationKz = partialResponse.AbbreviationKz
	universityResponse.MainImageUrl = partialResponse.MainImageUrl
	universityResponse.AddressLink = partialResponse.AddressLink
	universityResponse.DescriptionRu = partialResponse.DescriptionRu
	universityResponse.DescriptionKz = partialResponse.DescriptionKz
	universityResponse.Rating = partialResponse.Rating
	universityResponse.MinScore = partialResponse.MinScore
	universityResponse.Gallery = partialResponse.Gallery
	universityResponse.CityId = partialResponse.CityId

	id, err := models.AddUniversity(&universityResponse)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, validationErr := range validationErrors {
				errors[validationErr.Field()] = "Field validation error: " + validationErr.Tag()
			}
			c.Data["json"] = errors
		} else {
			c.Data["json"] = map[string]string{"error": "Validation failed: " + err.Error()}
		}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	file, header, err := c.GetFile("MainImageUrl")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to get main image file: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}
	defer file.Close()

	universityID := strconv.FormatInt(id, 10)
	mainImagePath := fmt.Sprintf("Universities/%s/%s", universityID, header.Filename)

	uploadedMainImageURL, err := models.UploadFileToCloud(mainImagePath, file)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to upload main image: " + err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	err = models.UpdateUniversityImageURL(id, uploadedMainImageURL)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to update university main image URL: " + err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	galleryFiles := c.Ctx.Request.MultipartForm.File["Gallery"]
	var galleryURLs []string

	for _, galleryFileHeader := range galleryFiles {
		galleryFile, err := galleryFileHeader.Open()
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Failed to open gallery file: " + err.Error()}
			c.Ctx.Output.SetStatus(400)
			c.ServeJSON()
			return
		}
		defer galleryFile.Close()

		galleryFilePath := fmt.Sprintf("Universities/%s/Gallery/%s", universityID, galleryFileHeader.Filename)
		uploadedGalleryURL, err := models.UploadFileToCloud(galleryFilePath, galleryFile)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Failed to upload gallery image: " + err.Error()}
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}

		galleryURLs = append(galleryURLs, uploadedGalleryURL)
	}

	err = models.AddGalleryImages(id, galleryURLs)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to add gallery images: " + err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]int64{"id": id}
	c.Ctx.Output.SetStatus(200)
	c.ServeJSON()
}

// GetForAdmin возвращает информацию о университете по его ID.
// @Title GetForAdmin
// @Description Получение информации о университете по ID.
// @Param	id		path	int	true	"ID университета для получения информации"
// @Success 200 {object} models.University	"Информация о университете"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
func (c *UniversityController) GetForAdmin() {
	id, _ := c.GetInt(":id")
	university, err := models.GetUniversityByIdForAdmin(id)
	if err == nil {
		c.Data["json"] = university
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetForUser возвращает информацию о университете по его ID.
// @Title GetForUser
// @Description Получение информации о университете по ID.
// @Param	id		path	int	true	"ID университета для получения информации"
// @Param   language    query   string  false   "Язык для получения информации (ru/kz)"
// @Success 200 {object} models.University	"Информация о университете"
// @Failure 400 некорректный ID или другая ошибка
// @router /:id [get]
func (c *UniversityController) GetForUser() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.Data["json"] = "Некорректный ID университета"
		c.ServeJSON()
		return
	}

	language := c.Ctx.Input.Header("lang")

	university, err := models.GetUniversityByIdForUser(id, language)
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
// @Param page query int true "Page number"
// @Param per_page query int true "Items per page"
// @Success 200 {object} map[string]interface{}	"Список университетов с пагинацией"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *UniversityController) GetAll() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	page, err := c.GetInt("page", 1)
	if err != nil || page < 1 {
		c.CustomAbort(http.StatusBadRequest, "Invalid page number")
		return
	}

	perPage, err := c.GetInt("per_page", 10)
	if err != nil || perPage < 1 {
		c.CustomAbort(http.StatusBadRequest, "Invalid per_page value")
		return
	}

	universities, totalCount, totalPage, currentPage, err := models.GetAllUniversities(language, page, perPage)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"universities": universities,
			"total_count":  totalCount,
			"total_page":   totalPage,
			"current_page": currentPage,
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetAllForAdmin возвращает список всех университетов.
// @Title GetAllForAdmin
// @Description Получение списка всех университетов.
// @Success 200 {array} models.University	"Список университетов"
// @Failure 400 ошибка получения списка или другая ошибка
// @router / [get]
func (c *UniversityController) GetAllForAdmin() {
	universities, err := models.GetAllUniversitiesForAdmin()
	if err == nil {
		c.Data["json"] = universities
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update @Title Update
// @Description Update the specified fields of a university
// @Param NameRu formData string false "Russian name of the university"
// @Param NameKz formData string false "Kazakh name of the university"
// @Param UniversityStatusRu formData string false "Russian university status"
// @Param UniversityStatusKz formData string false "Kazakh university status"
// @Param Website formData string false "Website of the university"
// @Param CallCenterNumber formData string false "Call center number"
// @Param WhatsAppNumber formData string false "WhatsApp number"
// @Param Address formData string false "Address of the university"
// @Param UniversityCode formData string false "University code"
// @Param StudyFormatRu formData string false "Russian study format"
// @Param StudyFormatKz formData string false "Kazakh study format"
// @Param AbbreviationRu formData string false "Russian abbreviation"
// @Param AbbreviationKz formData string false "Kazakh abbreviation"
// @Param AddressLink formData string false "Address link"
// @Param DescriptionRu formData string false "Russian description"
// @Param DescriptionKz formData string false "Kazakh description"
// @Param Rating formData int false "Rating of the university"
// @Param CityId formData int false "City ID"
// @Param MainImageUrl formData file false "Main image of the university"
// @Param Gallery formData file false "Gallery images of the university"
// @Param ServiceIds formData string false "Comma-separated list of service IDs"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @router /universities/{id} [put]
func (c *UniversityController) Update() {
	universityId, uniIderr := c.GetInt(":id")
	if uniIderr != nil {
		c.CustomAbort(400, "Invalid university ID")
		return
	}

	err := c.Ctx.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse form data: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	var partialResponse models.UpdateUniversityPartial
	partialResponse.Id = universityId

	if err := c.ParseForm(&partialResponse); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse form data: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	if partialResponse.Id == 0 {
		c.Data["json"] = map[string]string{"error": "Invalid university ID"}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	serviceIdsStr := c.GetString("ServiceIds")
	var serviceIds []int
	if serviceIdsStr != "" {
		serviceIdsStr = strings.Trim(serviceIdsStr, "[]")
		serviceIdsStrs := strings.Split(serviceIdsStr, ",")
		for _, serviceIDStr := range serviceIdsStrs {
			id, err := strconv.Atoi(strings.TrimSpace(serviceIDStr))
			if err != nil {
				c.Data["json"] = map[string]string{"error": "Invalid service ID: " + serviceIDStr}
				c.Ctx.Output.SetStatus(400)
				c.ServeJSON()
				return
			}
			serviceIds = append(serviceIds, id)
		}
	}

	university, err := models.GetUniversityByID(universityId)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "University not found: " + err.Error()}
		c.Ctx.Output.SetStatus(404)
		c.ServeJSON()
		return
	}

	// Update fields if they are provided
	if partialResponse.NameRu != "" {
		university.NameRu = partialResponse.NameRu
	}
	if partialResponse.NameKz != "" {
		university.NameKz = partialResponse.NameKz
	}
	if partialResponse.UniversityStatusRu != "" {
		university.UniversityStatusRu = partialResponse.UniversityStatusRu
	}
	if partialResponse.UniversityStatusKz != "" {
		university.UniversityStatusKz = partialResponse.UniversityStatusKz
	}
	if partialResponse.Website != "" {
		university.Website = partialResponse.Website
	}
	if partialResponse.CallCenterNumber != "" {
		university.CallCenterNumber = partialResponse.CallCenterNumber
	}
	if partialResponse.WhatsAppNumber != "" {
		university.WhatsAppNumber = partialResponse.WhatsAppNumber
	}
	if partialResponse.Address != "" {
		university.Address = partialResponse.Address
	}
	if partialResponse.UniversityCode != "" {
		university.UniversityCode = partialResponse.UniversityCode
	}
	if partialResponse.StudyFormatRu != "" {
		university.StudyFormatRu = partialResponse.StudyFormatRu
	}
	if partialResponse.StudyFormatKz != "" {
		university.StudyFormatKz = partialResponse.StudyFormatKz
	}
	if partialResponse.AbbreviationRu != "" {
		university.AbbreviationRu = partialResponse.AbbreviationRu
	}
	if partialResponse.AbbreviationKz != "" {
		university.AbbreviationKz = partialResponse.AbbreviationKz
	}
	if partialResponse.AddressLink != "" {
		university.AddressLink = partialResponse.AddressLink
	}
	if partialResponse.DescriptionRu != "" {
		university.DescriptionRu = partialResponse.DescriptionRu
	}
	if partialResponse.DescriptionKz != "" {
		university.DescriptionKz = partialResponse.DescriptionKz
	}
	if partialResponse.Rating != "" {
		university.Rating = partialResponse.Rating
	}
	if partialResponse.MinScore != 0 {
		university.MinEntryScore = partialResponse.MinScore
	}
	if partialResponse.CityId != 0 {
		university.City.Id = partialResponse.CityId
	}

	file, header, err := c.GetFile("MainImageUrl")
	if err == nil {
		defer file.Close()
		mainImagePath := fmt.Sprintf("Universities/%d/%s", universityId, header.Filename)
		uploadedMainImageURL, err := models.UploadFileToCloud(mainImagePath, file)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Failed to upload main image: " + err.Error()}
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}
		university.MainImageUrl = uploadedMainImageURL
	}

	galleryFiles := c.Ctx.Request.MultipartForm.File["Gallery"]
	var galleryURLs []string

	for _, galleryFileHeader := range galleryFiles {
		galleryFile, err := galleryFileHeader.Open()
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Failed to open gallery file: " + err.Error()}
			c.Ctx.Output.SetStatus(400)
			c.ServeJSON()
			return
		}
		defer galleryFile.Close()

		galleryFilePath := fmt.Sprintf("Universities/%d/Gallery/%s", universityId, galleryFileHeader.Filename)
		uploadedGalleryURL, err := models.UploadFileToCloud(galleryFilePath, galleryFile)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Failed to upload gallery image: " + err.Error()}
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}

		galleryURLs = append(galleryURLs, uploadedGalleryURL)
	}

	err = models.UpdateUniversityGallery(universityId, galleryURLs)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to update gallery images: " + err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	if err := models.UpdateUniversity(university); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to update university: " + err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	var services []*models.Service
	for _, serviceID := range serviceIds {
		service, err := models.GetServiceByID(serviceID)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Service not found: " + err.Error()}
			c.Ctx.Output.SetStatus(404)
			c.ServeJSON()
			return
		}
		services = append(services, service)
	}

	if err := models.UpdateUniversityServices(universityId, services); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to update university services: " + err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"status": "success"}
	c.Ctx.Output.SetStatus(200)
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
// @router /assignspec/:universityId/:specialityId [post]
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
// @router /assignspecialities/:universityId [post]
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
// @router /assignserv/:universityId [post]
func (c *UniversityController) AddServicesToUniversity() {
	universityId, _ := c.GetInt(":universityId")
	_ = c.Ctx.Input.CopyBody(512)
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
// @Success 200 {object} models.UniversitySearchResult "Список найденных университетов с информацией о пагинации"
// @Failure 400 {string} string "400 ошибка поиска или другая ошибка"
// @router /search [get]
func (c *UniversityController) SearchUniversities() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}
	params := make(map[string]interface{})
	if minScore, err := c.GetInt("min_score"); err == nil {
		params["min_score"] = minScore
	}
	if avgFee, err := c.GetInt("avg_fee"); err == nil {
		params["avg_fee"] = avgFee
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
	if sort := c.GetString("sort"); sort == "name_asc" || sort == "name_desc" {
		params["sort"] = sort
	}
	if name := c.GetString("name"); name != "" {
		params["name"] = name
	}
	if status := c.GetString("status"); status != "" {
		params["status"] = status
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

	if term, err := c.GetInt("term"); err == nil {
		params["term"] = term
	}

	log.Printf("Received parameters map: %+v", params)

	result, err := models.SearchUniversities(params, language)
	if err == nil {
		c.Data["json"] = result
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// DeleteSpecialityFromUniversity
// @Title DeleteSpecialityFromUniversity
// @Description удаляет взаимосвязь между университетом и специальностью
// @Param   university_id    path    int     true        "ID университета"
// @Param   speciality_id    path    int     true        "ID специальности"
// @Success 200 {string} "Success"
// @Failure 400 {string} "Invalid input"
// @Failure 404 {string} "Not Found"
// @router /deletespec/:universityId/:speciality_id [delete]
func (c *UniversityController) DeleteSpecialityFromUniversity() {
	universityIDStr := c.Ctx.Input.Param(":university_id")
	specialityIDStr := c.Ctx.Input.Param(":speciality_id")

	universityID, err := strconv.Atoi(universityIDStr)
	if err != nil {
		fmt.Println(universityID)
		c.CustomAbort(http.StatusBadRequest, "Invalid university ID")
		return
	}

	specialityID, err := strconv.Atoi(specialityIDStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid speciality ID")
		return
	}

	o := orm.NewOrm()

	sql := `DELETE FROM speciality_university WHERE university_id = ? AND speciality_id = ?`
	res, err := o.Raw(sql, universityID, specialityID).Exec()
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	num, err := res.RowsAffected()
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	if num == 0 {
		c.CustomAbort(http.StatusNotFound, "No relation found between university and speciality")
		return
	}

	c.Data["json"] = "Success"
	c.ServeJSON()
}

// DeleteGalleryPhoto deletes a gallery photo from a university.
// @Title DeleteGalleryPhoto
// @Description Deletes a photo from the university's gallery by photo ID and university ID.
// @Param	uniId		path	int		true	"University ID"
// @Param	photoId		path	int		true	"Photo ID"
// @Success 200 {object} map[string]string "message": "Photo deleted successfully"
// @Failure 400 {string} string "Invalid university ID or Invalid photo ID"
// @Failure 404 {string} string "University not found"
// @Failure 500 {string} string "Internal Server Error"
// @router /{uniId}/gallery/{photoId} [delete]
func (c *UniversityController) DeleteGalleryPhoto() {
	uniId, err := c.GetInt(":uniId")
	if err != nil {
		c.CustomAbort(400, "Invalid university ID")
		return
	}

	photoId, err := c.GetInt(":photoId")
	if err != nil {
		c.CustomAbort(400, "Invalid photo ID")
		return
	}

	o := orm.NewOrm()

	university := &models.University{Id: uniId}
	if err := o.Read(university); err != nil {
		c.CustomAbort(404, "University not found")
		return
	}

	if err := university.RemoveGalleryPhoto(photoId); err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = map[string]string{"message": "Photo deleted successfully"}
	c.ServeJSON()
}

// GetUniNames
// @Title Get University Names by Language
// @Description Returns a list of university IDs and names in the specified language.
// @Param	lang	query	string	false	"Language code (e.g., 'ru' for Russian, 'kz' for Kazakh). Defaults to 'kz' if not provided."
// @Success 200 {array} models.GetUniNamesResponse "Successful response with a list of university IDs and names."
// @Failure 500 {object} map[string]string "Internal Server Error with an error message."
// @Router /universities/names [get]
func (c *UniversityController) GetUniNames() {
	lang := c.Ctx.Input.Header("lang")
	if lang == "" {
		lang = "kz"
	}

	unis, err := models.GetUniversityNames(lang)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = unis
	c.ServeJSON()
}
