package main

import (
	"context"
	"log"
	"net/http"
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
		Date time.Time `json:"date" bson:"date"`
	} `json:"meta" bson:"meta"`
}

func main() {
	connectionString, err := connstring.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	dbname := connectionString.Database
	if dbname == "" {
		dbname = "new_database"
	}

	client, err := mongo.Connect(context.Background(), uri, nil)
	if err != nil {
		log.Fatal(err)
	}

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

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func hello(c echo.Context) error {
	var doc Document
	docs, err := doc.findAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, docs)
}

func (doc Document) findAll() ([]Document, error) {
	var docs []Document

	cursor, err := coll.Find(
		context.Background(),
		// bson.NewDocument(bson.EC.String("title", "New Title")),
		nil,
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

func (doc Document) insertOne() error {
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
			),
		)); err != nil {
		return err
	}

	return nil
}

func addDocument(c echo.Context) error {
	doc := new(Document)
	if err := c.Bind(doc); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	doc.Meta.Date = time.Now()

	if err := doc.insertOne(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, doc)
}
