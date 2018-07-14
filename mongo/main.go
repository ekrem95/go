package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/connstring"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var (
	uri        = "mongodb://localhost:27017"
	collection = "documents"

	db   *mongo.Database
	coll *mongo.Collection
)

// Document ...
type Document struct {
	Title string   `json:"title" bson:"title"`
	Data  string   `json:"data" bson:"data"`
	Tags  []string `json:"tags" bson:"tags"`
	Meta  struct {
		Date  time.Time `json:"date" bson:"date"`
		Count int32     `json:"count" bson:"count"`
	} `json:"meta" bson:"meta"`
}

func main() {
	// parse uri by using connstring.Parse()
	connectionString, err := connstring.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	// if database is not specified in connectionString
	// set database name
	dbname := connectionString.Database
	if dbname == "" {
		dbname = "new_database"
	}

	// connect to mongo
	client, err := mongo.Connect(context.Background(), uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	// set global database and collection values
	db = client.Database(dbname, nil)
	coll = db.Collection(collection)

	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/add-document", addDocument)
	e.GET("/find", find)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// list all documents
func hello(c echo.Context) error {
	var doc Document
	docs, err := doc.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, docs)
}

// add one document
func addDocument(c echo.Context) error {
	doc := new(Document)
	if err := c.Bind(doc); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	doc.Meta.Date = time.Now()
	doc.Meta.Count = int32(time.Now().Second())

	if err := doc.InsertOne(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, doc)
}

// find multiple documents that matches params
func find(c echo.Context) error {
	title := c.QueryParam("title")
	countString := c.QueryParam("count")

	count, err := strconv.ParseInt(countString, 10, 32)
	if err != nil {
		count = 0
	}

	var doc Document
	docs, err := doc.Find(title, int32(count))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, docs)
}

// FindAll ...
func (doc Document) FindAll() ([]Document, error) {
	var docs []Document

	cursor, err := coll.Find(
		context.Background(),
		bson.NewDocument(),
	)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// InsertOne ...
func (doc Document) InsertOne() error {
	if _, err := coll.InsertOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("title", doc.Title),
			bson.EC.String("data", doc.Data),
			bson.EC.ArrayFromElements("tags",
				bson.VC.String("cool"),
				bson.VC.String("stuff"),
			),
			bson.EC.SubDocumentFromElements("meta",
				bson.EC.Time("date", doc.Meta.Date),
				bson.EC.Int32("count", doc.Meta.Count),
			),
		)); err != nil {
		return err
	}

	return nil
}

// Find ...
func (doc Document) Find(title string, count int32) ([]Document, error) {
	var docs []Document

	cursor, err := coll.Find(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("title", title),
			bson.EC.SubDocumentFromElements("meta.count",
				bson.EC.Int32("$gt", count),
			)),
	)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}
