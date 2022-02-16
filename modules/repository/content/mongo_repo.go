package content

import (
	"boilerplate-golang-v2/business/content"
	"boilerplate-golang-v2/util"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//MongoDBRepository The implementation of content.Repository object
type MongoDBRepository struct {
	col *mongo.Collection
}

type collection struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Tags        []string           `bson:"tags"`
	CreatedAt   time.Time          `bson:"created_at"`
	CreatedBy   string             `bson:"created_by"`
	ModifiedAt  time.Time          `bson:"modified_at"`
	ModifiedBy  string             `bson:"modified_by"`
	Version     int                `bson:"version"`
}

func newCollection(content content.Content) (*collection, error) {
	var UIDString string
	if content.ID != "" {
		UIDString = content.ID
	} else {
		UIDString = util.GenerateID()
	}
	objectID, err := primitive.ObjectIDFromHex(UIDString)

	if err != nil {
		return nil, err
	}

	return &collection{
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

func (col *collection) MongoDBToContent() content.Content {
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

//NewMongoDBRepository Generate mongo DB content repository
func NewMongoDBRepository(db *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		db.Collection("contents"),
	}
}

//FindContentByID Find content based on given ID. Its return nil if not found
func (repo *MongoDBRepository) FindContentByID(ID string) (*content.Content, error) {
	var col collection

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		//if cannot be convert means that ID will be never found
		return nil, nil
	}

	filter := bson.M{
		"_id": objectID,
	}

	if err := repo.col.FindOne(context.TODO(), filter).Decode(&col); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	content := col.MongoDBToContent()
	return &content, nil
}

//FindAllByTag Find all contents based on given tag. Its return empty array if not found
func (repo *MongoDBRepository) FindAllByTag(tag string) ([]content.Content, error) {
	filter := bson.M{
		"tags": bson.M{
			"$all": [1]string{tag},
		},
	}

	cursor, err := repo.col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var contents []content.Content

	for cursor.Next(context.TODO()) {
		var col collection
		if err = cursor.Decode(&col); err != nil {
			return nil, err
		}

		contents = append(contents, col.MongoDBToContent())
	}

	return contents, nil
}

//InsertContent Insert new content into database. Its return content id if success
func (repo *MongoDBRepository) InsertContent(content content.Content) (ID string, err error) {
	col, err := newCollection(content)

	if err != nil {
		return "0", err
	}

	_, err = repo.col.InsertOne(context.TODO(), col)

	if err != nil {
		return "0", err
	}

	return col.ID.Hex(), nil
}

//UpdateContent Update existing content in database
func (repo *MongoDBRepository) UpdateContent(content content.Content, currentVersion int) error {
	col, err := newCollection(content)

	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     col.ID,
		"version": currentVersion,
	}

	updated := bson.M{
		"$set": col,
	}

	_, err = repo.col.UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		return err
	}

	return nil
}
