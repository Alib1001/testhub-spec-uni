package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testhub-spec-uni/models"

	beego "github.com/beego/beego/v2/server/web"
)

// ServiceController handles operations related to services
type ServiceController struct {
	beego.Controller
}

// @Title GetAllServices
// @Description Get all services
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.ServiceResponse
// @Failure 500 Internal server error
// @router / [get]
func (c *ServiceController) GetAllServices() {
	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	services, err := models.GetAllServices(language)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = services
	c.ServeJSON()
}

// @Title GetServiceById
// @Description Get service by ID
// @Param id path int true "Service ID"
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {object} models.ServiceResponse
// @Failure 404 Not found
// @router /:id [get]
func (c *ServiceController) GetServiceById() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
		return
	}

	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	service, err := models.GetServiceById(id, language)
	if err != nil {
		c.CustomAbort(http.StatusNotFound, err.Error())
		return
	}

	c.Data["json"] = service
	c.ServeJSON()
}

// @Title SearchServices
// @Description Search for services by name
// @Param prefix query string true "Prefix for service name"
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.ServiceResponse
// @Failure 400 Invalid input
// @router /search [get]
func (c *ServiceController) SearchServices() {
	prefix := c.GetString("prefix")
	if prefix == "" {
		c.CustomAbort(http.StatusBadRequest, "Prefix is required")
		return
	}

	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	services, err := models.SearchServicesByName(prefix, language)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = services
	c.ServeJSON()
}

// @Title AddService
// @Description Add a new service
// @Param   body    body    models.Service   true        "Service data"
// @Success 200 {int} models.Service.Id
// @Failure 400 Invalid input
// @router / [post]
func (c *ServiceController) AddService() {
	_ = c.Ctx.Input.CopyBody(1024)
	var service models.Service
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &service); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := models.AddService(&service)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = id
	c.ServeJSON()
}

// @Title UpdateService
// @Description Update service information
// @Param id path int true "Service ID"
// @Param body body models.Service true "Service data"
// @Success 200 {string} "Update successful"
// @Failure 400 Invalid input
// @router /:id [put]
func (c *ServiceController) UpdateService() {
	idStr := c.Ctx.Input.Param(":id")
	_ = c.Ctx.Input.CopyBody(1024)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
		return
	}

	var service models.Service
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &service); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	service.Id = id

	var updatedFields map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &updatedFields); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	fields := make([]string, 0, len(updatedFields))
	for field := range updatedFields {
		fields = append(fields, field)
	}

	if err := models.UpdateService(&service, fields...); err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = "Update successful"
	c.ServeJSON()
}

// @Title DeleteService
// @Description Delete a service by ID
// @Param id path int true "Service ID"
// @Success 200 {string} "Delete successful"
// @Failure 400 Invalid service ID
// @router /:id [delete]
func (c *ServiceController) DeleteService() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
		return
	}

	if err := models.DeleteService(id); err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = "Delete successful"
	c.ServeJSON()
}

// @Title GetServicesByUniversityId
// @Description Get services by university ID
// @Param id path int true "University ID"
// @Param lang header string true "Язык для получения данных, 'ru' или 'kz'"
// @Success 200 {array} models.ServiceResponse
// @Failure 400 Invalid input
// @router /university/:id [get]
func (c *ServiceController) GetServicesByUniversityId() {
	universityIdStr := c.Ctx.Input.Param(":id")
	universityId, err := strconv.Atoi(universityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid university ID")
		return
	}

	language := c.Ctx.Input.Header("lang")
	if language != "ru" && language != "kz" {
		c.CustomAbort(http.StatusBadRequest, "Invalid or unsupported language")
		return
	}

	services, err := models.GetServicesByUniversityId(universityId, language)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = services
	c.ServeJSON()
}

// @Title AddServiceToUniversity
// @Description Add a service to a university
// @Param serviceId path int true "Service ID"
// @Param universityId path int true "University ID"
// @Success 200 {string} "Add successful"
// @Failure 400 Invalid input
// @router /:serviceId/university/:universityId [post]
func (c *ServiceController) AddServiceToUniversity() {
	serviceIdStr := c.Ctx.Input.Param(":serviceId")
	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
		return
	}

	universityIdStr := c.Ctx.Input.Param(":universityId")
	universityId, err := strconv.Atoi(universityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid university ID")
		return
	}

	if err := models.AddServiceToUniversity(serviceId, universityId); err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = "Add successful"
	c.ServeJSON()
}
