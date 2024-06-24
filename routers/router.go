package routers

import (
	"testhub-spec-uni/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Subjects
	beego.Router("/subjects", &controllers.SubjectController{}, "post:Create")
	beego.Router("/subjects/:id", &controllers.SubjectController{}, "get:Get")
	beego.Router("/subjects", &controllers.SubjectController{}, "get:GetAll")
	beego.Router("/subjects/:id", &controllers.SubjectController{}, "put:Update")
	beego.Router("/subjects/:id", &controllers.SubjectController{}, "delete:Delete")

	// Specialties
	beego.Router("/specialties", &controllers.SpecialtyController{}, "post:Create")
	beego.Router("/specialties/:id", &controllers.SpecialtyController{}, "get:Get")
	beego.Router("/specialties", &controllers.SpecialtyController{}, "get:GetAll")
	beego.Router("/specialties/:id", &controllers.SpecialtyController{}, "put:Update")
	beego.Router("/specialties/:id", &controllers.SpecialtyController{}, "delete:Delete")

	// Universities
	beego.Router("/universities", &controllers.UniversityController{}, "post:Create")
	beego.Router("/universities/:id", &controllers.UniversityController{}, "get:Get")
	beego.Router("/universities", &controllers.UniversityController{}, "get:GetAll")
	beego.Router("/universities/:id", &controllers.UniversityController{}, "put:Update")
	beego.Router("/universities/:id", &controllers.UniversityController{}, "delete:Delete")

	// Cities
	beego.Router("/cities", &controllers.CityController{}, "post:Create")
	beego.Router("/cities/:id", &controllers.CityController{}, "get:Get")
	beego.Router("/cities", &controllers.CityController{}, "get:GetAll")
	beego.Router("/cities/:id", &controllers.CityController{}, "put:Update")
	beego.Router("/cities/:id", &controllers.CityController{}, "delete:Delete")

	// Quotas
	beego.Router("/quotas", &controllers.QuotaController{}, "post:Create")
	beego.Router("/quotas/:id", &controllers.QuotaController{}, "get:Get")
	beego.Router("/quotas", &controllers.QuotaController{}, "get:GetAll")
	beego.Router("/quotas/:id", &controllers.QuotaController{}, "put:Update")
	beego.Router("/quotas/:id", &controllers.QuotaController{}, "delete:Delete")

}
