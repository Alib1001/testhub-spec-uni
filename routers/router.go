package routers

import (
	"testhub-spec-uni/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/subjects",
			beego.NSRouter("/", &controllers.SubjectController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SubjectController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "delete:Delete"),
		),

		beego.NSNamespace("/specialities",
			beego.NSRouter("/", &controllers.SpecialityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.SpecialityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SpecialityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.SpecialityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SpecialityController{}, "delete:Delete"),
			beego.NSRouter("/:specialityId/subjects/:subjectId", &controllers.SpecialityController{}, "post:AddSubject"),
		),

		beego.NSNamespace("/universities",
			beego.NSRouter("/", &controllers.UniversityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.UniversityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "delete:Delete"),
		),

		beego.NSNamespace("/cities",
			beego.NSRouter("/", &controllers.CityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.CityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.CityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.CityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.CityController{}, "delete:Delete"),
		),

		beego.NSNamespace("/quotas",
			beego.NSRouter("/", &controllers.QuotaController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "get:Get"),
			beego.NSRouter("/", &controllers.QuotaController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "delete:Delete"),
			beego.NSRouter("/all/:id", &controllers.QuotaController{}, "get:GetQuotaWithSpecialities"),
			beego.NSRouter("/:quota_id/specialities/:speciality_id", &controllers.QuotaController{}, "post:AddSpecialityToQuota"),
		),

		beego.NSNamespace("/users",
			beego.NSRouter("/login", &controllers.AuthController{}, "get:GetLogin;post:PostLogin"),
			beego.NSRouter("/", &controllers.UserController{}),
			beego.NSRouter("/:id", &controllers.UserController{}, "get:GetUserByID"),
			beego.NSRouter("/:id", &controllers.UserController{}, "put:UpdateUserByID"),
			beego.NSRouter("/:id", &controllers.UserController{}, "delete:DeleteUser"),
		),
	)

	beego.AddNamespace(ns)
}
