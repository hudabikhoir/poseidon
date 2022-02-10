package content

import (
	"boilerplate-golang-v2/business"
	"boilerplate-golang-v2/business/content/spec"
	"time"

	validator "github.com/go-playground/validator/v10"
)

//Repository ingoing port for content
type Repository interface {
	//FindContentByID If data not found will return nil without error
	FindContentByID(ID string) (content *Content, err error)

	//FindAllByTag If no data match with the given tag, will return empty slice instead of nil
	FindAllByTag(tag string) (contents []Content, err error)

	//InsertContent Insert new content into storage
	InsertContent(content Content) (id string, err error)

	//UpdateContent if data not found will return core.ErrZeroAffected
	UpdateContent(content Content, currentVersion int) (err error)
}

//Service outgoing port for content
type Service interface {
	GetContentByID(ID string) (content *Content, err error)

	GetContentsByTag(tag string) (contents []Content, err error)

	CreateContent(upsertcontentSpec spec.UpsertContentSpec, createdBy string) (id string, err error)

	UpdateContent(ID string, upsertcontentSpec spec.UpsertContentSpec, currentVersion int, modifiedBy string) (err error)
}

//=============== The implementation of those interface put below =======================

type service struct {
	repository Repository
	validate   *validator.Validate
}

//NewService Construct content service object
func NewService(repository Repository) Service {
	return &service{
		repository,
		validator.New(),
	}
}

//GetContentByID Get content by given ID, return nil if not exist
func (s *service) GetContentByID(ID string) (*Content, error) {
	return s.repository.FindContentByID(ID)
}

//GetContentsByTag Get all contents by given tag, return zero array if not match
func (s *service) GetContentsByTag(tag string) ([]Content, error) {
	contents, err := s.repository.FindAllByTag(tag)
	if err != nil || contents == nil {
		return []Content{}, err
	}

	return contents, err
}

//CreateContent Create new content and store into database
func (s *service) CreateContent(upsertcontentSpec spec.UpsertContentSpec, createdBy string) (string, error) {
	err := s.validate.Struct(upsertcontentSpec)

	if err != nil {
		return "", business.ErrInvalidSpec
	}

	content := NewContent(
		upsertcontentSpec.Name,
		upsertcontentSpec.Description,
		upsertcontentSpec.Tags,
		createdBy,
		time.Now(),
	)
	ID, err := s.repository.InsertContent(content)
	if err != nil {
		return "", err
	}

	return ID, nil
}

//UpdateContent Update existing content in the database.
//Will return ErrNotFound when content is not exists or ErrConflict if data version is not match
func (s *service) UpdateContent(ID string, upsertcontentSpec spec.UpsertContentSpec, currentVersion int, modifiedBy string) error {
	err := s.validate.Struct(upsertcontentSpec)

	if err != nil || len(ID) == 0 {
		return business.ErrInvalidSpec
	}

	//get the content first to make sure data is exist
	content, err := s.repository.FindContentByID(ID)

	if err != nil {
		return err
	} else if content == nil {
		return business.ErrNotFound
	} else if content.Version != currentVersion {
		return business.ErrHasBeenModified
	}

	newContent := content.ModifyContent(upsertcontentSpec.Name, upsertcontentSpec.Description, upsertcontentSpec.Tags, modifiedBy, time.Now())

	return s.repository.UpdateContent(newContent, currentVersion)
}
