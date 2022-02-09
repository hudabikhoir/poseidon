package content

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"boilerplate-golang-v2/business"
	"boilerplate-golang-v2/business/content"
)

//MySQLRepository The implementation of content.Repository object
type MySQLRepository struct {
	db *sql.DB
}

//NewMySQLRepository Generate mongo DB content repository
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db,
	}
}

//FindContentByID Find content based on given ID. Its return nil if not found
func (repo *MySQLRepository) FindContentByID(ID string) (*content.Content, error) {
	var content content.Content

	selectQuery := `SELECT id, name, description, created_at, created_by, modified_at, modified_by, version, COALESCE(tags, "")
		FROM content i
		LEFT JOIN (
			SELECT content_id, 
			GROUP_CONCAT(tag) as tags
			FROM content_tag GROUP BY content_id
		)AS it ON i.id = it.content_id
		WHERE i.id = ?`

	var tags string
	err := repo.db.
		QueryRow(selectQuery, ID).
		Scan(
			&content.ID, &content.Name, &content.Description,
			&content.CreatedAt, &content.CreatedBy,
			&content.ModifiedAt, &content.ModifiedBy,
			&content.Version, &tags)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	content.Tags = constructTagArray(tags)

	return &content, nil
}

//FindAllByTag Find all contents based on given tag. Its return empty array if not found
func (repo *MySQLRepository) FindAllByTag(tag string) ([]content.Content, error) {
	//TODO: if feel have a performance issue in tag grouping, move the logic from db to here
	selectQuery := `SELECT id, name, description, created_at, created_by, modified_at, modified_by, version, COALESCE(tags, "")
		FROM content i
		LEFT JOIN (
			SELECT content_id, 
			GROUP_CONCAT(tag) as tags
			FROM content_tag GROUP BY content_id
		)AS it ON i.id = it.content_id
		WHERE i.id IN (
			SELECT content_id
			FROM content_tag
			WHERE tag = ?	
		)`

	row, err := repo.db.Query(selectQuery, tag)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var contents []content.Content

	for row.Next() {
		var content content.Content
		var tags string

		err := row.Scan(
			&content.ID, &content.Name, &content.Description,
			&content.CreatedAt, &content.CreatedBy,
			&content.ModifiedAt, &content.ModifiedBy,
			&content.Version, &tags)

		if err != nil {
			return nil, err
		}

		content.Tags = constructTagArray(tags)
		contents = append(contents, content)
	}

	if err != nil {
		return nil, err
	}

	return contents, nil
}

//InsertContent Insert new content into database. Its return content id if success
func (repo *MySQLRepository) InsertContent(content content.Content) error {
	fmt.Println("masuk repository")
	ctx := context.Background()
	fmt.Println("repo db:", repo.db)
	tx, err := repo.db.BeginTx(ctx, nil)
	fmt.Println("err repo:", err)
	if err != nil {
		return err
	}

	contentQuery := fmt.Sprintf(`INSERT INTO content (name, description, created_at, created_by, modified_at, modified_by, version) VALUES ('%v', '%v', NOW(), '%v', NOW(), '%v', '%v')`,
		content.Name,
		content.Description,
		content.CreatedBy,
		content.ModifiedBy,
		content.Version)

	fmt.Println("contentQuery:", contentQuery)
	if err != nil {
		return err
	}

	_, err = tx.Exec(contentQuery)

	if err != nil {
		tx.Rollback()
		return err
	}

	tagQuery := `INSERT INTO content_tag (id, content_id, tag) VALUES (1, 2, 'kesehatan')`

	// for _, tag := range content.Tags {
	_, err = tx.Exec(tagQuery)
	// _, err = tx.Exec(tagQuery, content.ID, tag)

	if err != nil {
		tx.Rollback()
		return err
	}
	// }

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

//UpdateContent Update existing content in database
func (repo *MySQLRepository) UpdateContent(content content.Content, currentVersion int) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	contentInsertQuery := `UPDATE content 
		SET
			name = ?,
			description = ?,
			modified_at = ?,
			modified_by = ?,
			version = ?
		WHERE id = ? AND version = ?`

	res, err := tx.Exec(contentInsertQuery,
		content.Name,
		content.Description,
		content.ModifiedAt,
		content.ModifiedBy,
		content.Version,
		content.ID,
		currentVersion,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		tx.Rollback()
		return err
	}

	if affected == 0 {
		tx.Rollback()
		return business.ErrZeroAffected
	}

	//TODO: maybe better if we only delete the record that we need to delete
	//add logic slice to find which deleted and which want to added
	tagDeleteQuery := "DELETE FROM content_tag WHERE content_id = ?"
	_, err = tx.Exec(tagDeleteQuery, content.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tagUpsertQuery := "INSERT INTO content_tag (content_id, tag) VALUES (?, ?)"

	for _, tag := range content.Tags {
		_, err = tx.Exec(tagUpsertQuery, content.ID, tag)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func constructTagArray(tags string) []string {
	if tags == "" {
		return make([]string, 0)
	}

	return strings.Split(tags, ",")
}
