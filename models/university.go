package models

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testhub-spec-uni/conf"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/elastic/go-elasticsearch/esapi"
)

type University struct {
	Id               int      `orm:"auto"`
	UniversityCode   string   `orm:"size(64)"`
	Name             string   `orm:"size(128)"`
	Abbreviation     string   `orm:"size(64)"`
	UniversityStatus string   `orm:"size(64)"`
	Address          string   `orm:"size(256)"`
	Website          string   `orm:"size(128)"`
	SocialMediaList  []string `orm:"-"`
	ContactList      []string `orm:"-"`
	AverageFee       int
	HasMilitaryDept  bool
	HasDormitory     bool
	ProfileImageUrl  string `orm:"size(256)"`
	MinEntryScore    int
	PhotosUrlList    []string      `orm:"-"`
	Description      string        `orm:"type(text)"`
	Specialities     []*Speciality `orm:"rel(m2m);rel_table(speciality_university)"`
	City             *City         `orm:"rel(fk)"`
	CreatedAt        time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt        time.Time     `orm:"auto_now;type(datetime)"`
}

type UniversitySearchResponse struct {
	Hits struct {
		Hits []struct {
			Source University `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
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

	err = IndexUniversity(university)
	if err != nil {
		return id, fmt.Errorf("university added but failed to index in Elasticsearch: %v", err)
	}

	return id, nil
}

func GetUniversityById(id int) (*University, error) {
	o := orm.NewOrm()
	university := &University{Id: id}
	err := o.Read(university)
	return university, err
}

func GetAllUniversities() ([]*University, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").All(&universities)
	return universities, err
}

func UpdateUniversity(university *University) error {
	o := orm.NewOrm()
	_, err := o.Update(university)
	return err
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

func IndexUniversity(university *University) error {
	// Convert university to JSON
	data, err := json.Marshal(university)
	if err != nil {
		return err
	}

	// Include speciality IDs in the indexed data
	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return err
	}
	dataMap["speciality_ids"] = getSpecialityIDs(university)

	// Convert back to JSON after updating
	updatedData, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	// Create Elasticsearch index request
	req := esapi.IndexRequest{
		Index:      "universities",
		DocumentID: fmt.Sprintf("%d", university.Id),
		Body:       strings.NewReader(string(updatedData)),
		Refresh:    "true",
	}

	// Execute the request
	res, err := req.Do(context.Background(), conf.EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Handle response
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		}
		return fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	return nil
}

func getSpecialityIDs(university *University) []int {
	var specialityIDs []int
	for _, speciality := range university.Specialities {
		specialityIDs = append(specialityIDs, speciality.Id)
	}
	return specialityIDs
}

func SearchUniversitiesByName(prefix string) ([]University, error) {
	var results []University

	query := fmt.Sprintf(`{
        "query": {
            "query_string": {
                "query": "%s*",
                "fields": ["Name", "Abbreviation", "UniversityCode"]
            }
        }
    }`, prefix)
	res, err := conf.EsClient.Search(
		conf.EsClient.Search.WithContext(context.Background()),
		conf.EsClient.Search.WithIndex("universities"),
		conf.EsClient.Search.WithBody(strings.NewReader(query)),
		conf.EsClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return results, err
		} else {
			return results, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var sr UniversitySearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return results, err
	}

	for _, hit := range sr.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}

func SearchUniversities(params map[string]interface{}) ([]*University, error) {
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
	}

	universities, err = filterByMinScore(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByAvgFee(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByHasMilitaryDept(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByHasDormitory(params, universities)
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
		return filtered, nil
	}
	return universities, nil
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
		return filtered, nil
	}
	return universities, nil
}

func filterByHasMilitaryDept(params map[string]interface{}, universities []*University) ([]*University, error) {
	if hasMilitaryDept, ok := params["has_military_dept"].(bool); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.HasMilitaryDept == hasMilitaryDept {
				filtered = append(filtered, uni)
			}
		}
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
		return filtered, nil
	}
	return universities, nil
}

func filterByCityID(params map[string]interface{}, universities []*University) ([]*University, error) {
	if cityID, ok := params["city_id"].(int); ok {
		var filtered []*University
		for _, uni := range universities {
			if uni.City != nil && uni.City.Id == cityID {
				filtered = append(filtered, uni)
			}
		}
		return filtered, nil
	}
	return universities, nil
}

func filterBySpecialityID(params map[string]interface{}, universities []*University) ([]*University, error) {
	if specialityID, ok := params["speciality_id"].(int); ok {
		var filtered []*University
		for _, uni := range universities {
			for _, spec := range uni.Specialities {
				fmt.Printf("University: %d, Speciality: %d\n", uni.Id, spec.Id) // Debugging line
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

	var firstSubjectId interface{}
	var secondSubjectId interface{}

	if firstSubjectIdInt, ok := params["first_subject_id"].(int); ok {
		firstSubjectId = firstSubjectIdInt
	} else {
		firstSubjectId = nil
	}

	if secondSubjectIdInt, ok := params["second_subject_id"].(int); ok {
		secondSubjectId = secondSubjectIdInt
	} else {
		secondSubjectId = nil
	}

	query := `
        SELECT DISTINCT u.*
        FROM university u
        JOIN speciality_university su ON u.id = su.university_id
        JOIN speciality s ON su.speciality_id = s.id
        JOIN subject_pair sp ON s.subject_pair_id = sp.id
        WHERE 1=1`

	var args []interface{}
	argCount := 1

	if firstSubjectId != nil {
		query += " AND (sp.subject1_id = $" + strconv.Itoa(argCount) + " OR $" + strconv.Itoa(argCount) + " IS NULL)"
		args = append(args, firstSubjectId)
		argCount++
	}

	if secondSubjectId != nil {
		query += " AND (sp.subject2_id = $" + strconv.Itoa(argCount) + " OR $" + strconv.Itoa(argCount) + " IS NULL)"
		args = append(args, secondSubjectId)
		argCount++
	}

	var filtered []*University
	_, err := o.Raw(query, args...).QueryRows(&filtered)
	if err != nil {
		return nil, err
	}

	return filtered, nil
}
