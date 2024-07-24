package models

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type University struct {
	Id                 int      `orm:"auto"`
	UniversityCode     string   `orm:"size(64)"`
	Name               string   `orm:"size(128)"`
	Abbreviation       string   `orm:"size(64)"`
	AbbreviationRu     string   `orm:"size(64)"`
	AbbreviationKz     string   `orm:"size(64)"`
	UniversityStatus   string   `orm:"size(64)"`
	UniversityStatusRu string   `orm:"size(64)"`
	UniversityStatusKz string   `orm:"size(64)"`
	Address            string   `orm:"size(256)"`
	Website            string   `orm:"size(128)"`
	SocialMediaList    []string `orm:"-"`
	ContactList        []string `orm:"-"`
	AverageFee         int
	ProfileImageUrl    string `orm:"size(256)"`
	MinEntryScore      int
	PhotosUrlList      []string      `orm:"-"`
	Description        string        `orm:"type(text)"`
	Specialities       []*Speciality `orm:"rel(m2m);rel_table(speciality_university)"`
	Services           []*Service    `orm:"rel(m2m);rel_table(university_service)"`
	PointStats         []*PointStat  `orm:"reverse(many)"`
	City               *City         `orm:"rel(fk)"`
	CreatedAt          time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt          time.Time     `orm:"auto_now;type(datetime)"`
	CallCenterNumber   string        `orm:"size(64)"`
	WhatsAppNumber     string        `orm:"size(64)"`
	StudyFormat        string        `orm:"size(64)"`
	StudyFormatRu      string        `orm:"size(64)"`
	StudyFormatKz      string        `orm:"size(64)"`
	AddressLink        string        `orm:"size(256)"`
}

type UniversitySearchResult struct {
	Universities []*University `json:"universities"`
	Page         int           `json:"page"`
	TotalPages   int           `json:"total_pages"`
	TotalCount   int           `json:"total_count"`
}

func init() {
	orm.RegisterModel(new(University))
}

func AddUniversity(university *University) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(university)
	if err != nil {
		return 0, err
	}
	university.Id = int(id)

	return id, nil
}

func GetUniversityById(id int) (*University, error) {
	o := orm.NewOrm()
	university := &University{Id: id}
	err := o.Read(university)
	if err != nil {
		return nil, err
	}

	// Load related specialities and services
	if _, err := o.LoadRelated(university, "Specialities"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(university, "Services"); err != nil {
		return nil, err
	}

	return university, nil
}

func GetAllUniversities() ([]*University, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").All(&universities)
	return universities, err
}

func UpdateUniversityFields(university *University) error {
	o := orm.NewOrm()
	existingUniversity := &University{Id: university.Id}

	// Check if the university exists
	if err := o.Read(existingUniversity); err != nil {
		return fmt.Errorf("university with ID %d not found: %v", university.Id, err)
	}

	// Prepare fields to be updated
	updateFields := []string{}

	// Check and add non-zero fields to update list
	if university.UniversityCode != "" {
		updateFields = append(updateFields, "UniversityCode")
	}
	if university.Name != "" {
		updateFields = append(updateFields, "Name")
	}
	if university.Abbreviation != "" {
		updateFields = append(updateFields, "Abbreviation")
	}
	if university.AbbreviationRu != "" {
		updateFields = append(updateFields, "AbbreviationRu")
	}
	if university.AbbreviationKz != "" {
		updateFields = append(updateFields, "AbbreviationKz")
	}
	if university.UniversityStatus != "" {
		updateFields = append(updateFields, "UniversityStatus")
	}
	if university.UniversityStatusRu != "" {
		updateFields = append(updateFields, "UniversityStatusRu")
	}
	if university.UniversityStatusKz != "" {
		updateFields = append(updateFields, "UniversityStatusKz")
	}
	if university.Address != "" {
		updateFields = append(updateFields, "Address")
	}
	if university.Website != "" {
		updateFields = append(updateFields, "Website")
	}
	if university.AverageFee != 0 {
		updateFields = append(updateFields, "AverageFee")
	}
	if university.ProfileImageUrl != "" {
		updateFields = append(updateFields, "ProfileImageUrl")
	}
	if university.MinEntryScore != 0 {
		updateFields = append(updateFields, "MinEntryScore")
	}
	if university.Description != "" {
		updateFields = append(updateFields, "Description")
	}
	if university.CallCenterNumber != "" {
		updateFields = append(updateFields, "CallCenterNumber")
	}
	if university.WhatsAppNumber != "" {
		updateFields = append(updateFields, "WhatsAppNumber")
	}
	if university.StudyFormat != "" {
		updateFields = append(updateFields, "StudyFormat")
	}
	if university.StudyFormatRu != "" {
		updateFields = append(updateFields, "StudyFormatRu")
	}
	if university.StudyFormatKz != "" {
		updateFields = append(updateFields, "StudyFormatKz")
	}
	if university.AddressLink != "" {
		updateFields = append(updateFields, "AddressLink")
	}

	if len(updateFields) > 0 {
		_, err := o.Update(university, updateFields...)
		if err != nil {
			return fmt.Errorf("failed to update university fields: %v", err)
		}
	}

	return nil
}

func DeleteUniversity(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&University{Id: id})
	return err
}

func GetUniversitiesInCity(cityId int) ([]*University, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").
		Filter("City__Id", cityId).
		All(&universities)
	return universities, err
}

func AssignCityToUniversity(universityId int, cityId int) error {
	o := orm.NewOrm()

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}
	city := &City{Id: cityId}
	if err := o.Read(city); err != nil {
		return err
	}

	university.City = city

	if _, err := o.Update(university); err != nil {
		return err
	}

	return nil
}
func AddSpecialityToUniversity(specialityId, universityId int) error {
	o := orm.NewOrm()

	speciality := &Speciality{Id: specialityId}
	if err := o.Read(speciality); err != nil {
		return err
	}

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}

	exist := o.QueryM2M(university, "Specialities").Exist(speciality)
	if exist {
		return fmt.Errorf("speciality with ID %d is already assigned to university with ID %d", specialityId, universityId)
	}

	_, err := o.QueryM2M(university, "Specialities").Add(speciality)
	if err != nil {
		return err
	}

	o.LoadRelated(university, "Specialities")
	fmt.Printf("Specialities for university %d: %v\n", universityId, university.Specialities)

	return nil
}

func getSpecialityIDs(university *University) []int {
	var specialityIDs []int
	for _, speciality := range university.Specialities {
		specialityIDs = append(specialityIDs, speciality.Id)
	}
	return specialityIDs
}

func getServiceIDs(university *University) []int {
	var serviceIDs []int
	for _, service := range university.Services {
		serviceIDs = append(serviceIDs, service.Id)
	}
	return serviceIDs
}

func AddSpecialitiesToUniversity(specialityIds []int, universityId int) error {
	o := orm.NewOrm()

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}

	for _, specialityId := range specialityIds {
		speciality := &Speciality{Id: specialityId}
		if err := o.Read(speciality); err != nil {
			return err
		}

		exist := o.QueryM2M(university, "Specialities").Exist(speciality)
		if exist {
			continue // Skip already assigned specialities
		}

		if _, err := o.QueryM2M(university, "Specialities").Add(speciality); err != nil {
			return err
		}
	}

	o.LoadRelated(university, "Specialities")
	fmt.Printf("Specialities for university %d: %v\n", universityId, university.Specialities)

	return nil
}

func AddServicesToUniversity(serviceIds []int, universityId int) error {
	o := orm.NewOrm()

	university := &University{Id: universityId}
	if err := o.Read(university); err != nil {
		return err
	}

	for _, serviceId := range serviceIds {
		service := &Service{Id: serviceId}
		if err := o.Read(service); err != nil {
			return err
		}

		exist := o.QueryM2M(university, "Services").Exist(service)
		if exist {
			continue // Skip already assigned services
		}

		if _, err := o.QueryM2M(university, "Services").Add(service); err != nil {
			return err
		}
	}

	o.LoadRelated(university, "Services")
	fmt.Printf("Services for university %d: %v\n", universityId, university.Services)

	return nil
}

func SearchUniversities(params map[string]interface{}) (*UniversitySearchResult, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").All(&universities)
	if err != nil {
		return nil, err
	}

	for _, uni := range universities {
		if _, err := o.LoadRelated(uni, "Specialities"); err != nil {
			return nil, err
		}
		if _, err := o.LoadRelated(uni, "Services"); err != nil {
			return nil, err
		}
	}

	universities, err = filterByMinScore(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByAvgFee(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterUniversitiesByName(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterBySubjects(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByCityID(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByStudyFormat(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterBySpecialityIDs(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterBySpecialityID(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterBySortOrder(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByServiceIDs(params, universities)
	if err != nil {
		return nil, err
	}

	result, err := paginateUniversities(universities, params)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SearchUniversities: total universities after filtering: %d\n", len(result.Universities))
	return result, nil
}

func filterUniversitiesByName(params map[string]interface{}, universities []*University) ([]*University, error) {
	prefix, ok := params["name"].(string)
	if !ok || prefix == "" {
		return universities, nil
	}

	o := orm.NewOrm()
	searchPattern := fmt.Sprintf("%%%s%%", prefix)

	var matchingUniversities []*University
	_, err := o.Raw(`
		SELECT * 
		FROM university 
		WHERE name LIKE ? 
		OR abbreviation LIKE ? 
		OR university_code LIKE ?
	`, searchPattern, searchPattern, searchPattern).QueryRows(&matchingUniversities)
	if err != nil {
		return universities, err
	}

	// Create a map for quick lookup
	matchingUniversityMap := make(map[int]*University)
	for _, uni := range matchingUniversities {
		matchingUniversityMap[uni.Id] = uni
	}

	// Filter the original list of universities
	var filteredUniversities []*University
	for _, uni := range universities {
		if _, exists := matchingUniversityMap[uni.Id]; exists {
			filteredUniversities = append(filteredUniversities, uni)
		}
	}

	return filteredUniversities, nil
}

func filterBySortOrder(params map[string]interface{}, universities []*University) ([]*University, error) {
	if sortOrder, ok := params["sort"].(string); ok {
		switch sortOrder {
		case "avg_fee_asc":
			sort.Slice(universities, func(i, j int) bool {
				return universities[i].AverageFee < universities[j].AverageFee
			})
		case "avg_fee_desc":
			sort.Slice(universities, func(i, j int) bool {
				return universities[i].AverageFee > universities[j].AverageFee
			})
		default:
			return universities, fmt.Errorf("invalid sort order: %s", sortOrder)
		}
	}
	return universities, nil
}
func filterByMinScore(params map[string]interface{}, universities []*University) ([]*University, error) {
	if minScore, ok := params["min_score"].(int); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.MinEntryScore >= minScore {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByMinScore: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

func filterByAvgFee(params map[string]interface{}, universities []*University) ([]*University, error) {
	if avgFee, ok := params["avg_fee"].(int); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.AverageFee >= avgFee {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByAvgFee: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

/*
func filterByHasMilitaryDept(params map[string]interface{}, universities []*University) ([]*University, error) {
	if hasMilitaryDept, ok := params["has_military_dept"].(bool); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.HasMilitaryDept == hasMilitaryDept {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByHasMilitaryDept: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

func filterByHasDormitory(params map[string]interface{}, universities []*University) ([]*University, error) {
	if hasDormitory, ok := params["has_dormitory"].(bool); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.HasDormitory == hasDormitory {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByHasDormitory: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

*/

func filterByCityID(params map[string]interface{}, universities []*University) ([]*University, error) {
	if cityID, ok := params["city_id"].(int); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.City != nil && uni.City.Id == cityID {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByCityID: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

func filterBySpecialityIDs(params map[string]interface{}, universities []*University) ([]*University, error) {
	if specialityIDs, ok := params["speciality_ids"].([]int); ok {
		var filtered []*University
		for _, uni := range universities {
			matches := 0
			for _, spec := range uni.Specialities {
				for _, id := range specialityIDs {
					if spec.Id == id {
						matches++
						break
					}
				}
			}
			if matches == len(specialityIDs) {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterBySpecialityIDs: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}
func filterByServiceIDs(params map[string]interface{}, universities []*University) ([]*University, error) {
	if serviceIDs, ok := params["service_ids"].([]int); ok {
		var filtered []*University
		for _, uni := range universities {
			fmt.Println("ID: ", uni.Services)
			matches := 0
			for _, service := range uni.Services {
				for _, id := range serviceIDs {
					if service.Id == id {
						matches++
						break
					}
				}
			}
			if matches == len(serviceIDs) {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByServiceIDs: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

func filterByStudyFormat(params map[string]interface{}, universities []*University) ([]*University, error) {
	if studyFormat, ok := params["study_format"].(string); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.StudyFormat == studyFormat {
				filtered = append(filtered, uni)
			}
		}
		fmt.Printf("filterByStudyFormat: filtered %d universities\n", len(universities)-len(filtered))
		return filtered, nil
	}
	return universities, nil
}

func filterBySpecialityID(params map[string]interface{}, universities []*University) ([]*University, error) {
	if specialityID, ok := params["speciality_id"].(int); ok {
		var filtered []*University
		for _, uni := range universities {
			for _, spec := range uni.Specialities {
				if spec.Id == specialityID {
					filtered = append(filtered, uni)
					break
				}
			}
		}
		return filtered, nil
	}
	return universities, nil
}

func filterBySubjects(params map[string]interface{}, universities []*University) ([]*University, error) {
	o := orm.NewOrm()

	// Check if subject IDs are provided in params
	firstSubjectId, firstOk := params["first_subject_id"].(int)
	secondSubjectId, secondOk := params["second_subject_id"].(int)

	// If both subject IDs are not provided, return original list of universities
	if !firstOk && !secondOk {
		return universities, nil
	}

	// Prepare query and arguments
	query :=
		`SELECT DISTINCT u.*
        FROM university u
        JOIN speciality_university su ON u.id = su.university_id
        JOIN speciality s ON su.speciality_id = s.id
        JOIN subject_pair sp ON s.subject_pair_id = sp.id
        WHERE 1=1`

	var args []interface{}
	argCount := 1

	if firstOk {
		query += " AND (sp.subject1_id = $" + strconv.Itoa(argCount) + " OR $" + strconv.Itoa(argCount) + " IS NULL)"
		args = append(args, firstSubjectId)
		argCount++
	}

	if secondOk {
		query += " AND (sp.subject2_id = $" + strconv.Itoa(argCount) + " OR $" + strconv.Itoa(argCount) + " IS NULL)"
		args = append(args, secondSubjectId)
		argCount++
	}

	var filtered []*University
	_, err := o.Raw(query, args...).QueryRows(&filtered)
	if err != nil {
		return nil, err
	}

	fmt.Printf("filterBySubjects: filtered %d universities\n", len(filtered))
	return filtered, nil
}

func paginateUniversities(universities []*University, params map[string]interface{}) (*UniversitySearchResult, error) {
	totalCount := len(universities)

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
		universities = []*University{}
	} else if end >= totalCount {
		universities = universities[start:totalCount]
	} else {
		universities = universities[start:end]
	}

	result := &UniversitySearchResult{
		Universities: universities,
		Page:         page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
	}

	return result, nil
}
