package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Speciality struct {
	Id              int    `orm:"auto" json:"id"`
	Name            string `orm:"size(128)" json:"name"`
	NameRu          string `orm:"size(128)" json:"NameRu"`
	NameKz          string `orm:"size(128)" json:"NameKz"`
	AbbreviationRu  string `json:"AbbreviationRu" validate:"required"`
	AbbreviationKz  string `json:"AbbreviationKz" validate:"required"`
	Code            string `orm:"size(64)" json:"code"`
	VideoLink       string `orm:"size(256)" json:"video_link"`
	Description     string `orm:"type(text)" json:"description"`
	DescriptionRu   string `orm:"type(text)" json:"DescriptionRu"`
	DescriptionKz   string `orm:"type(text)" json:"DescriptionKz"`
	Degree          string `orm:"size(128)" json:"degree"`
	Scholarship     bool
	SubjectPair     *SubjectPair            `orm:"rel(fk);on_delete(set_null);null" json:"subject_pair,omitempty"`
	UniversityTerms []*SpecialityUniversity `orm:"reverse(many)" json:"university_terms,omitempty"` // Множественная связь через SpecialityUniversity
	CreatedAt       time.Time               `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt       time.Time               `orm:"auto_now;type(datetime)" json:"updated_at"`
	PointStats      []*PointStat            `orm:"reverse(many)" json:"point_stats,omitempty"`
}

type GetSpecialityResponse struct {
	Id          int          `orm:"auto" json:"id"`
	Name        string       `orm:"size(128)" json:"name"`
	Degree      string       `orm:"size(128)" json:"degree"`
	SubjectPair *SubjectPair `orm:"rel(fk);on_delete(set_null);null" json:"subject_pair,omitempty"`
	PointStats  []*PointStat `orm:"reverse(many)" json:"point_stats,omitempty"`
	AvgSalary   int          `json:"avg_salary"`
	Code        string       `orm:"size(64)" json:"code"`
	Scholarship bool
}

type GetSpecialityForAdmResponse struct {
	Id             int          `orm:"auto" json:"id"`
	NameRu         string       `orm:"size(128)" json:"NameRu"`
	NameKz         string       `orm:"size(128)" json:"NameKz"`
	DescriptionRu  string       `orm:"type(text)" json:"DescriptionRu"`
	DescriptionKz  string       `orm:"type(text)" json:"DescriptionKz"`
	AbbreviationRu string       `json:"AbbreviationRu" validate:"required"`
	AbbreviationKz string       `json:"AbbreviationKz" validate:"required"`
	Degree         string       `orm:"size(128)" json:"degree"`
	SubjectPair    *SubjectPair `orm:"rel(fk);on_delete(set_null);null" json:"subject_pair,omitempty"`
	PointStats     []*PointStat `orm:"reverse(many)" json:"point_stats,omitempty"`
	AvgSalary      int          `json:"avg_salary"`
	Code           string       `orm:"size(64)" json:"code"`
	Scholarship    bool
}

type AddSpecialityResponse struct {
	Id             int    `form:"Id"`
	NameRu         string `form:"NameRu" validate:"required"`
	NameKz         string `form:"NameKz" validate:"required"`
	AbbreviationRu string `form:"AbbreviationRu" validate:"required"`
	AbbreviationKz string `form:"AbbreviationKz" validate:"required"`
	Subject1       int    `form:"Subject_1" validate:"required"`
	Subject2       int    `form:"Subject_2" validate:"required"`
	Degree         string `form:"Degree" validate:"required"`
	Code           string `form:"Code" validate:"required"`
	DescriptionRu  string `form:"DescriptionRu"`
	DescriptionKz  string `form:"DescriptionKz"`
	Scholarship    bool   `form:"Scholarship"`
}
type UpdateSpecialityResponse struct {
	Id             int    `form:"id"`
	NameRu         string `form:"NameRu"`
	NameKz         string `form:"NameKz"`
	AbbreviationRu string `form:"AbbreviationRu"`
	AbbreviationKz string `form:"AbbreviationKz"`
	Degree         string `form:"Degree"`
	Code           string `form:"Code"`
	DescriptionRu  string `form:"DescriptionRu"`
	DescriptionKz  string `form:"DescriptionKz"`
	Scholarship    bool   `form:"Scholarship"`
	Subject1       int    `form:"Subject_1"`
	Subject2       int    `form:"Subject_2"`
}

type AnnualGrant struct {
	Year       int `json:"year"`
	GrantCount int `json:"grant_count"`
}
type AnnualPoints struct {
	Year          int `json:"year"`
	MinScore      int `json:"min_score"`
	MinGrantScore int `json:"min_grant_score"`
}

type GetByUniResponseForUser struct {
	SpecialityID    int            `orm:"column(speciality_id)" json:"speciality_id"`
	SpecialityName  string         `orm:"column(speciality_name)" json:"speciality_name"`
	UniversityName  string         `orm:"column(university_name)" json:"university_name"`
	EducationFormat string         `orm:"column(education_format)" json:"education_format"`
	Code            string         `orm:"column(code)" json:"speciality_code"`
	Price           int            `orm:"column(price)" json:"price"`
	Degree          string         `orm:"column(degree)" json:"degree"`
	Scholarship     bool           `orm:"column(scholarship)" json:"scholarship"`
	AvgSalary       int            `orm:"column(avg_salary)" json:"avg_salary"`
	MinScore        int            `json:"min_score" orm:"column(min_score)"`
	GrantCount      int            `json:"grant_count"`
	SubjectNames    []string       `json:"subject_names"`
	AnnualPoints    []AnnualPoints `json:"annual_points"`
	AnnualGrants    []AnnualGrant  `json:"annual_grants"`
	Term            int            `json:"term" orm:"column(term)"`
	Subject1ID      int            `orm:"column(subject1_id)" json:"-"`
	Subject2ID      int            `orm:"column(subject2_id)" json:"-"`
}

type GetByUniResponseForAdm struct {
	SpecialityID      int            `json:"speciality_id" orm:"column(speciality_id)"`
	SpecialityNameRu  string         `json:"speciality_name_ru" orm:"column(speciality_name_ru)"`
	SpecialityNameKz  string         `json:"speciality_name_kz" orm:"column(speciality_name_kz)"`
	UniversityNameRu  string         `json:"university_name_ru" orm:"column(university_name_ru)"`
	UniversityNameKz  string         `json:"university_name_kz" orm:"column(university_name_kz)"`
	EducationFormatRu string         `json:"education_format_ru" orm:"column(education_format_ru)"`
	EducationFormatKz string         `json:"education_format_kz" orm:"column(education_format_kz)"`
	Code              string         `orm:"column(code)" json:"speciality_code"`
	Price             int            `orm:"column(price)" json:"price"`
	Degree            string         `json:"degree" orm:"column(degree)"`
	Scholarship       string         `json:"scholarship" orm:"column(scholarship)"`
	AvgSalary         int            `json:"avg_salary" orm:"column(avg_salary)"`
	MinScore          int            `json:"min_score" orm:"column(min_score)"`
	GrantCount        int            `json:"grant_count" orm:"column(grant_count)"`
	Subject1ID        int            `json:"-" orm:"column(subject1_id)"`
	Subject2ID        int            `json:"-" orm:"column(subject2_id)"`
	Term              int            `json:"term" orm:"column(term)"`
	SubjectNamesRu    []string       `json:"subject_names_ru"`
	SubjectNamesKz    []string       `json:"subject_names_kz"`
	AnnualPoints      []AnnualPoints `json:"annual_points"`
	AnnualGrants      []AnnualGrant  `json:"annual_grants"`
}

type SpecialitySearchResult struct {
	Specialities []*Speciality `json:"specialities"`
	Page         int           `json:"page"`
	TotalPages   int           `json:"total_pages"`
	TotalCount   int           `json:"total_count"`
}

type IUniverStatistic struct {
	Id         int `json:"id"`
	GrantCount int `json:"grant_count"`
	MinScore   int `json:"min_score"`
	Year       int `json:"year"`
	AvgSalary  int `json:"avg_salary"`
	Price      int `json:"price"`
}

type IUniverSpecialtyShortcut struct {
	SpecialityID      int                `json:"speciality_id"`
	SpecialityNameRu  string             `json:"speciality_name_ru"`
	SpecialityNameKz  string             `json:"speciality_name_kz"`
	SpecialityCode    string             `json:"speciality_code"`
	Degree            string             `json:"degree"`
	EducationFormatRu string             `json:"education_format_ru"`
	EducationFormatKz string             `json:"education_format_kz"`
	Term              int                `json:"term"`
	EduLang           string             `json:"edu_lang"`
	Statistics        []IUniverStatistic `json:"statistics"`
}

func init() {
	orm.RegisterModel(new(Speciality))
}

func AddSpecialityFromFormData(data *AddSpecialityResponse) (int64, error) {
	o := orm.NewOrm()

	subject1 := Subject{Id: data.Subject1}
	subject2 := Subject{Id: data.Subject2}

	if err := o.Read(&subject1); err != nil {
		return 0, fmt.Errorf("subject 1 not found: %v", err)
	}
	if err := o.Read(&subject2); err != nil {
		return 0, fmt.Errorf("subject 2 not found: %v", err)
	}

	subjectPair := SubjectPair{
		Subject1: &subject1,
		Subject2: &subject2,
	}

	if _, err := o.Insert(&subjectPair); err != nil {
		return 0, err
	}

	speciality := Speciality{
		NameRu:         data.NameRu,
		NameKz:         data.NameKz,
		AbbreviationRu: data.AbbreviationRu,
		AbbreviationKz: data.AbbreviationKz,
		Degree:         data.Degree,
		Code:           data.Code,
		DescriptionKz:  data.DescriptionKz,
		DescriptionRu:  data.DescriptionRu,
		Scholarship:    data.Scholarship,
		SubjectPair:    &subjectPair,
	}

	id, err := o.Insert(&speciality)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetSpecialityById(id int, language string) (*Speciality, error) {
	o := orm.NewOrm()
	var speciality Speciality

	err := o.Raw(`SELECT * FROM speciality WHERE id = ?`, id).QueryRow(&speciality)
	if err != nil {
		return nil, err
	}

	switch language {
	case "ru":
		speciality.Name = speciality.NameRu
		speciality.Description = speciality.DescriptionRu
	case "kz":
		speciality.Name = speciality.NameKz
		speciality.Description = speciality.DescriptionKz
	}

	if speciality.SubjectPair != nil {
		var subjectPair SubjectPair
		err = o.Raw(`SELECT * FROM subject_pair WHERE id = ?`, speciality.SubjectPair.Id).QueryRow(&subjectPair)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}
		speciality.SubjectPair = &subjectPair
	}

	var pointStats []*PointStat
	_, err = o.Raw(`SELECT * FROM point_stat WHERE speciality_id = ?`, id).QueryRows(&pointStats)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}
	speciality.PointStats = pointStats

	return &speciality, nil
}

func GetAllSpecialities(language string) ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	_, err := o.QueryTable("speciality").All(&specialities)
	if err != nil {
		return nil, err
	}

	for _, speciality := range specialities {
		if speciality.SubjectPair != nil {
			err := o.Read(speciality.SubjectPair)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
		}

		switch language {
		case "ru":
			speciality.Name = speciality.NameRu
			speciality.Description = speciality.DescriptionRu
		case "kz":
			speciality.Name = speciality.NameKz
			speciality.Description = speciality.DescriptionKz
		}
	}

	return specialities, nil
}

func UpdateSpecialityFromFormData(data *UpdateSpecialityResponse) error {
	o := orm.NewOrm()
	o.Begin()

	speciality := Speciality{Id: data.Id}
	if err := o.Read(&speciality); err != nil {
		o.Rollback()
		return fmt.Errorf("speciality not found: %v", err)
	}

	if data.NameRu != "" {
		speciality.NameRu = data.NameRu
	}
	if data.NameKz != "" {
		speciality.NameKz = data.NameKz
	}
	if data.AbbreviationRu != "" {
		speciality.AbbreviationRu = data.AbbreviationRu
	}
	if data.AbbreviationKz != "" {
		speciality.AbbreviationKz = data.AbbreviationKz
	}
	if data.Degree != "" {
		speciality.Degree = data.Degree
	}
	if data.Code != "" {
		speciality.Code = data.Code
	}
	if data.DescriptionRu != "" {
		speciality.DescriptionRu = data.DescriptionRu
	}
	if data.DescriptionKz != "" {
		speciality.DescriptionKz = data.DescriptionKz
	}
	speciality.Scholarship = data.Scholarship

	if data.Subject1 != 0 && data.Subject2 != 0 {
		subject1 := Subject{Id: data.Subject1}
		subject2 := Subject{Id: data.Subject2}

		if err := o.Read(&subject1); err != nil {
			o.Rollback() // Rollback transaction on error
			return fmt.Errorf("subject 1 not found: %v", err)
		}
		if err := o.Read(&subject2); err != nil {
			o.Rollback() // Rollback transaction on error
			return fmt.Errorf("subject 2 not found: %v", err)
		}

		subjectPair := SubjectPair{
			Subject1: &subject1,
			Subject2: &subject2,
		}

		if _, err := o.Insert(&subjectPair); err != nil {
			o.Rollback() // Rollback transaction on error
			return fmt.Errorf("failed to insert subject pair: %v", err)
		}
		speciality.SubjectPair = &subjectPair
	}

	if _, err := o.Update(&speciality); err != nil {
		o.Rollback() // Rollback transaction on error
		return fmt.Errorf("failed to update speciality: %v", err)
	}

	o.Commit() // Commit transaction
	return nil
}

func DeleteSpeciality(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Speciality{Id: id})
	return err
}

func SearchSpecialities(params map[string]interface{}, language string) (*SpecialitySearchResult, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	_, err := o.QueryTable("speciality").All(&specialities)
	if err != nil {
		return nil, err
	}

	specialities, err = filterSpecialitiesBySubjectPair(params, specialities)
	if err != nil {
		return nil, err
	}

	specialities, err = filterSpecialitiesInUniversity(params, specialities)
	if err != nil {
		return nil, err
	}

	specialities, err = filterSpecialitiesByName(params, specialities, language)
	if err != nil {
		return nil, err
	}

	// Пагинация
	result, err := paginateSpecialities(specialities, params, language)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SearchSpecialities: total specialities after filtering: %d\n", len(result.Specialities))
	return result, nil
}
func GetSpecialitiesInUniversityForUser(universityId int, language string) ([]GetByUniResponseForUser, error) {
	o := orm.NewOrm()
	var results []GetByUniResponseForUser

	query := `
        WITH speciality_data AS (
            SELECT DISTINCT 
                s.id as speciality_id,
                CASE WHEN ? = 'ru' THEN s.name_ru ELSE s.name_kz END as speciality_name,
                CASE WHEN ? = 'ru' THEN u.name_ru ELSE u.name_kz END as university_name,
                CASE WHEN ? = 'ru' THEN u.study_format_ru ELSE u.study_format_kz END as education_format,
                s.code,
                ps.price,
                s.degree,
                s.scholarship,
                ps.avg_salary,
                sp.subject1_id,
                sp.subject2_id,
                ps.year,
                ps.min_score,
                ps.min_grant_score,
                ps.grant_count,
                su.term
            FROM 
                speciality s
                INNER JOIN speciality_university su ON s.id = su.speciality_id
                INNER JOIN university u ON su.university_id = u.id
                LEFT JOIN point_stat ps ON s.id = ps.speciality_id AND u.id = ps.university_id
                LEFT JOIN subject_pair sp ON s.subject_pair_id = sp.id
            WHERE 
                u.id = ?
        )
        SELECT
            speciality_id,
            speciality_name,
            university_name,
            education_format,
            code,
            price,
            degree,
            scholarship,
            avg_salary,
            subject1_id,
            subject2_id,
            MIN(min_score) AS min_score,
            SUM(grant_count) AS grant_count,
            term
        FROM
            speciality_data
        GROUP BY
            speciality_id, speciality_name, university_name, education_format, code, price, degree, scholarship, avg_salary, subject1_id, subject2_id, term
    `

	_, err := o.Raw(query, language, language, language, universityId).QueryRows(&results)
	if err != nil {
		return nil, err
	}

	for i := range results {
		var subjectNames []string

		if results[i].Subject1ID != 0 {
			var subject1 Subject
			err := o.QueryTable("subject").Filter("id", results[i].Subject1ID).One(&subject1)
			if err == nil {
				if language == "ru" {
					subjectNames = append(subjectNames, subject1.NameRu)
				} else {
					subjectNames = append(subjectNames, subject1.NameKz)
				}
			}
		}

		if results[i].Subject2ID != 0 {
			var subject2 Subject
			err := o.QueryTable("subject").Filter("id", results[i].Subject2ID).One(&subject2)
			if err == nil {
				if language == "ru" {
					subjectNames = append(subjectNames, subject2.NameRu)
				} else {
					subjectNames = append(subjectNames, subject2.NameKz)
				}
			}
		}

		results[i].SubjectNames = subjectNames

		annualPoints := []AnnualPoints{}
		annualGrants := []AnnualGrant{}
		var latestGrantCount int

		var pointStats []PointStat
		_, err := o.QueryTable("point_stat").
			Filter("speciality_id", results[i].SpecialityID).
			Filter("university_id", universityId).
			All(&pointStats)
		if err != nil {
			return nil, err
		}

		for _, ps := range pointStats {
			annualPoints = append(annualPoints, AnnualPoints{
				Year:          ps.Year,
				MinScore:      ps.MinScore,
				MinGrantScore: ps.MinGrantScore,
			})
			annualGrants = append(annualGrants, AnnualGrant{
				Year:       ps.Year,
				GrantCount: ps.GrantCount,
			})
			if latestGrantCount == 0 || ps.Year > annualGrants[len(annualGrants)-1].Year {
				latestGrantCount = ps.GrantCount
			}
		}

		results[i].AnnualPoints = annualPoints
		results[i].AnnualGrants = annualGrants
		results[i].GrantCount = latestGrantCount
	}

	if len(results) == 0 {
		return []GetByUniResponseForUser{}, nil
	}

	return results, nil
}

func GetSpecialitiesInUniversityForAdmin(universityID int) ([]IUniverSpecialtyShortcut, error) {
	var specialities []IUniverSpecialtyShortcut

	query := `
        SELECT
            s.id as speciality_id,
            s.name_ru as speciality_name_ru,
            s.name_kz as speciality_name_kz,
            s.code as speciality_code,
            s.degree,
            u.study_format_ru as education_format_ru,
            u.study_format_kz as education_format_kz,
            su.term,
            su.edu_lang,
            COALESCE(json_agg(json_build_object(
                'id', ls.id,
                'grant_count', ls.grant_count,
                'min_score', ls.min_score,
                'year', ls.year,
                'avg_salary', ls.avg_salary,
                'price', ls.price
            )) FILTER (WHERE ls.id IS NOT NULL), '[]') AS statistics
        FROM
            uni_spec.speciality s
            JOIN uni_spec.speciality_university su ON su.speciality_id = s.id
            JOIN uni_spec.university u ON su.university_id = u.id
            LEFT JOIN uni_spec.point_stat ls ON s.id = ls.speciality_id
        WHERE su.university_id = ?
        GROUP BY s.id, s.name_ru, s.name_kz, s.code, s.degree, u.study_format_ru, u.study_format_kz, su.term, su.edu_lang
    `

	// Execute the query
	o := orm.NewOrm()
	rawSeter := o.Raw(query, universityID)
	var maps []orm.Params
	_, err := rawSeter.Values(&maps)
	if err != nil {
		return nil, err
	}

	for _, m := range maps {
		specialityID, _ := strconv.Atoi(m["speciality_id"].(string))
		var statistics []IUniverStatistic
		json.Unmarshal([]byte(m["statistics"].(string)), &statistics)

		term, err := strconv.Atoi(m["term"].(string))
		if err != nil {
			return nil, err
		}

		speciality := IUniverSpecialtyShortcut{
			SpecialityID:      specialityID,
			SpecialityNameRu:  m["speciality_name_ru"].(string),
			SpecialityNameKz:  m["speciality_name_kz"].(string),
			SpecialityCode:    m["speciality_code"].(string),
			Degree:            m["degree"].(string),
			EducationFormatRu: m["education_format_ru"].(string),
			EducationFormatKz: m["education_format_kz"].(string),
			Term:              term,
			EduLang:           m["edu_lang"].(string),
			Statistics:        statistics,
		}

		specialities = append(specialities, speciality)
	}

	return specialities, nil
}

func AssociateSpecialityWithSubjectPair(specialityId int, subjectPairId int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("speciality").Filter("id", specialityId).Update(orm.Params{
		"subject_pair_id": subjectPairId,
	})
	return err
}

func GetSubjectPairsBySpecialityId(specialityId int, language string) ([]*SubjectPair, error) {
	o := orm.NewOrm()
	var specialities []*Speciality

	_, err := o.QueryTable("speciality").Filter("id", specialityId).All(&specialities)
	if err != nil {
		return nil, err
	}

	var subjectPairs []*SubjectPair
	for _, speciality := range specialities {
		if speciality.SubjectPair != nil {
			err := o.Read(speciality.SubjectPair)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
			subjectPairs = append(subjectPairs, speciality.SubjectPair)
		}
	}

	for _, pair := range subjectPairs {
		if pair.Subject1 != nil {
			err := o.Read(pair.Subject1)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
		}
		if pair.Subject2 != nil {
			err := o.Read(pair.Subject2)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
		}
	}

	return subjectPairs, nil
}

func GetSpecialitiesBySubjectPair(subject1Id, subject2Id int, language string) ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality

	_, err := o.Raw(`
		SELECT sp.*
		FROM subject_pair spair
		JOIN speciality sp ON spair.id = sp.subject_pair_id
		WHERE spair.subject1_id = ? AND spair.subject2_id = ?
	`, subject1Id, subject2Id).QueryRows(&specialities)

	if err != nil {
		return nil, err
	}

	// Переназначение полей в зависимости от языка
	for _, speciality := range specialities {
		switch language {
		case "ru":
			speciality.Name = speciality.NameRu
			speciality.Description = speciality.DescriptionRu
		case "kz":
			speciality.Name = speciality.NameKz
			speciality.Description = speciality.DescriptionKz
		}
	}

	return specialities, nil
}

func filterSpecialitiesBySubjectPair(params map[string]interface{}, specialities []*Speciality) ([]*Speciality, error) {
	subject1Id, ok1 := params["subject1_id"].(int)
	subject2Id, ok2 := params["subject2_id"].(int)
	if !ok1 || !ok2 {
		return specialities, nil
	}

	return GetSpecialitiesBySubjectPair(subject1Id, subject2Id, "")
}

func filterSpecialitiesInUniversity(params map[string]interface{}, specialities []*Speciality) ([]*Speciality, error) {
	universityId, ok := params["university_id"].(int)
	if !ok {
		return specialities, nil
	}

	// Получаем язык из параметров, если он существует, иначе используем язык по умолчанию
	language, ok := params["language"].(string)
	if !ok {
		language = "ru" // язык по умолчанию
	}

	// Получаем специальности из университета
	specialityResponses, err := GetSpecialitiesInUniversityForUser(universityId, language)
	if err != nil {
		return nil, err
	}

	// Преобразуем specialityResponses в []*Speciality
	var filteredSpecialities []*Speciality
	for _, sr := range specialityResponses {
		speciality := &Speciality{
			Id:          sr.SpecialityID,
			Name:        sr.SpecialityName,
			Degree:      sr.Degree,
			Scholarship: sr.Scholarship,
		}

		// Динамическое обновление имени в зависимости от языка
		if language == "ru" {
			speciality.NameRu = sr.SpecialityName
		} else if language == "kz" {
			speciality.NameKz = sr.SpecialityName
		}

		filteredSpecialities = append(filteredSpecialities, speciality)
	}

	return filteredSpecialities, nil
}

func filterSpecialitiesByName(params map[string]interface{}, specialities []*Speciality, language string) ([]*Speciality, error) {
	prefix, ok := params["name"].(string)
	if !ok || prefix == "" {
		return specialities, nil
	}

	var results []*Speciality
	o := orm.NewOrm()
	query := "SELECT * FROM speciality WHERE name LIKE ?"
	searchPattern := fmt.Sprintf("%s%%", prefix)
	_, err := o.Raw(query, searchPattern).QueryRows(&results)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[int]*Speciality)
	for _, speciality := range results {
		resultMap[speciality.Id] = speciality
	}

	var filteredSpecialities []*Speciality
	for _, speciality := range specialities {
		if _, found := resultMap[speciality.Id]; found {
			filteredSpecialities = append(filteredSpecialities, speciality)
		}
	}

	// Переназначение полей в зависимости от языка
	for _, speciality := range filteredSpecialities {
		switch language {
		case "ru":
			speciality.Name = speciality.NameRu
			speciality.Description = speciality.DescriptionRu
		case "kz":
			speciality.Name = speciality.NameKz
			speciality.Description = speciality.DescriptionKz
		}
	}

	return filteredSpecialities, nil
}

func paginateSpecialities(specialities []*Speciality, params map[string]interface{}, language string) (*SpecialitySearchResult, error) {
	totalCount := len(specialities)

	page := 1
	if p, ok := params["page"].(int); ok && p > 0 {
		page = p
	}

	perPage := 10
	if pp, ok := params["per_page"].(int); ok && pp > 0 {
		perPage = pp
	}

	totalPages := (totalCount + perPage - 1) / perPage

	start := (page - 1) * perPage
	end := start + perPage

	if start >= totalCount {
		specialities = []*Speciality{}
	} else if end >= totalCount {
		specialities = specialities[start:totalCount]
	} else {
		specialities = specialities[start:end]
	}

	// Переназначение полей в зависимости от языка
	for _, speciality := range specialities {
		switch language {
		case "ru":
			speciality.Name = speciality.NameRu
			speciality.Description = speciality.DescriptionRu
		case "kz":
			speciality.Name = speciality.NameKz
			speciality.Description = speciality.DescriptionKz
		}
	}

	result := &SpecialitySearchResult{
		Specialities: specialities,
		Page:         page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
	}

	return result, nil
}
