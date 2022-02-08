package content

import (
	"boilerplate-golang-v2/business/content"
	"context"
	"time"

	"github.com/go-kivik/kivik"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CouchDBRepository The implementation of content.Repository object
type CouchDBRepository struct {
	col *kivik.DB
}

type couchdbCol struct {
	ID          primitive.ObjectID `json:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tags        []string           `json:"tags"`
	CreatedAt   time.Time          `json:"created_at"`
	CreatedBy   string             `json:"created_by"`
	ModifiedAt  time.Time          `json:"modified_at"`
	ModifiedBy  string             `json:"modified_by"`
	Version     int                `json:"version"`
}

type CouchDBQuery struct {
	Selector map[string]interface{} `json:"selector"`
	Limit    uint64                 `json:"limit,omitempty"`
	Skip     uint64                 `json:"skip,omitempty"`
	Fields   []string               `json:"fields,omitempty"`
	Sort     []map[string]string    `json:"sort,omitempty"`
}

func newCouchdbCollection(content content.Content) (*couchdbCol, error) {
	objectID, err := primitive.ObjectIDFromHex(content.ID)

	if err != nil {
		return nil, err
	}

	return &couchdbCol{
		objectID,
		content.Name,
		content.Description,
		content.Tags,
		content.CreatedAt,
		content.CreatedBy,
		content.ModifiedAt,
		content.ModifiedBy,
		content.Version,
	}, nil
}

func (col *couchdbCol) CouchDBToContent() content.Content {
	var content content.Content
	content.ID = col.ID.Hex()
	content.Name = col.Name
	content.Description = col.Description
	content.Tags = col.Tags
	content.CreatedAt = col.CreatedAt
	content.CreatedBy = col.CreatedBy
	content.ModifiedAt = col.ModifiedAt
	content.ModifiedBy = col.ModifiedBy
	content.Version = col.Version

	return content
}

//NewCouchDBRepository Generate mongo DB article repository
func NewCouchDBRepository(client *kivik.Client) *CouchDBRepository {
	return &CouchDBRepository{
		col: client.DB(context.TODO(), "feedfund_contents"),
	}
}

//FindContentByID Find content based on given ID. Its return nil if not found
func (repo *CouchDBRepository) FindContentByID(ID string) (*content.Content, error) {
	// var col couchdbCol

	// objectID, err := primitive.ObjectIDFromHex(ID)
	// if err != nil {
	// 	//if cannot be convert means that ID will be never found
	// 	return nil, nil
	// }

	// filter := bson.M{
	// 	"_id": objectID,
	// }

	// if err := repo.col.FindOne(context.TODO(), filter).Decode(&col); err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return nil, nil
	// 	}

	// 	return nil, err
	// }

	// content := col.ToContent()
	// return &content, nil
	return nil, nil
}

//FindAllByTag Find all contents based on given tag. Its return empty array if not found
func (repo *CouchDBRepository) FindAllByTag(tag string) ([]content.Content, error) {
	var col couchdbCol

	contents := []content.Content{}
	query := CouchDBQuery{
		Selector: map[string]interface{}{
			"system_success": true,
		},
	}

	row, err := repo.col.Find(context.TODO(), query)

	if err != nil {
		return contents, err
	}
	defer row.Close()

	for row.Next() {
		if err = row.ScanDoc(&col); err != nil {
			return contents, err
		}
		contents = append(contents, col.CouchDBToContent())
	}

	return contents, nil
}

//InsertContent Insert new content into database. Its return content id if success
func (repo *CouchDBRepository) InsertContent(content content.Content) error {
	// col, err := newCollection(content)

	// if err != nil {
	// 	return err
	// }

	// _, err = repo.col.InsertOne(context.TODO(), col)

	// if err != nil {
	// 	return err
	// }

	return nil
}

//UpdateContent Update existing content in database
func (repo *CouchDBRepository) UpdateContent(content content.Content, currentVersion int) error {
	// col, err := newCollection(content)

	// if err != nil {
	// 	return err
	// }

	// filter := bson.M{
	// 	"_id":     col.ID,
	// 	"version": currentVersion,
	// }

	// updated := bson.M{
	// 	"$set": col,
	// }

	// _, err = repo.col.UpdateOne(context.TODO(), filter, updated)
	// if err != nil {
	// 	return err
	// }

	return nil
}
