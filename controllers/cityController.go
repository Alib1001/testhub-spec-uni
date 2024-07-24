package controllers

import (
	"encoding/json"
	"testhub-spec-uni/models"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type CityResponse struct {
	Id        int       `json:"Id"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type CityController struct {
	beego.Controller
}

// Create adds a new city to the database.
// @Title Create
// @Description Создание нового города.
// @Param	body	body	models.City	true	"JSON с данными о городе"
// @Success 200 {object} map[string]int64 {"id": 1} "ID созданного города"
// @Failure 400 {string} string "400 ошибка разбора JSON или другая ошибка"
// @router / [post]
func (c *CityController) Create() {
	var city models.City

	requestBody := c.Ctx.Input.CopyBody(1024)

	err := json.Unmarshal(requestBody, &city)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	id, err := models.AddCity(&city)
	if err == nil {
		city.Id = int(id)
		c.Data["json"] = map[string]int64{"id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Get возвращает информацию о городе по его ID.
// @Title Get
// @Description Получение информации о городе по ID.
// @Param	id		path	int	true	"ID города для получения информации"
// @Param	lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} models.City "Информация о городе"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /:id [get]
func (c *CityController) Get() {
	id, _ := c.GetInt(":id")
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.Data["json"] = "Invalid or unsupported language"
		c.ServeJSON()
		return
	}

	city, err := models.GetCityById(id, language)
	if err == nil {
		response := CityResponse{
			Id:        city.Id,
			Name:      city.Name,
			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		}
		c.Data["json"] = response
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetAll возвращает список всех городов на указанном языке.
// @Title GetAll
// @Description Получение списка всех городов на указанном языке.
// @Param  lang  header  string  true  "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.City "Список городов"
// @Failure 400 {string} string "400 ошибка получения списка или другая ошибка"
// @router / [get]
func (c *CityController) GetAll() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.Data["json"] = "Invalid or unsupported language"
		c.ServeJSON()
		return
	}

	cities, err := models.GetAllCitiesByLanguage(language)
	if err == nil {
		var response []CityResponse
		for _, city := range cities {
			response = append(response, CityResponse{
				Id:        city.Id,
				Name:      city.Name,
				CreatedAt: city.CreatedAt,
				UpdatedAt: city.UpdatedAt,
			})
		}
		c.Data["json"] = response
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Update updates information about a city by its ID.
// @Title Update
// @Description Обновление информации о городе по ID.
// @Param	id		path	int	true	"ID города для обновления информации"
// @Param	body	body	models.City	true	"JSON с обновленными данными о городе"
// @Success 200 string "Обновление успешно выполнено"
// @Failure 400 {string} string "400 некорректный ID, ошибка разбора JSON или другая ошибка"
// @router /:id [put]
func (c *CityController) Update() {
	id, _ := c.GetInt(":id")
	var city models.City
	json.Unmarshal(c.Ctx.Input.RequestBody, &city)
	city.Id = id
	err := models.UpdateCity(&city)
	if err == nil {
		c.Data["json"] = "Update successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete deletes a city by its ID.
// @Summary Delete a city
// @Description Удаление города по ID
// @ID delete-city-by-id
// @Param id path int true "ID города для удаления"
// @Success 200 {string} string "Успешное удаление"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @Router /:id [delete]
func (c *CityController) Delete() {
	id, _ := c.GetInt(":id")
	err := models.DeleteCity(id)
	if err == nil {
		c.Data["json"] = "Delete successful"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetWithUniversities возвращает информацию о городе вместе с университетами по его ID.
// @Title GetWithUniversities
// @Description Получение информации о городе с университетами по ID.
// @Param	id		path	int	true	"ID города для получения информации"
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} models.City "Информация о городе с университетами"
// @Failure 400 {string} string "400 некорректный ID или другая ошибка"
// @router /info/:id [get]
func (c *CityController) GetWithUniversities() {
	id, _ := c.GetInt(":id")
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.Data["json"] = "Invalid or unsupported language"
		c.ServeJSON()
		return
	}

	city, err := models.GetCityWithUniversities(id, language)
	if err == nil {
		response := CityResponse{
			Id:        city.Id,
			Name:      city.Name,
			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		}
		c.Data["json"] = response
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// SearchCities ищет города по имени.
// @Title SearchCities
// @Description Поиск городов по имени.
// @Param	name		query	string	true	"Имя города для поиска"
// @Param	lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} CityResponse "Список найденных городов"
// @Failure 400 {string} string "400 ошибка поиска или другая ошибка"
// @router /search [get]
func (c *CityController) SearchCities() {
	name := c.GetString("name")
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.Data["json"] = "Invalid or unsupported language"
		c.ServeJSON()
		return
	}

	cities, err := models.SearchCitiesByName(name, language)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	var cityResponses []CityResponse
	for _, city := range cities {
		var cityName string
		switch language {
		case "ru":
			cityName = city.NameRu
		case "kz":
			cityName = city.NameKz
		default:
			cityName = city.Name
		}

		cityResponses = append(cityResponses, CityResponse{
			Id:        city.Id,
			Name:      cityName,
			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		})
	}

	c.Data["json"] = cityResponses
	c.ServeJSON()
}
