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
// @Success 200 {array} models.Service
// @Failure 500 Internal server error
// @router / [get]
func (c *ServiceController) GetAllServices() {
	services, err := models.GetAllServices()
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = services
	c.ServeJSON()
}

// @Title GetServiceById
// @Description Get service by ID
// @Param   id    path    int   true        "Service ID"
// @Success 200 {object} models.Service
// @Failure 404 Not found
// @router /:id [get]
func (c *ServiceController) GetServiceById() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
	}

	service, err := models.GetServiceById(id)
	if err != nil {
		c.CustomAbort(http.StatusNotFound, err.Error())
	}

	c.Data["json"] = service
	c.ServeJSON()
}

// GetServicesByUniversityId retrieves services associated with a university by its ID
// @Title GetServicesByUniversityId
// @Description Retrieve services associated with a university by its ID
// @Param   universityId     query    int   true        "University ID"
// @Success 200 {object}  []Service
// @Failure 400 Invalid university ID
// @router /getbyuni/:universityId [get]
func (c *ServiceController) GetServicesByUniversityId() {
	universityIdStr := c.GetString(":universityId")
	universityId, err := strconv.Atoi(universityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid university ID")
	}

	services, err := models.GetServicesByUniversityId(universityId)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
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
	}

	id, err := models.AddService(&service)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = id
	c.ServeJSON()
}

// @Title DeleteService
// @Description Delete a service by ID
// @Param   id    path    int   true        "Service ID"
// @Success 200 {string} delete success!
// @Failure 404 Not found
// @router /:id [delete]
func (c *ServiceController) DeleteService() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
	}

	err = models.DeleteService(id)
	if err != nil {
		c.CustomAbort(http.StatusNotFound, err.Error())
	}

	c.Data["json"] = "delete success!"
	c.ServeJSON()
}

// @Title UpdateService
// @Description Update an existing service
// @Param   body    body    models.Service   true        "Service data"
// @Success 200 {string} update success!
// @Failure 400 Invalid input
// @router / [put]
func (c *ServiceController) UpdateService() {
	var service models.Service
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &service); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid input")
	}

	err := models.UpdateService(&service)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = "update success!"
	c.ServeJSON()
}

// @Title AddServiceToUniversity
// @Description Bind a service to a university
// @Param   serviceId    path    int   true        "Service ID"
// @Param   universityId path    int   true        "University ID"
// @Success 200 {string} bind success!
// @Failure 400 Invalid input
// @router /bind/:serviceId/:universityId [post]
func (c *ServiceController) AddServiceToUniversity() {
	serviceIdStr := c.Ctx.Input.Param(":serviceId")
	universityIdStr := c.Ctx.Input.Param(":universityId")

	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid service ID")
	}

	universityId, err := strconv.Atoi(universityIdStr)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid university ID")
	}

	err = models.AddServiceToUniversity(serviceId, universityId)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = "bind success!"
	c.ServeJSON()
}

// @Title SearchServices
// @Description Search for services by name
// @Param   prefix    query    string   true        "Prefix for service name"
// @Success 200 {array} models.Service
// @Failure 400 Invalid input
// @router /search [get]
func (c *ServiceController) SearchServices() {
	prefix := c.GetString("prefix")
	if prefix == "" {
		c.CustomAbort(http.StatusBadRequest, "Prefix is required")
	}

	services, err := models.SearchServicesByName(prefix)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = services
	c.ServeJSON()
}
