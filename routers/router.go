// @APIVersion 1.0.2
// @Title Testhub universities  API
// @Description API for Testhub universities.
// @Contact superalibek123@gmail.com

package routers

import (
	"testhub-spec-uni/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/subjects",
			beego.NSInclude(&controllers.SubjectController{}),
			beego.NSRouter("/", &controllers.SubjectController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SubjectController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SubjectController{}, "delete:Delete"),
			beego.NSRouter("/secubjects/:firstSubjectId", &controllers.SubjectController{}, "get:GetAllowedSecondSubjects"),
			beego.NSRouter("/search", &controllers.SubjectController{}, "get:SearchSubjectsByName"),
		),
		beego.NSNamespace("/subjectpairs",
			beego.NSInclude(&controllers.SubjectPairController{}),
			beego.NSRouter("/add/:firstSubjectId/:secondSubjectId", &controllers.SubjectPairController{}, "post:Add"),
			beego.NSRouter("/:id", &controllers.SubjectPairController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SubjectPairController{}, "get:GetAll"),
			beego.NSRouter("/:id/:firstSubjectId/:secondSubjectId", &controllers.SubjectPairController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SubjectPairController{}, "delete:Delete"),
			beego.NSRouter("/get/:firstSubjectId/:secondSubjectId", &controllers.SubjectPairController{}, "get:GetBySubjectIds"),

			//beego.NSRouter("/by_speciality/:specialityId", &controllers.SubjectPairController{}, "get:GetSubjectPairsBySpecialityID"),
		),

		beego.NSNamespace("/specialities",
			beego.NSInclude(&controllers.SpecialityController{}),
			beego.NSRouter("/", &controllers.SpecialityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.SpecialityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.SpecialityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.SpecialityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.SpecialityController{}, "delete:Delete"),
			beego.NSRouter("/search", &controllers.SpecialityController{}, "get:SearchSpecialities"),
			beego.NSRouter("/byuni/:universityId", &controllers.SpecialityController{}, "get:GetByUniversity"),
			beego.NSRouter("/bysubjects/:subject1_id/:subject2_id", &controllers.SpecialityController{}, "get:GetSpecialitiesBySubjectPair"),
			beego.NSRouter("/associatepair/:speciality_id/:subject_pair_id", &controllers.SpecialityController{}, "put:AssociateSpecialityWithSubjectPair"),
			beego.NSRouter("/byspec/:speciality_id", &controllers.SpecialityController{}, "get:GetSubjectPairsBySpecialityId"),
			beego.NSRouter("/addpointstat/:universityId/:specialityId", &controllers.SpecialityController{}, "post:AddPointStat"),
			beego.NSRouter("/pointstatsbyparams/:universityId/:specialityId", &controllers.SpecialityController{}, "get:GetPointStatsByUniversityAndSpeciality"),

			//beego.NSRouter("/subject_combinations/:id", &controllers.SpecialityController{}, "get:GetSubjectsCombinationForSpeciality"),
			//beego.NSRouter("/:specialityId/subjects/:subjectId", &controllers.SpecialityController{}, "post:AddSubject"),
		),
		beego.NSNamespace("/universities",
			beego.NSInclude(&controllers.UniversityController{}),
			beego.NSRouter("/", &controllers.UniversityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.UniversityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.UniversityController{}, "delete:Delete"),
			beego.NSRouter("/assigncity/:universityId/:cityId", &controllers.UniversityController{}, "put:AssignCityToUniversity"),
			beego.NSRouter("/assignspec/:universityId/:specialityId", &controllers.UniversityController{}, "post:AddSpecialityToUniversity"),
			beego.NSRouter("/assignspecialities/:universityId", &controllers.UniversityController{}, "post:AddSpecialitiesToUniversity"),
			beego.NSRouter("/assignserv/:universityId", &controllers.UniversityController{}, "post:AddServicesToUniversity"),
			beego.NSRouter("/search", &controllers.UniversityController{}, "get:SearchUniversities"),
			beego.NSRouter("/deletespec/:university_id/:speciality_id", &controllers.UniversityController{}, "delete:DeleteSpecialityFromUniversity"),
		),

		beego.NSNamespace("/cities",
			beego.NSInclude(&controllers.CityController{}),
			beego.NSRouter("/", &controllers.CityController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.CityController{}, "get:Get"),
			beego.NSRouter("/", &controllers.CityController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.CityController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.CityController{}, "delete:Delete"),
			beego.NSRouter("/info/:id", &controllers.CityController{}, "get:GetWithUniversities"),
			beego.NSRouter("/search", &controllers.CityController{}, "get:SearchCities"),
		),

		beego.NSNamespace("/quotas",
			beego.NSInclude(&controllers.QuotaController{}),
			beego.NSRouter("/", &controllers.QuotaController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "get:Get"),
			beego.NSRouter("/", &controllers.QuotaController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.QuotaController{}, "delete:Delete"),
			beego.NSRouter("/all/:id", &controllers.QuotaController{}, "get:GetQuotaWithSpecialities"),
			beego.NSRouter("/:quota_id/specialities/:speciality_id", &controllers.QuotaController{}, "post:AddSpecialityToQuota"),
		),

		beego.NSNamespace("/services",
			beego.NSInclude(&controllers.ServiceController{}),
			beego.NSRouter("/search", &controllers.ServiceController{}, "get:SearchServices"),
			beego.NSRouter("/", &controllers.ServiceController{}, "post:AddService"),
			beego.NSRouter("/", &controllers.ServiceController{}, "get:GetAllServices"),
			beego.NSRouter("/:id", &controllers.ServiceController{}, "get:GetServiceById"),
			beego.NSRouter("/:id", &controllers.ServiceController{}, "delete:DeleteService"),
			beego.NSRouter("/:id", &controllers.ServiceController{}, "put:UpdateService"),
			beego.NSRouter("/bind/:serviceId/:universityId", &controllers.ServiceController{}, "post:AddServiceToUniversity"),
			beego.NSRouter("/getbyuni/:universityId", &controllers.ServiceController{}, "get:GetServicesByUniversityId"),
		),

		/**
		beego.NSNamespace("/users",
			beego.NSRouter("/", &controllers.UserController{}),
			beego.NSRouter("/:id", &controllers.UserController{}, "get:GetUserByID"),
			beego.NSRouter("/:id", &controllers.UserController{}, "put:UpdateUserByID"),
			beego.NSRouter("/:id", &controllers.UserController{}, "delete:DeleteUser"),
		),

		beego.NSRouter("/login", &controllers.AuthController{}, "get:GetLogin;post:PostLogin"),
		**/
	)

	beego.AddNamespace(ns)
}
