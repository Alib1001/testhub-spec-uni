package controllers

import (
	"encoding/json"
	"fmt"
	"testhub-spec-uni/models"

	"github.com/beego/beego/v2/server/web"
)

// UserController handles operations related to users.

type UserController struct {
	web.Controller
}

// @Title Get
// @Description Get all users
// @Success 200 {object} []models.User
// @Failure 500 Internal Server Error
// @router / [get]
func (c *UserController) Get() {
	users, err := models.GetAllUsers()
	if err != nil {
		c.Abort("500")
	}
	c.Data["json"] = users

	c.ServeJSON()
}

// @Title GetUserByID
// @Description Get user by ID
// @Param   id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 500 Internal Server Error
// @router /:id [get]
func (c *UserController) GetUserByID() {
	id, _ := c.GetInt(":id")
	user, err := models.GetUserById(id)
	if err != nil {
		c.Abort("500")
	}

	c.Data["json"] = user
	c.ServeJSON()
}

// @Title CreateUser
// @Description Create a new user
// @Param   user body models.User true "User object that needs to be added"
// @Success 201 {string} message
// @Failure 400 Bad request
// @Failure 409 Conflict - Username already taken
// @router / [post]
func (c *UserController) Post() {
	requestBody := string(c.Ctx.Input.CopyBody(1024))

	var user models.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		c.Abort("400")
		fmt.Println("Request Body:", requestBody)
		return
	}

	existingUser, _ := models.GetUserByUsername(user.Username)
	if existingUser.Id != 0 {
		c.Ctx.Output.SetStatus(409)
		c.Data["json"] = map[string]string{"error": "Username already taken"}
		c.ServeJSON()
		return
	}

	err = models.AddUser(user)
	if err != nil {
		c.Abort("500")
		return
	}
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]string{"message": "User created successfully"}
	c.ServeJSON()
}

// @Title UpdateUser
// @Description Update user by ID
// @Param   id path int true "User ID"
// @Param   user body models.User true "Updated user object"
// @Success 200 {string} message
// @Failure 400 Bad request
// @Failure 404 User not found
// @router /:id [put]
func (c *UserController) Put() {
	var user models.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Abort("400")
	}
	err = models.UpdateUser(user)
	if err != nil {
		c.Abort("500")
	}
}

// @Title UpdateUserByID
// @Description Update user by ID
// @Param   id path int true "User ID"
// @Param   user body models.User true "Updated user object"
// @Success 200 {string} message
// @Failure 400 Bad request
// @Failure 404 User not found
// @router /:id [put]
func (c *UserController) UpdateUserByID() {
	requestBody := string(c.Ctx.Input.CopyBody(1024))

	id, _ := c.GetInt(":id")

	var user models.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Abort("400")
	}

	existingUser, _ := models.GetUserByUsername(user.Username)
	if existingUser.Id != 0 && existingUser.Id != id {
		c.Ctx.Output.SetStatus(409)
		c.Data["json"] = map[string]string{"error": "Username already taken"}
		c.ServeJSON()
		return
	}

	user.Id = id
	err = models.UpdateUser(user)
	if err != nil {
		c.Abort("500")
	}
	fmt.Println("Request Body:", requestBody)
}

// @Title DeleteUserByID
// @Description Delete user by ID
// @Param   id path int true "User ID"
// @Success 200 {string} message
// @Failure 404 User not found
// @router /:id [delete]
func (c *UserController) DeleteUser() {
	id, _ := c.GetInt(":id")
	err := models.DeleteUser(id)
	if err != nil {
		c.Abort("500")
	}
}
