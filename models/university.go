package models

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type University struct {
	Id                 int      `orm:"auto"`
	UniversityCode     string   `orm:"size(64)"`
	Name               string   `orm:"size(128)"`
	NameRu             string   `orm:"size(128)" json:"-"`
	NameKz             string   `orm:"size(128)" json:"-"`
	Abbreviation       string   `orm:"size(64)"`
	AbbreviationRu     string   `orm:"size(64)" json:"-"`
	AbbreviationKz     string   `orm:"size(64)" json:"-"`
	UniversityStatus   string   `orm:"size(64)" json:"-"`
	UniversityStatusRu string   `orm:"size(64)" json:"-"`
	UniversityStatusKz string   `orm:"size(64)" json:"-"`
	Address            string   `orm:"size(256)"`
	Website            string   `orm:"size(128)"`
	SocialMediaList    []string `orm:"-"`
	ContactList        []string `orm:"-"`
	AverageFee         int
	MainImageUrl       string `orm:"size(256)"`
	MinEntryScore      int
	PhotosUrlList      []string                `orm:"-"`
	Description        string                  `orm:"type(text)"`
	DescriptionRu      string                  `orm:"type(text)" json:"-"`
	DescriptionKz      string                  `orm:"type(text)" json:"-"`
	Specialities       []*SpecialityUniversity `orm:"reverse(many);on_delete(cascade)"`
	Services           []*Service              `orm:"rel(m2m);rel_table(university_service);on_delete(cascade)"`
	PointStats         []*PointStat            `orm:"reverse(many);on_delete(cascade)"`
	City               *City                   `orm:"rel(fk)"`
	CreatedAt          time.Time               `orm:"auto_now_add;type(datetime)"`
	UpdatedAt          time.Time               `orm:"auto_now;type(datetime)"`
	CallCenterNumber   string                  `orm:"size(64)"`
	WhatsAppNumber     string                  `orm:"size(64)"`
	StudyFormat        string                  `orm:"size(64)"`
	StudyFormatRu      string                  `orm:"size(64)"`
	StudyFormatKz      string                  `orm:"size(64)"`
	AddressLink        string                  `orm:"size(256)"`
	Email              string                  `orm:"size(64)"`
	Rating             string                  `orm:"size(64)"`
	Gallery            []*Gallery              `orm:"reverse(many);on_delete(cascade)"`
	Popular            bool
}

type UniversitySearchResult struct {
	Universities []*GetAllUniversityResponse `json:"universities"`
	Page         int                         `json:"page"`
	TotalPages   int                         `json:"total_pages"`
	TotalCount   int                         `json:"total_count"`
}

type GetAllUniversityResponse struct {
	Id               int    `json:"Id"`
	Name             string `json:"Name"`
	ImageUrl         string `json:"ImageUrl"`
	Address          string `json:"Address"`
	UniversityCode   string `json:"UniversityCode"`
	SpecialityCount  int    `json:"SpecialityCount"`
	UniversityStatus string `json:"UniversityStatus"`
	MinScore         int    `json:"MinScore"`
	Rating           string `json:"Rating"`
}
type GetUniNamesResponse struct {
	Id           int    `json:"Id"`
	Abbreviation string `json:"Abbreviation"`
}

type GetAllUniversityForAdminResponse struct {
	Id                 int    `json:"Id"`
	NameRu             string `json:"NameRu"`
	NameKz             string `json:"NameKz"`
	ImageUrl           string `json:"ImageUrl"`
	Address            string `json:"Address"`
	UniversityCode     string `json:"UniversityCode"`
	SpecialityCount    int    `json:"SpecialityCount"`
	UniversityStatusRu string `json:"UniversityStatusRu"`
	UniversityStatusKz string `json:"UniversityStatusKz"`
	MinScore           int    `json:"MinScore"`
	Rating             string `json:"Rating"`
}
type GetByIdUniversityResponseForAdmin struct {
	Id                 int                `json:"Id"`
	NameRu             string             `json:"NameRu" validate:"required"`
	NameKz             string             `json:"NameKz" validate:"required"`
	UniversityStatusRu string             `json:"UniversityStatusRu" validate:"required"`
	UniversityStatusKz string             `json:"UniversityStatusKz" validate:"required"`
	Website            string             `json:"Website" validate:"required,url"`
	Email              string             `json:"Email"`
	CallCenterNumber   string             `json:"CallCenterNumber" validate:"required"`
	WhatsAppNumber     string             `json:"WhatsAppNumber" validate:"required"`
	Address            string             `json:"Address" validate:"required"`
	UniversityCode     string             `json:"UniversityCode" validate:"required"`
	StudyFormatRu      string             `json:"StudyFormatRu" validate:"required"`
	StudyFormatKz      string             `json:"StudyFormatKz" validate:"required"`
	AbbreviationRu     string             `json:"AbbreviationRu" validate:"required"`
	AbbreviationKz     string             `json:"AbbreviationKz" validate:"required"`
	MainImageUrl       string             `json:"MainImageUrl" validate:"required,url"`
	AddressLink        string             `json:"AddressLink" validate:"required"`
	DescriptionRu      string             `json:"DescriptionRu" validate:"required"`
	DescriptionKz      string             `json:"DescriptionKz" validate:"required"`
	Rating             string             `json:"Rating" validate:"required"`
	MinScore           int                `json:"MinScore" `
	Gallery            []*GalleryResponse `json:"Gallery"`
	Services           []*Service         `json:"Services"`
	City               *City              `json:"City" validate:"required"`
}

type GetByIdUniversityResponseForUser struct {
	Id               int                       `json:"Id"`
	Name             string                    `json:"Name"`
	Abbreviation     string                    `json:"Abbreviation"`
	MainImageUrl     string                    `json:"MainImageUrl"`
	Website          string                    `json:"Website"`
	Email            string                    `json:"Email"`
	CallCenterNumber string                    `json:"CallCenterNumber"`
	WhatsAppNumber   string                    `json:"WhatsAppNumber"`
	Address          string                    `json:"Address"`
	AddressLink      string                    `json:"AddressLink"`
	Description      string                    `json:"Description"`
	Services         []*ServiceResponseForUser `json:"Services"`
	Gallery          []*GalleryResponse        `json:"Gallery"`
}
type AddUUniversityResponse struct {
	Id                 int      `form:"Id"`
	NameRu             string   `form:"NameRu" validate:"required"`
	NameKz             string   `form:"NameKz" validate:"required"`
	UniversityStatusRu string   `form:"UniversityStatusRu" validate:"required"`
	UniversityStatusKz string   `form:"UniversityStatusKz" validate:"required"`
	Website            string   `form:"Website" validate:"required,url"`
	CallCenterNumber   string   `form:"CallCenterNumber" validate:"required"`
	WhatsAppNumber     string   `form:"WhatsAppNumber" validate:"required"`
	Email              string   `form:"Email"`
	Address            string   `form:"Address" validate:"required"`
	UniversityCode     string   `form:"UniversityCode" validate:"required"`
	StudyFormatRu      string   `form:"StudyFormatRu" validate:"required"`
	StudyFormatKz      string   `form:"StudyFormatKz" validate:"required"`
	AbbreviationRu     string   `form:"AbbreviationRu" validate:"required"`
	AbbreviationKz     string   `form:"AbbreviationKz" validate:"required"`
	MainImageUrl       string   `form:"MainImageUrl"`
	AddressLink        string   `form:"AddressLink" validate:"required"`
	DescriptionRu      string   `form:"DescriptionRu" validate:"required"`
	DescriptionKz      string   `form:"DescriptionKz" validate:"required"`
	Rating             string   `form:"Rating" validate:"required"`
	MinScore           int      `form:"MinScore" validate:"required"`
	Gallery            []string `form:"Gallery"`
	ServiceIds         []int    `form:"ServiceIds"`
	CityId             int      `form:"CityId" validate:"required"`
}

type AddUniversityPartial struct {
	NameRu             string   `form:"NameRu" validate:"required"`
	NameKz             string   `form:"NameKz" validate:"required"`
	UniversityStatusRu string   `form:"UniversityStatusRu" validate:"required"`
	UniversityStatusKz string   `form:"UniversityStatusKz" validate:"required"`
	Website            string   `form:"Website" validate:"required,url"`
	Email              string   `form:"Email"`
	CallCenterNumber   string   `form:"CallCenterNumber" validate:"required"`
	WhatsAppNumber     string   `form:"WhatsAppNumber" validate:"required"`
	Address            string   `form:"Address" validate:"required"`
	UniversityCode     string   `form:"UniversityCode" validate:"required"`
	StudyFormatRu      string   `form:"StudyFormatRu" validate:"required"`
	StudyFormatKz      string   `form:"StudyFormatKz" validate:"required"`
	AbbreviationRu     string   `form:"AbbreviationRu" validate:"required"`
	AbbreviationKz     string   `form:"AbbreviationKz" validate:"required"`
	MainImageUrl       string   `form:"MainImageUrl"`
	AddressLink        string   `form:"AddressLink" validate:"required"`
	DescriptionRu      string   `form:"DescriptionRu" validate:"required"`
	DescriptionKz      string   `form:"DescriptionKz" validate:"required"`
	Rating             string   `form:"Rating" validate:"required"`
	MinScore           int      `form:"MinScore" validate:"required"`
	Gallery            []string `form:"Gallery"`
	CityId             int      `form:"CityId" validate:"required"`
}

type UpdateUniversityResponse struct {
	Id                 int      `form:"Id"`
	NameRu             string   `form:"NameRu"`
	NameKz             string   `form:"NameKz"`
	UniversityStatusRu string   `form:"UniversityStatusRu"`
	UniversityStatusKz string   `form:"UniversityStatusKz"`
	Website            string   `form:"Website"`
	CallCenterNumber   string   `form:"CallCenterNumber"`
	WhatsAppNumber     string   `form:"WhatsAppNumber"`
	Address            string   `form:"Address"`
	UniversityCode     string   `form:"UniversityCode"`
	StudyFormatRu      string   `form:"StudyFormatRu"`
	StudyFormatKz      string   `form:"StudyFormatKz"`
	AbbreviationRu     string   `form:"AbbreviationRu"`
	AbbreviationKz     string   `form:"AbbreviationKz"`
	MainImageUrl       string   `form:"MainImageUrl"`
	Email              string   `form:"Email"`
	AddressLink        string   `form:"AddressLink"`
	DescriptionRu      string   `form:"DescriptionRu"`
	DescriptionKz      string   `form:"DescriptionKz"`
	Rating             string   `form:"Rating"`
	MinScore           int      `form:"MinScore"`
	Gallery            []string `form:"Gallery"`
	ServiceIds         []int    `form:"ServiceIds"`
	CityId             int      `form:"CityId"`
}

type UpdateUniversityPartial struct {
	Id                 int                     `form:"Id"`
	NameRu             string                  `form:"NameRu"`
	NameKz             string                  `form:"NameKz"`
	UniversityStatusRu string                  `form:"UniversityStatusRu"`
	UniversityStatusKz string                  `form:"UniversityStatusKz"`
	Website            string                  `form:"Website"`
	CallCenterNumber   string                  `form:"CallCenterNumber"`
	WhatsAppNumber     string                  `form:"WhatsAppNumber"`
	Address            string                  `form:"Address"`
	UniversityCode     string                  `form:"UniversityCode"`
	StudyFormatRu      string                  `form:"StudyFormatRu"`
	StudyFormatKz      string                  `form:"StudyFormatKz"`
	AbbreviationRu     string                  `form:"AbbreviationRu"`
	AbbreviationKz     string                  `form:"AbbreviationKz"`
	MainImageUrl       string                  `form:"MainImageUrl"`
	Email              string                  `form:"Email"`
	AddressLink        string                  `form:"AddressLink"`
	DescriptionRu      string                  `form:"DescriptionRu"`
	DescriptionKz      string                  `form:"DescriptionKz"`
	Rating             string                  `form:"Rating"`
	MinScore           int                     `form:"MinScore"`
	Gallery            []string                `form:"Gallery"`
	SpecialityTerms    []*SpecialityUniversity `orm:"reverse(many)" json:"speciality_terms,omitempty"`
	CityId             int                     `form:"CityId"`
}

type LocalizedFields struct {
	Name   string
	Status string
}

func init() {
	orm.RegisterModel(new(University))
}
func AddUniversity(universityResponse *AddUUniversityResponse) (int64, error) {
	validate := validator.New()

	err := validate.Struct(universityResponse)
	if err != nil {
		return 0, err
	}

	var city City
	o := orm.NewOrm()
	err = o.QueryTable("city").Filter("Id", universityResponse.CityId).One(&city)
	if err != nil {
		return 0, err
	}

	dbUniversity := &University{
		NameRu:             universityResponse.NameRu,
		NameKz:             universityResponse.NameKz,
		UniversityStatusRu: universityResponse.UniversityStatusRu,
		UniversityStatusKz: universityResponse.UniversityStatusKz,
		Website:            universityResponse.Website,
		Email:              universityResponse.Email,
		CallCenterNumber:   universityResponse.CallCenterNumber,
		WhatsAppNumber:     universityResponse.WhatsAppNumber,
		Address:            universityResponse.Address,
		UniversityCode:     universityResponse.UniversityCode,
		StudyFormatRu:      universityResponse.StudyFormatRu,
		StudyFormatKz:      universityResponse.StudyFormatKz,
		AbbreviationRu:     universityResponse.AbbreviationRu,
		AbbreviationKz:     universityResponse.AbbreviationKz,
		MainImageUrl:       universityResponse.MainImageUrl,
		AddressLink:        universityResponse.AddressLink,
		DescriptionRu:      universityResponse.DescriptionRu,
		DescriptionKz:      universityResponse.DescriptionKz,
		Rating:             universityResponse.Rating,
		MinEntryScore:      universityResponse.MinScore,
		City:               &city,
	}

	id, err := o.Insert(dbUniversity)
	if err != nil {
		return 0, err
	}
	universityResponse.Id = int(id)

	for _, galleryURL := range universityResponse.Gallery {
		gallery := &Gallery{
			University: dbUniversity,
			PhotoUrl:   galleryURL,
		}
		_, err := o.Insert(gallery)
		if err != nil {
			return 0, err
		}
	}

	for _, serviceID := range universityResponse.ServiceIds {
		service := &Service{Id: serviceID}
		m2m := o.QueryM2M(dbUniversity, "Services")
		_, err := m2m.Add(service)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func UpdateUniversityImageURL(id int64, imageURL string) error {
	o := orm.NewOrm()

	// Create an instance of University with the provided ID
	university := University{Id: int(id)}

	// Read the university record from the database
	if err := o.Read(&university); err != nil {
		if err == orm.ErrNoRows {
			return errors.New("university not found")
		}
		fmt.Printf("Error reading university: %v\n", err)
		return fmt.Errorf("error reading university: %v", err)
	}

	// Update the MainImageUrl field
	university.MainImageUrl = imageURL
	if _, err := o.Update(&university, "MainImageUrl"); err != nil {
		// Log the error for debugging purposes
		fmt.Printf("Error updating university image URL: %v\n", err)
		return fmt.Errorf("error updating university image URL: %v", err)
	}

	return nil
}

func AddGalleryImages(universityID int64, galleryURLs []string) error {
	o := orm.NewOrm()
	university := &University{Id: int(universityID)}

	if err := o.Read(university); err != nil {
		if err == orm.ErrNoRows {
			return errors.New("university not found")
		}
		return err
	}

	for _, galleryURL := range galleryURLs {
		gallery := &Gallery{
			University: university,
			PhotoUrl:   galleryURL,
		}
		if _, err := o.Insert(gallery); err != nil {
			return err
		}
	}

	return nil
}

func GetUniversityByIdForAdmin(id int) (*GetByIdUniversityResponseForAdmin, error) {
	o := orm.NewOrm()
	university := &University{Id: id}
	err := o.Read(university)
	if err != nil {
		return nil, err
	}

	if _, err := o.LoadRelated(university, "Services"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(university, "Gallery"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(university, "City"); err != nil {
		return nil, err
	}

	var galleryResponses []*GalleryResponse
	for _, gallery := range university.Gallery {
		galleryResponses = append(galleryResponses, &GalleryResponse{
			Id:       gallery.Id,
			PhotoUrl: gallery.PhotoUrl,
		})
	}

	response := &GetByIdUniversityResponseForAdmin{
		Id:                 university.Id,
		NameRu:             university.NameRu,
		NameKz:             university.NameKz,
		Email:              university.Email,
		Address:            university.Address,
		UniversityStatusRu: university.UniversityStatusRu,
		UniversityStatusKz: university.UniversityStatusKz,
		Website:            university.Website,
		CallCenterNumber:   university.CallCenterNumber,
		WhatsAppNumber:     university.WhatsAppNumber,
		UniversityCode:     university.UniversityCode,
		StudyFormatRu:      university.StudyFormatRu,
		StudyFormatKz:      university.StudyFormatKz,
		AbbreviationRu:     university.AbbreviationRu,
		AbbreviationKz:     university.AbbreviationKz,
		MainImageUrl:       university.MainImageUrl,
		Gallery:            galleryResponses,
		AddressLink:        university.AddressLink,
		DescriptionRu:      university.DescriptionRu,
		DescriptionKz:      university.DescriptionKz,
		Rating:             university.Rating,
		MinScore:           university.MinEntryScore,
		Services:           university.Services,
		City:               university.City,
	}

	return response, nil
}

func GetUniversityByIdForUser(id int, language string) (*GetByIdUniversityResponseForUser, error) {
	o := orm.NewOrm()
	university := &University{Id: id}
	err := o.Read(university)
	if err != nil {
		return nil, err
	}

	if _, err := o.LoadRelated(university, "Services"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(university, "Gallery"); err != nil {
		return nil, err
	}

	var galleryResponses []*GalleryResponse
	for _, gallery := range university.Gallery {
		galleryResponses = append(galleryResponses, &GalleryResponse{
			Id:       gallery.Id,
			PhotoUrl: gallery.PhotoUrl,
		})
	}

	var name, abbreviation, description string
	switch language {
	case "ru":
		name = university.NameRu
		description = university.DescriptionRu
		abbreviation = university.AbbreviationRu
	case "kz":
		name = university.NameKz
		description = university.DescriptionKz
		abbreviation = university.AbbreviationKz
	default:
		name = university.NameKz
		description = university.DescriptionKz
		abbreviation = university.AbbreviationKz
	}

	serviceResponses := []*ServiceResponseForUser{}
	for _, service := range university.Services {
		serviceResponse, err := GetServiceById(service.Id, language)
		if err != nil {
			return nil, err
		}
		serviceResponses = append(serviceResponses, serviceResponse)
	}

	response := &GetByIdUniversityResponseForUser{
		Id:               university.Id,
		Name:             name,
		Website:          university.Website,
		Email:            university.Email,
		CallCenterNumber: university.CallCenterNumber,
		WhatsAppNumber:   university.WhatsAppNumber,
		Address:          university.Address,
		Abbreviation:     abbreviation,
		MainImageUrl:     university.MainImageUrl,
		AddressLink:      university.AddressLink,
		Description:      description,
		Gallery:          galleryResponses,
		Services:         serviceResponses,
	}

	return response, nil
}

func GetAllUniversities(language string, page int, perPage int) ([]*GetAllUniversityResponse, int64, int, int, error) {
	o := orm.NewOrm()
	var universities []*University

	offset := (page - 1) * perPage

	totalCount, err := o.QueryTable("university").Count()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	totalPage := int((totalCount + int64(perPage) - 1) / int64(perPage))

	_, err = o.QueryTable("university").Limit(perPage).Offset(offset).All(&universities)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	var responses []*GetAllUniversityResponse
	for _, university := range universities {
		switch language {
		case "ru":
			university.Name = university.NameRu
			university.Description = university.DescriptionRu
			university.UniversityStatus = university.UniversityStatusRu
		case "kz":
			university.Name = university.NameKz
			university.Description = university.DescriptionKz
			university.UniversityStatus = university.UniversityStatusKz
		default:
			university.Name = university.NameKz
			university.Description = university.DescriptionKz
		}

		if _, err := o.LoadRelated(university, "Specialities"); err != nil {
			return nil, 0, 0, 0, err
		}

		response := &GetAllUniversityResponse{
			Id:               university.Id,
			Name:             university.Name,
			ImageUrl:         university.MainImageUrl,
			Address:          university.Address,
			UniversityCode:   university.UniversityCode,
			SpecialityCount:  len(university.Specialities),
			UniversityStatus: university.UniversityStatus,
			MinScore:         university.MinEntryScore,
			Rating:           university.Rating,
		}

		responses = append(responses, response)
	}

	return responses, totalCount, totalPage, page, nil
}

func GetAllUniversitiesForAdmin() ([]*GetAllUniversityForAdminResponse, error) {
	o := orm.NewOrm()
	var universities []*University
	_, err := o.QueryTable("university").All(&universities)
	if err != nil {
		return nil, err
	}

	var responses []*GetAllUniversityForAdminResponse
	for _, university := range universities {
		if _, err := o.LoadRelated(university, "Specialities"); err != nil {
			return nil, err
		}

		response := &GetAllUniversityForAdminResponse{
			Id:                 university.Id,
			NameRu:             university.NameRu,
			NameKz:             university.NameKz,
			ImageUrl:           university.MainImageUrl,
			Address:            university.Address,
			UniversityCode:     university.UniversityCode,
			SpecialityCount:    len(university.Specialities),
			UniversityStatusRu: university.UniversityStatusRu,
			UniversityStatusKz: university.UniversityStatusKz,
			MinScore:           university.MinEntryScore,
			Rating:             university.Rating,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func GetUniversityByID(id int) (*University, error) {
	var university University
	if err := orm.NewOrm().QueryTable("university").Filter("id", id).RelatedSel().One(&university); err != nil {
		return nil, err
	}
	return &university, nil
}

func UpdateUniversity(university *University) error {
	_, err := orm.NewOrm().Update(university)
	return err
}

func UpdateUniversityGallery(universityID int, newGalleryURLs []string) error {
	o := orm.NewOrm()

	// Fetch existing gallery images for the university
	var existingGalleries []*Gallery
	_, err := o.QueryTable("gallery").Filter("university_id", universityID).All(&existingGalleries)
	if err != nil {
		return err
	}

	// Create a map of existing URLs for quick lookup
	existingURLs := make(map[string]bool)
	for _, gallery := range existingGalleries {
		existingURLs[gallery.PhotoUrl] = true
	}

	// Start a transaction to ensure atomicity
	err = o.Begin()
	if err != nil {
		return err
	}

	// Insert new gallery images that do not already exist
	for _, url := range newGalleryURLs {
		if !existingURLs[url] {
			gallery := &Gallery{
				University: &University{Id: universityID},
				PhotoUrl:   url,
			}
			if _, err := o.Insert(gallery); err != nil {
				o.Rollback() // Rollback transaction on error
				return err
			}
		}
	}

	// Remove images that are no longer in the new gallery URLs
	for _, gallery := range existingGalleries {
		if !contains(newGalleryURLs, gallery.PhotoUrl) {
			if _, err := o.Delete(gallery); err != nil {
				o.Rollback() // Rollback transaction on error
				return err
			}
		}
	}

	// Commit the transaction if everything is successful
	if err := o.Commit(); err != nil {
		return err
	}

	return nil
}

// Helper function to check if a slice contains a specific element
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func DeleteUniversity(id int) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		fmt.Println("Failed to start transaction:", err)
		return err
	}

	university := University{Id: id}
	if err := o.Read(&university); err != nil {
		o.Rollback()
		fmt.Println("Failed to read university:", err)
		return err
	}

	if university.MainImageUrl != "" {
		fmt.Println("Deleting main image from cloud:", university.MainImageUrl)
		err = deleteFileFromCloud(university.MainImageUrl)
		if err != nil {
			o.Rollback()
			fmt.Println("Failed to delete main image from cloud:", err)
			return err
		}
	}

	var galleries []*Gallery
	_, err = o.QueryTable("gallery").Filter("university_id", id).All(&galleries)
	if err != nil {
		o.Rollback() // Rollback transaction on error
		fmt.Println("Failed to fetch galleries:", err)
		return err
	}

	for _, gallery := range galleries {
		if gallery.PhotoUrl != "" {
			fmt.Println("Deleting gallery image from cloud:", gallery.PhotoUrl)
			err = deleteFileFromCloud(gallery.PhotoUrl)
			if err != nil {
				o.Rollback() // Rollback transaction on error
				fmt.Println("Failed to delete gallery image from cloud:", err)
				return err
			}
		}
	}

	// Delete the university
	_, err = o.Delete(&university)
	if err != nil {
		o.Rollback() // Rollback transaction on error
		fmt.Println("Failed to delete university:", err)
		return err
	}

	err = o.Commit() // Commit transaction
	if err != nil {
		fmt.Println("Failed to commit transaction:", err)
	}
	return err
}

func deleteFileFromCloud(filePath string) error {
	awsAccessKey, _ := beego.AppConfig.String("aws_access_key")
	awsSecretKey, _ := beego.AppConfig.String("aws_secret_key")
	bucket, _ := beego.AppConfig.String("bucket")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		Endpoint:         aws.String("https://chi-sextans.object.pscloud.io"),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to wait for file deletion: %v", err)
	}

	return nil
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
		specialityUniversity := &SpecialityUniversity{
			Speciality: &Speciality{Id: specialityId},
			University: university,
			Term:       4,
			EduLang:    "kz",
		}

		// Check if the relation already exists
		exist := o.QueryTable("speciality_university").Filter("speciality_id", specialityId).Filter("university_id", universityId).Exist()
		if exist {
			continue
		}

		// Insert the new SpecialityUniversity entry
		if _, err := o.Insert(specialityUniversity); err != nil {
			return err
		}
	}

	// Load the related Specialities
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

func SearchUniversities(params map[string]interface{}, language string) (*UniversitySearchResult, error) {
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

	universities, err = filterUniversitiesByName(params, universities, language)
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

	universities, err = filterByTerm(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterBySpecialityID(params, universities)
	if err != nil {
		return nil, err
	}

	universities, err = filterByStatus(params, universities, language)
	if err != nil {
		return nil, err
	}

	universities, err = filterBySortOrder(params, universities, language)
	if err != nil {
		return nil, err
	}

	universities, err = filterByServiceIDs(params, universities)
	if err != nil {
		return nil, err
	}

	var universityResponses []*GetAllUniversityResponse
	for _, uni := range universities {
		localized := selectLanguage(uni, language)
		uniResponse := &GetAllUniversityResponse{
			Id:               uni.Id,
			Name:             localized.Name,
			UniversityStatus: localized.Status,
			ImageUrl:         uni.MainImageUrl,
			Address:          uni.Address,
			UniversityCode:   uni.UniversityCode,
			SpecialityCount:  len(uni.Specialities),
			MinScore:         uni.MinEntryScore,
			Rating:           uni.Rating,
		}
		universityResponses = append(universityResponses, uniResponse)
	}

	result, err := paginateUniversities(universityResponses, params)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SearchUniversities: total universities after filtering: %d\n", len(result.Universities))
	return result, nil
}

func selectLanguage(uni *University, language string) LocalizedFields {
	switch language {
	case "ru":
		return LocalizedFields{
			Name:   uni.NameRu,
			Status: uni.UniversityStatusRu,
		}
	case "kz":
		return LocalizedFields{
			Name:   uni.NameKz,
			Status: uni.UniversityStatusKz,
		}
	default:
		return LocalizedFields{
			Name:   uni.UniversityStatusKz,
			Status: uni.UniversityStatusKz,
		}
	}
}
func filterUniversitiesByName(params map[string]interface{}, universities []*University, lang string) ([]*University, error) {
	prefix, ok := params["name"].(string)
	if !ok || prefix == "" {
		return universities, nil
	}

	if lang != "ru" && lang != "kz" {
		return universities, fmt.Errorf("invalid or missing lang parameter")
	}

	o := orm.NewOrm()
	searchPattern := fmt.Sprintf("%%%s%%", prefix)

	var nameField, abbreviationField string
	if lang == "ru" {
		nameField = "name_ru"
		abbreviationField = "abbreviation_ru"
	} else if lang == "kz" {
		nameField = "name_kz"
		abbreviationField = "abbreviation_kz"
	}

	var matchingUniversities []*University
	query := fmt.Sprintf(`
		SELECT * 
		FROM university 
		WHERE %s ILIKE ? 
		OR %s ILIKE ? 
		OR university_code ILIKE ?
	`, nameField, abbreviationField)

	_, err := o.Raw(query, searchPattern, searchPattern, searchPattern).QueryRows(&matchingUniversities)
	if err != nil {
		return universities, err
	}

	matchingUniversityMap := make(map[int]*University)
	for _, uni := range matchingUniversities {
		matchingUniversityMap[uni.Id] = uni
	}

	var filteredUniversities []*University
	for _, uni := range universities {
		if _, exists := matchingUniversityMap[uni.Id]; exists {
			filteredUniversities = append(filteredUniversities, uni)
		}
	}

	return filteredUniversities, nil
}

func filterBySortOrder(params map[string]interface{}, universities []*University, lang string) ([]*University, error) {
	if sortOrder, ok := params["sort"].(string); ok {
		switch sortOrder {
		case "name_asc":
			if lang == "ru" {
				sort.Slice(universities, func(i, j int) bool {
					return universities[i].NameRu < universities[j].NameRu
				})
			} else if lang == "kz" {
				sort.Slice(universities, func(i, j int) bool {
					return universities[i].NameKz < universities[j].NameKz
				})
			} else {
				return universities, fmt.Errorf("unsupported language: %s", lang)
			}
		case "name_desc":
			if lang == "ru" {
				sort.Slice(universities, func(i, j int) bool {
					return universities[i].NameRu > universities[j].NameRu
				})
			} else if lang == "kz" {
				sort.Slice(universities, func(i, j int) bool {
					return universities[i].NameKz > universities[j].NameKz
				})
			} else {
				return universities, fmt.Errorf("unsupported language: %s", lang)
			}
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

func filterByTerm(params map[string]interface{}, universities []*University) ([]*University, error) {
	if term, ok := params["term"].(int); ok {
		var filteredUniversities []*University

		fmt.Printf("Filtering by Term: %d\n", term)

		for _, uni := range universities {
			fmt.Printf("Checking university: %s (ID: %d)\n", uni.Name, uni.Id)

			matchingTermFound := false
			for _, specUni := range uni.Specialities {
				if specUni.Term == term {
					fmt.Printf(" -> Match found for Term: %d in university: %s (ID: %d)\n", term, uni.Name, uni.Id)
					matchingTermFound = true
					break
				}
			}

			if matchingTermFound {
				filteredUniversities = append(filteredUniversities, uni)
			}
		}

		fmt.Printf("Total universities after term filter: %d\n", len(filteredUniversities))
		return filteredUniversities, nil
	}

	// If no term is provided, return the original list of universities
	return universities, nil
}

func filterByStatus(params map[string]interface{}, universities []*University, language string) ([]*University, error) {
	statusFilter, ok := params["status"].(string)
	if !ok || statusFilter == "" {
		return universities, nil
	}
	var filteredUniversities []*University
	statusFilterLower := strings.ToLower(statusFilter)
	for _, uni := range universities {
		var universityStatus string
		switch language {
		case "ru":
			universityStatus = uni.UniversityStatusRu
		case "kz":
			universityStatus = uni.UniversityStatusKz
		default:
			universityStatus = uni.UniversityStatusKz
		}
		if strings.ToLower(universityStatus) == statusFilterLower {
			filteredUniversities = append(filteredUniversities, uni)
		}
	}

	return filteredUniversities, nil
}

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
	if specialityIDs, ok := params["speciality_ids"].([]int); ok && len(specialityIDs) > 0 {
		var filteredUniversities []*University

		fmt.Printf("Filtering by Speciality IDs: %v\n", specialityIDs)

		for _, uni := range universities {
			fmt.Printf("Checking university: %s (ID: %d)\n", uni.Name, uni.Id)

			specialityMap := make(map[int]bool)
			for _, specUni := range uni.Specialities {
				if containsInt(specialityIDs, specUni.Speciality.Id) {
					specialityMap[specUni.Speciality.Id] = true
				}
			}

			allSpecialitiesMatched := true
			for _, id := range specialityIDs {
				if !specialityMap[id] {
					allSpecialitiesMatched = false
					break
				}
			}

			if allSpecialitiesMatched {
				filteredUniversities = append(filteredUniversities, uni)
			}
		}

		return filteredUniversities, nil
	}

	return universities, nil
}
func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
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
			fmt.Println("studyFormat", studyFormat)
			fmt.Println("studyFormatRu", uni.StudyFormatRu)
			fmt.Println("studyFormatKz", uni.StudyFormatKz)
			if uni.StudyFormat == studyFormat || uni.StudyFormatRu == studyFormat || uni.StudyFormatKz == studyFormat {
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

	firstSubjectId, firstOk := params["first_subject_id"].(int)
	secondSubjectId, secondOk := params["second_subject_id"].(int)

	if !firstOk && !secondOk {
		return universities, nil
	}

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

func paginateUniversities(universities []*GetAllUniversityResponse, params map[string]interface{}) (*UniversitySearchResult, error) {
	totalCount := len(universities)

	// Determine the page number
	page := 1
	if p, ok := params["page"].(int); ok && p > 0 {
		page = p
	}

	// Determine the number of items per page
	perPage := 10
	if pp, ok := params["per_page"].(int); ok && pp > 0 {
		perPage = pp
	}

	// Calculate the total number of pages
	totalPages := (totalCount + perPage - 1) / perPage

	// Determine the slice range
	start := (page - 1) * perPage
	end := start + perPage

	// Handle out-of-range indices
	if start >= totalCount {
		universities = []*GetAllUniversityResponse{}
	} else if end > totalCount {
		universities = universities[start:totalCount]
	} else {
		universities = universities[start:end]
	}

	// Create the result struct
	result := &UniversitySearchResult{
		Universities: universities,
		Page:         page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
	}

	return result, nil
}

func UpdateUniversityServices(universityID int, services []*Service) error {
	o := orm.NewOrm()

	err := o.Begin()
	if err != nil {
		return err
	}

	_, err = o.Raw("DELETE FROM university_service WHERE university_id = ?", universityID).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	if len(services) > 0 {
		query := "INSERT INTO university_service (university_id, service_id) VALUES "
		values := make([]interface{}, 0)

		for i, service := range services {
			if i > 0 {
				query += ", "
			}
			query += "(?, ?)"
			values = append(values, universityID, service.Id)
		}

		_, err = o.Raw(query, values...).Exec()
		if err != nil {
			o.Rollback()
			return err
		}
	}

	err = o.Commit()
	if err != nil {
		o.Rollback()
		return err
	}

	return nil
}

func (u *University) RemoveGalleryPhoto(photoID int) error {
	o := orm.NewOrm()

	photo := &Gallery{Id: photoID}

	if err := o.QueryTable(new(Gallery)).Filter("Id", photoID).Filter("University__Id", u.Id).One(photo); err != nil {
		if err == orm.ErrNoRows {
			return fmt.Errorf("photo with ID %d not found in university gallery", photoID)
		}
		return fmt.Errorf("failed to find photo: %v", err)
	}

	if _, err := o.Delete(photo); err != nil {
		return fmt.Errorf("failed to delete photo: %v", err)
	}

	return nil
}

func UploadFileToCloud(filePath string, file multipart.File) (string, error) {
	awsAccessKey, _ := beego.AppConfig.String("aws_access_key")
	awsSecretKey, _ := beego.AppConfig.String("aws_secret_key")
	bucket, _ := beego.AppConfig.String("bucket")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		Endpoint:         aws.String("https://chi-sextans.object.pscloud.io"),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	uploader := s3.New(sess)

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	var contentType string
	switch {
	case strings.HasSuffix(filePath, ".svg"):
		contentType = "image/svg+xml"
	case strings.HasSuffix(filePath, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(filePath, ".webp"):
		contentType = "image/webp"
	default:
		contentType = http.DetectContentType(buf.Bytes())
	}

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filePath),
		Body:        bytes.NewReader(buf.Bytes()),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType), // Устанавливаем MIME-тип
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	fileURL := fmt.Sprintf("https://chi-sextans.object.pscloud.io/%s/%s", bucket, filePath)
	return fileURL, nil
}

func GetUniversityNames(lang string) ([]GetUniNamesResponse, error) {
	o := orm.NewOrm()
	var universities []University
	var response []GetUniNamesResponse

	_, err := o.QueryTable(new(University)).Filter("popular", true).All(&universities)
	if err != nil {
		return nil, err
	}

	for _, uni := range universities {
		var name string
		switch lang {
		case "ru":
			name = uni.AbbreviationRu
		case "kz":
			name = uni.AbbreviationKz
		default:
			name = uni.Name
		}

		response = append(response, GetUniNamesResponse{
			Id:           uni.Id,
			Abbreviation: name,
		})
	}

	return response, nil
}

// GetUniversityServices retrieves the services associated with a given university ID.
func GetUniversityServices(universityId int64) ([]*Service, error) {
	o := orm.NewOrm()
	var services []*Service

	// Assuming there's a many-to-many relationship set up between University and Service
	_, err := o.QueryTable("service").
		Filter("universities__university__id", universityId).
		All(&services)

	if err != nil {
		return nil, fmt.Errorf("error fetching services for university ID %d: %v", universityId, err)
	}

	return services, nil
}
