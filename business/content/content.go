package content

import "time"

//Content product content that available to rent or sell
type Content struct {
	ID          string
	Name        string
	Description string
	Tags        []string
	CreatedAt   time.Time
	CreatedBy   string
	ModifiedAt  time.Time
	ModifiedBy  string
	Version     int
}

//NewContent create new content
func NewContent(
	name string,
	description string,
	tags []string,
	creator string,
	createdAt time.Time) Content {

	return Content{
		Name:        name,
		Description: description,
		Tags:        tags,
		CreatedAt:   createdAt,
		CreatedBy:   creator,
		ModifiedAt:  createdAt,
		ModifiedBy:  creator,
		Version:     1,
	}
}

//ModifyContent update existing content data
func (oldContent *Content) ModifyContent(newName string, newDescription string, newTags []string, updater string, modifiedAt time.Time) Content {
	return Content{
		ID:          oldContent.ID,
		Name:        newName,
		Description: newDescription,
		Tags:        newTags,
		CreatedAt:   oldContent.CreatedAt,
		CreatedBy:   oldContent.CreatedBy,
		ModifiedAt:  modifiedAt,
		ModifiedBy:  updater,
		Version:     oldContent.Version + 1,
	}
}
