package routers

import (
	"testhub-spec-uni/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/subjects",
			beego.NSInclude(&controllers.SubjectController{}),
			beego.NSRouter("/", &controllers.SubjectController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SubjectController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "delete:Delete"),
		),

		beego.NSNamespace("/specialties",
			beego.NSInclude(&controllers.SpecialtyController{}),
			beego.NSRouter("/", &controllers.SpecialtyController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.SpecialtyController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SpecialtyController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.SpecialtyController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SpecialtyController{}, "delete:Delete"),
			beego.NSRouter("/:specialtyId/subjects/:subjectId", &controllers.SpecialtyController{}, "post:AddSubject"),
		),

		beego.NSNamespace("/universities",
			beego.NSInclude(&controllers.UniversityController{}),
			beego.NSRouter("/", &controllers.UniversityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.UniversityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "delete:Delete"),
		),

		beego.NSNamespace("/cities",
			beego.NSInclude(&controllers.CityController{}),
			beego.NSRouter("/", &controllers.CityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.CityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.CityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.CityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.CityController{}, "delete:Delete"),
		),

		beego.NSNamespace("/quotas",
			beego.NSInclude(&controllers.QuotaController{}),
			beego.NSRouter("/", &controllers.QuotaController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "get:Get"),
			beego.NSRouter("/", &controllers.QuotaController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "delete:Delete"),
			beego.NSRouter("/:quotaId/specialties/:specialtyId", &controllers.QuotaController{}, "post:AddSpecialtyToQuota"),
		),
	)

	beego.AddNamespace(ns)
}
