package util

import (
	"boilerplate-golang-v2/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/go-kivik/couchdb/v4"
	"github.com/go-kivik/kivik/v4"

	_ "github.com/lib/pq"
)

//DatabaseDriver Database driver enum
type DatabaseDriver string

const (
	//MongoDB DatabaseDriver
	MongoDB DatabaseDriver = "mongodb"
	//MySQL DatabaseDriver
	MySQL DatabaseDriver = "mysql"
	//CouchDB DatabaseDriver
	CouchDB DatabaseDriver = "couchdb"
	//Postgre DatabaseDriver
	PostgreSQL DatabaseDriver = "postgressql"
)

//DatabaseConnection Database connection
type DatabaseConnection struct {
	Driver DatabaseDriver

	//for MySQL
	MySQLDB *sql.DB

	//for MongoDB
	MongoDB     *mongo.Database
	mongoClient *mongo.Client

	//for couchdb
	CouchDBClient *kivik.Client

	PostgreSQL *sql.DB
}

//NewDatabaseConnection Create new database connection based on given config
func NewDatabaseConnection(config *config.AppConfig) *DatabaseConnection {
	var db DatabaseConnection
	//define the data repository
	switch config.Database.Driver {
	case "mysql":
		//initiate mysql db repository
		db.MySQLDB = newMysqlDB(config)
		db.Driver = MySQL
	case "sqlite":
		//initiate mysqlite db repository
		db.MySQLDB = newSQLiteDBClient(config)
		db.Driver = MySQL
	case "postgressql":
		//initiate postgreSQL db repository
		db.PostgreSQL = newPostgreSQL(config)
		db.Driver = PostgreSQL
	case "mongodb":
		// initiate mongodb repository
		db.mongoClient = newMongoDBClient(config)
		db.MongoDB = db.mongoClient.Database(config.Database.Name)
		db.Driver = MongoDB
	case "couchdb":
		//initiate mysql db repository
		db.CouchDBClient = newCouchDBClient(config)
		db.Driver = CouchDB
	default:
		panic("Unsupported nosql database driver")
	}

	return &db
}

//CloseConnection Close db connection
func (db *DatabaseConnection) CloseConnection() {
	if db.MySQLDB != nil {
		db.MySQLDB.Close()
	}

	if db.mongoClient != nil {
		db.mongoClient.Disconnect(context.Background())
	}
}

func newMysqlDB(config *config.AppConfig) *sql.DB {
	var uri string

	uri = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
		config.Database.Port,
		config.Database.Name)

	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	return db
}

func newSQLiteDBClient(config *config.AppConfig) *sql.DB {
	db := fmt.Sprintf("../%v", config.Database.Name)
	sqliteDatabase, _ := sql.Open("sqlite3", db) // Open the created SQLite File
	return sqliteDatabase
}

func newMongoDBClient(config *config.AppConfig) *mongo.Client {
	uri := "mongodb://"

	if config.Database.Username != "" {
		uri = fmt.Sprintf("%s%v:%v@", uri, config.Database.Username, config.Database.Password)
	}

	uri = fmt.Sprintf("%s%v:%v",
		uri,
		config.Database.Address,
		config.Database.Port)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	return client
}

func newCouchDBClient(config *config.AppConfig) *kivik.Client {
	client, err := kivik.New("couch", fmt.Sprintf(
		"https://%s:%s@%s/",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
	))
	if err != nil {
		panic(err)
	}

	return client
}

// newPostgreSQL return a client connection handle to a Postgre server.
func newPostgreSQL(config *config.AppConfig) *sql.DB {
	var uri string
	uri = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
		config.Database.Port,
		config.Database.Name)
	client, err := sql.Open("postgres", uri)
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.PingContext(ctx)
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	return client
}
