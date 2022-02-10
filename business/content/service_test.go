package content_test

import (
	"boilerplate-golang-v2/business"
	"boilerplate-golang-v2/business/content"
	"boilerplate-golang-v2/business/content/spec"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

var service content.Service
var content1, content2 content.Content
var insertSpec, updateSpec, failedSpec, errorSpec spec.UpsertContentSpec
var creator, updater, errorFindID string
var errorInsert error = errors.New("error on insert")
var errorFind error = errors.New("error on find")

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestGetContentByID(t *testing.T) {
	t.Run("Expect found the content", func(t *testing.T) {
		foundContent, _ := service.GetContentByID(content1.ID)
		if !reflect.DeepEqual(*foundContent, content1) {
			t.Error("Expect content has to be equal with content1", foundContent, content1)
		}
	})

	t.Run("Expect not found the content", func(t *testing.T) {
		content, err := service.GetContentByID("random")

		if err != nil {
			t.Error("Expect error is nil. Error: ", err)
		} else if content != nil {
			t.Error("Expect content must be not found (nil)")
		}
	})
}

func TestGetContentByTags(t *testing.T) {
	t.Run("Expect found the contents", func(t *testing.T) {
		contents, _ := service.GetContentsByTag("tag2")

		if len(contents) != 2 {
			t.Error("Expect content length must be two")
			t.FailNow()
		}

		if reflect.DeepEqual(contents[0], content1) {
			if !reflect.DeepEqual(contents[1], content2) {
				t.Error("Expect 2nd content is equal to content2")
			}
		} else if reflect.DeepEqual(contents[0], content2) {
			if !reflect.DeepEqual(contents[1], content1) {
				t.Error("Expect 2nd content is equal to content1")
			}
		} else {
			t.Error("Expect contents is content1 and content2")
		}
	})

	t.Run("Expect not found the contents", func(t *testing.T) {
		contents, err := service.GetContentsByTag("not-found-tag")

		if err != nil {
			t.Error("Expect error is nil", err)
		} else if contents == nil {
			t.Error("Expect contents is not nil")
		} else if len(contents) != 0 {
			t.Error("Expect contents is not found")
		}
	})
}

func TestCreateContent(t *testing.T) {
	t.Run("Expect success create content", func(t *testing.T) {
		id, err := service.CreateContent(insertSpec, creator)

		if err != nil {
			t.Error("Expext error is not nil. Error: ", err)
			t.FailNow()
		}

		for _, tag := range insertSpec.Tags {
			contents, _ := service.GetContentsByTag(tag)

			if len(contents) == 0 {
				t.Error("Expect at least one content when search by given tag: ", tag)
				continue
			}

			isFound := false
			for _, content := range contents {
				if content.ID == id {
					isFound = true
					break
				}
			}

			if !isFound {
				t.Error("Expect found inserted content when search by given tag: ", tag)
			}
		}

		newContent, _ := service.GetContentByID(id)

		if newContent == nil {
			t.Error("Expect content is not nil after inserted")
			t.FailNow()
		}

		if newContent.Name != insertSpec.Name {
			t.Error("Expect name is equal as given")
		}

		if newContent.Description != insertSpec.Description {
			t.Error("Expect description is equal as given")
		}

		if !reflect.DeepEqual(newContent.Tags, insertSpec.Tags) {
			t.Error("Expect tags is equal as given")
		}

		if newContent.CreatedBy != creator {
			t.Error("Expect created by is equal to " + creator)
		}

		if newContent.ModifiedBy != creator {
			t.Error("Expect modified by is equal to " + creator)
		}

		if newContent.CreatedAt != newContent.ModifiedAt {
			t.Error("Expect created at and modified at is equal")
		}

		if newContent.Version != 1 {
			t.Error("Expect version is equal to 1")
		}
	})

	t.Run("Expect failed create content on spec", func(t *testing.T) {
		_, err := service.CreateContent(failedSpec, creator)

		if err == nil {
			t.Error("Expect error is not nil")
		} else if err != business.ErrInvalidSpec {
			t.Error("Expect error invalid spec. Error is: ", err)
		}
	})

	t.Run("Expect failed create content on repository", func(t *testing.T) {
		_, err := service.CreateContent(errorSpec, creator)

		if err == nil {
			t.Error("Expect error is not nil")
		} else if err != errorInsert {
			t.Error("Expect error on insert. Error is: ", err)
		}
	})
}

func TestUpdateContent(t *testing.T) {
	t.Run("Expect success update content", func(t *testing.T) {
		id := content2.ID
		version := content2.Version
		oldTags := content2.Tags

		service.UpdateContent(id, updateSpec, version, updater)

		//find the old tag that doesn't exist in new updated tags
		var invalidateTags []string
		for _, tag := range oldTags {
			isExistOnUpdatedTag := false

			for _, updatedTag := range updateSpec.Tags {
				if tag == updatedTag {
					isExistOnUpdatedTag = true
					break
				}
			}

			if !isExistOnUpdatedTag {
				invalidateTags = append(invalidateTags, tag)
			}
		}

		//verify the invalidated tag is not contain the content anymore
		for _, invalidateTag := range invalidateTags {
			tagContents, _ := service.GetContentsByTag(invalidateTag)
			isFound := false

			for _, tagContent := range tagContents {
				if tagContent.ID == id {
					isFound = true
					break
				}
			}

			if isFound {
				t.Error("Expect not found when search by old invalidate tag: ", invalidateTag)
			}
		}

		contents, _ := service.GetContentsByTag(updateSpec.Tags[0])

		isFound := false
		for _, content := range contents {
			if content.ID == id {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error("Expect found inserted content when search by given tag: ", updateSpec.Tags[0])
		}

		updatedContent, _ := service.GetContentByID(content2.ID)

		if updatedContent == nil {
			t.Error("Expect content is not nil after updated")
			t.FailNow()
		}

		if updatedContent.Name != updateSpec.Name {
			t.Error("Expect name is equal as given")
		}

		if updatedContent.Description != updateSpec.Description {
			t.Error("Expect description is equal as given")
		}

		if !reflect.DeepEqual(updatedContent.Tags, updateSpec.Tags) {
			t.Error("Expect tags is equal as given")
		}

		if updatedContent.CreatedBy != content2.CreatedBy {
			t.Error("Expect created by is equal to " + content2.CreatedBy)
		}

		if updatedContent.ModifiedBy != updater {
			t.Error("Expect modified by is equal to " + updater)
		}

		if updatedContent.CreatedAt == updatedContent.ModifiedAt {
			t.Error("Expect created at and modified at is not equal")
		}

		if updatedContent.Version != content2.Version+1 {
			t.Error("Expect version was increase by one")
		}
	})

	t.Run("Expect failed update content on spec", func(t *testing.T) {
		err := service.UpdateContent(content2.ID, failedSpec, content2.Version, updater)

		if err == nil {
			t.Error("Expect error is not nil")
		} else if err != business.ErrInvalidSpec {
			t.Error("Expect error invalid spec. Error is: ", err)
		}
	})

	t.Run("Expect failed update content on not found", func(t *testing.T) {
		err := service.UpdateContent("not-found", updateSpec, 1, updater)

		if err == nil {
			t.Error("Expect error is not nil")
		} else if err != business.ErrNotFound {
			t.Error("Expect error content not found. Error is: ", err)
		}
	})

	t.Run("Expect failed update content on wrong version", func(t *testing.T) {
		err := service.UpdateContent(content1.ID, updateSpec, content1.Version+1, updater)

		if err == nil {
			t.Error("Expect error is not nil")
		} else if err != business.ErrHasBeenModified {
			t.Error("Expect error content has been modified. Error is: ", err)
		}
	})

	t.Run("Expect failed update content on repository", func(t *testing.T) {
		err := service.UpdateContent(errorFindID, updateSpec, 1, updater)

		if err == nil {
			t.Error("Expect error is not nil")
		} else if err != errorFind {
			t.Error("Expect error on insert. Error is: ", err)
		}
	})
}

func setup() {
	//initialize content1
	content1.ID = "5f350b7d21148431abc65290"
	content1.Name = "Content one"
	content1.Description = "Description one"
	content1.Tags = []string{"tag1", "tag2"}
	content1.Version = 1
	content1.CreatedAt = time.Now()
	content1.CreatedBy = "creator one"
	content1.ModifiedAt = content1.CreatedAt
	content1.ModifiedBy = content1.CreatedBy

	//initialize content 2
	content2.ID = "5f351360ac84a3bb1baee057"
	content2.Name = "Content two"
	content2.Description = "Description two"
	content2.Tags = []string{"tag2", "tag3", "tag4"}
	content2.Version = 2
	content2.CreatedAt = time.Now().Add(time.Minute * -15)
	content2.CreatedBy = "creator two"
	content2.ModifiedAt = time.Now()
	content2.ModifiedBy = "updater two"

	repo := newInMemoryRepository()
	service = content.NewService(&repo)

	insertSpec.Name = "New Content"
	insertSpec.Description = "New Description"
	insertSpec.Tags = []string{"tag99"}

	updateSpec.Name = "Update Content"
	updateSpec.Description = "Update Description"
	updateSpec.Tags = []string{"tag99-updated"}

	failedSpec.Name = ""
	failedSpec.Description = "Failed Description"
	failedSpec.Tags = []string{}

	errorSpec.Name = "Error Content"
	errorSpec.Description = "Error Description"
	errorSpec.Tags = []string{}

	creator = "creator"
	updater = "updater"

	errorFindID = "error-find-id"
}

type inMemoryRepository struct {
	contentByID  map[string]content.Content
	contentByTag map[string][]content.Content
}

func newInMemoryRepository() inMemoryRepository {
	var repo inMemoryRepository
	repo.contentByID = make(map[string]content.Content)
	repo.contentByTag = make(map[string][]content.Content)

	repo.contentByID[content1.ID] = content1
	repo.contentByID[content2.ID] = content2

	for _, tag := range content1.Tags {
		contents := repo.contentByTag[tag]
		repo.contentByTag[tag] = append(contents, content1)
	}

	for _, tag := range content2.Tags {
		contents := repo.contentByTag[tag]
		repo.contentByTag[tag] = append(contents, content2)
	}

	return repo
}

func (repo *inMemoryRepository) FindContentByID(ID string) (*content.Content, error) {
	if ID == errorFindID {
		return nil, errorFind
	}

	content, ok := repo.contentByID[ID]
	if !ok {
		return nil, nil
	}

	return &content, nil
}

func (repo *inMemoryRepository) FindAllByTag(tag string) ([]content.Content, error) {
	var contents []content.Content
	contents, ok := repo.contentByTag[tag]

	if !ok {
		return contents, nil
	}

	return contents, nil
}

func (repo *inMemoryRepository) InsertContent(content content.Content) (string, error) {
	if content.Name == errorSpec.Name {
		return "", errorInsert
	}

	repo.contentByID[content.ID] = content

	for _, tag := range content.Tags {
		contents := repo.contentByTag[tag]
		repo.contentByTag[tag] = append(contents, content)
	}
	return content.ID, nil
}

func (repo *inMemoryRepository) UpdateContent(content content.Content, currentVersion int) error {
	//cleanup old tag first
	oldContent := repo.contentByID[content.ID]

	//cleanup the old tags first
	for _, tag := range oldContent.Tags {
		tagContents, _ := repo.FindAllByTag(tag)

		contentIndex := -1
		for idx, tagContent := range tagContents {
			if tagContent.ID == content.ID {
				contentIndex = idx
				break
			}
		}

		if contentIndex != -1 {
			tagContents = append(tagContents[:contentIndex], tagContents[contentIndex+1:]...)
		}

		repo.contentByTag[tag] = tagContents
	}

	repo.contentByID[content.ID] = content

	//adding the new tag
	for _, tag := range content.Tags {
		contents := repo.contentByTag[tag]
		repo.contentByTag[tag] = append(contents, content)
	}
	return nil
}
