package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-playground/validator/v10"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"testhub-spec-uni/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UniversityController обрабатывает запросы для работы с университетами.
type UniversityController struct {
	beego.Controller
}

func (c *UniversityController) Create() {
	err := c.Ctx.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse form data: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	var universityResponse models.AddUUniversityResponse

	if err := c.ParseForm(&universityResponse); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse form data: " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

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

	// Handle main image upload
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

	uploadedMainImageURL, err := uploadFileToCloud(mainImagePath, file)
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

	// Handle gallery images upload
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
		uploadedGalleryURL, err := uploadFileToCloud(galleryFilePath, galleryFile)
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

func uploadFileToCloud(filePath string, file multipart.File) (string, error) {
	awsAccessKey, _ := beego.AppConfig.String("aws_access_key")
	awsSecretKey, _ := beego.AppConfig.String("aws_secret_key")
	bucket, _ := beego.AppConfig.String("bucket")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		Endpoint:         aws.String("https://chi-sextans.object.pscloud.io"),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	uploader := s3.New(sess)

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	fileURL := fmt.Sprintf("https://chi-sextans.object.pscloud.io/%s/%s", bucket, filePath)
	return fileURL, nil
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
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}
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
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}
	universities, err := models.GetAllUniversities(language)
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
// @Success 200 {string} "Обновление успешно выполнено"
// @Failure 400 {string} "Некорректный ID, ошибка разбора JSON или другая ошибка"
// @router /:id [put]
func (c *UniversityController) Update() {
	// Получение ID из URL-параметра
	id, err := c.GetInt(":id")
	requestBody := c.Ctx.Input.CopyBody(2048)
	if err != nil {
		c.Data["json"] = "Некорректный ID"
		c.ServeJSON()
		return
	}

	// Чтение и разбор JSON тела запроса
	var university models.University
	err = json.Unmarshal(requestBody, &university)
	if err != nil {
		c.Data["json"] = "Ошибка разбора JSON: " + err.Error()
		c.ServeJSON()
		return
	}

	// Установка ID университета
	university.Id = id

	// Обновление информации о университете
	err = models.UpdateUniversityFields(&university)
	if err != nil {
		c.Data["json"] = "Ошибка обновления: " + err.Error()
		c.ServeJSON()
		return
	}

	// Успешное обновление
	c.Data["json"] = "Обновление успешно выполнено"
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

	result, err := models.SearchUniversities(params, language)
	if err == nil {
		c.Data["json"] = result
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
