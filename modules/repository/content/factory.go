package content

import (
	"boilerplate-golang-v2/business/content"
	"boilerplate-golang-v2/util"
)

//RepositoryFactory Will return business.content.Repository based on active database connection
func RepositoryFactory(dbCon *util.DatabaseConnection) content.Repository {
	var contentRepo content.Repository
	if dbCon.Driver == util.MySQL {
		contentRepo = NewMySQLRepository(dbCon.PostgreSQL)
	} else if dbCon.Driver == util.MongoDB {
		contentRepo = NewMongoDBRepository(dbCon.MongoDB)
	}

	return contentRepo
}
