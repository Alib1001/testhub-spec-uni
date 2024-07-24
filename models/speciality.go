package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Speciality struct {
	Id           int           `orm:"auto"`
	Name         string        `orm:"size(128)"`
	Code         string        `orm:"size(64)"`
	VideoLink    string        `orm:"size(256)"`
	Description  string        `orm:"type(text)"`
	Universities []*University `orm:"reverse(many)"`
	SubjectPair  *SubjectPair  `orm:"rel(fk);on_delete(set_null);null"`
	CreatedAt    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time     `orm:"auto_now;type(datetime)"`
	PointStats   []*PointStat  `orm:"reverse(many)"`
}
type SpecialitySearchResult struct {
	Specialities []*Speciality `json:"specialities"`
	Page         int           `json:"page"`
	TotalPages   int           `json:"total_pages"`
	TotalCount   int           `json:"total_count"`
}

func init() {
	orm.RegisterModel(new(Speciality))
}

func AddSpeciality(speciality *Speciality) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(speciality)

	if err != nil {
		return 0, err
	}

	speciality.Id = int(id)

	return id, err
}

func GetSpecialityById(id int) (*Speciality, error) {
	o := orm.NewOrm()
	speciality := &Speciality{Id: id}
	err := o.QueryTable("speciality").Filter("Id", id).RelatedSel().One(speciality)
	if err != nil {
		return nil, err
	}

	if speciality.SubjectPair != nil {
		err = o.Read(speciality.SubjectPair)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}
	}

	_, err = o.QueryTable("point_stat").Filter("Speciality__Id", id).RelatedSel().All(&speciality.PointStats)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}

	return speciality, nil
}

func GetAllSpecialities() ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	_, err := o.QueryTable("speciality").All(&specialities)
	if err != nil {
		return nil, err
	}

	// Загрузка связанных SubjectPair для каждой специальности
	for _, speciality := range specialities {
		if speciality.SubjectPair != nil {
			err := o.Read(speciality.SubjectPair)
			if err != nil && err != orm.ErrNoRows {
				return nil, err
			}
		}
	}

	return specialities, nil
}

func UpdateSpeciality(speciality *Speciality) error {
	o := orm.NewOrm()
	_, err := o.Update(speciality)
	return err
}

func DeleteSpeciality(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Speciality{Id: id})
	return err
}

func SearchSpecialities(params map[string]interface{}) (*SpecialitySearchResult, error) {
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

	specialities, err = filterSpecialitiesByName(params, specialities)
	if err != nil {
		return nil, err
	}

	// Пагинация
	result, err := paginateSpecialities(specialities, params)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SearchSpecialities: total specialities after filtering: %d\n", len(result.Specialities))
	return result, nil
}

func GetSpecialitiesInUniversity(universityId int) ([]*Speciality, error) {
	o := orm.NewOrm()
	var specialities []*Speciality
	num, err := o.QueryTable("speciality").
		Filter("Universities__University__Id", universityId).
		All(&specialities)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Number of specialities found: %d\n", num) // Добавьте отладочное сообщение
	for _, speciality := range specialities {
		fmt.Printf("Speciality: %+v\n", speciality) // Выводите каждую специальность для отладки
	}
	return specialities, err
}

func AssociateSpecialityWithSubjectPair(specialityId int, subjectPairId int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("speciality").Filter("id", specialityId).Update(orm.Params{
		"subject_pair_id": subjectPairId,
	})
	return err
}

func GetSubjectPairsBySpecialityId(specialityId int) ([]*SubjectPair, error) {
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

func GetSpecialitiesBySubjectPair(subject1Id, subject2Id int) ([]*Speciality, error) {
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

	return specialities, nil
}

func filterSpecialitiesBySubjectPair(params map[string]interface{}, specialities []*Speciality) ([]*Speciality, error) {
	subject1Id, ok1 := params["subject1_id"].(int)
	subject2Id, ok2 := params["subject2_id"].(int)
	if !ok1 || !ok2 {
		return specialities, nil
	}

	return GetSpecialitiesBySubjectPair(subject1Id, subject2Id)
}

func filterSpecialitiesInUniversity(params map[string]interface{}, specialities []*Speciality) ([]*Speciality, error) {
	universityId, ok := params["university_id"].(int)
	if !ok {
		return specialities, nil
	}

	return GetSpecialitiesInUniversity(universityId)
}

func filterSpecialitiesByName(params map[string]interface{}, specialities []*Speciality) ([]*Speciality, error) {
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

	return filteredSpecialities, nil
}

func paginateSpecialities(specialities []*Speciality, params map[string]interface{}) (*SpecialitySearchResult, error) {
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

	result := &SpecialitySearchResult{
		Specialities: specialities,
		Page:         page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
	}

	return result, nil
}
