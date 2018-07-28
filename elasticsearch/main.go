package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	uuid "github.com/satori/go.uuid"
	validator "gopkg.in/go-playground/validator.v9"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string                `json:"user" validate:"required"`
	Message  string                `json:"message" validate:"required"`
	Retweets int                   `json:"retweets" validate:"required"`
	Image    string                `json:"image,omitempty"`
	Date     time.Time             `json:"date,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"date":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

var (
	ctx    context.Context
	client *elastic.Client

	elasticAddress = "http://127.0.0.1:9200"
	elasticIndex   = "twitter"
	elasticType    = "tweet"
	validate       = validator.New()
)

func main() {
	r := gin.Default()
	r.GET("/tweets", find)
	r.GET("/tweet", findOne)
	r.GET("/search", search)
	r.POST("/new", insert)
	r.POST("/bulk", bulk)
	r.GET("/update", update)
	r.GET("/delete-index", deleteIndex)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

func init() {
	var err error
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx = context.Background()

	// Obtain a client and connect to the Elasticsearch installation
	for {
		client, err = elastic.NewClient(elastic.SetURL(elasticAddress), elastic.SetSniff(false))
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(3 * time.Second)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(elasticAddress).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(elasticIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(elasticIndex).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func search(c *gin.Context) {
	// Parse request
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query not specified"})
		return
	}
	skip := 0
	take := 10

	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err == nil {
		take = i
	}

	esQuery := elastic.NewMultiMatchQuery(query, "user", "message").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := client.Search().
		Index(elasticIndex).
		Query(esQuery).
		From(skip).Size(take).
		Do(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var tweets []Tweet
	for _, hit := range result.Hits.Hits {
		var tweet Tweet
		json.Unmarshal(*hit.Source, &tweet)
		tweets = append(tweets, tweet)
	}
	c.JSON(http.StatusOK, gin.H{
		"time":      fmt.Sprintf("%d", result.TookInMillis),
		"hits":      fmt.Sprintf("%d", result.Hits.TotalHits),
		"documents": tweets,
	})
}

func bulk(c *gin.Context) {
	var tweets []Tweet
	c.BindJSON(&tweets)

	for _, tweet := range tweets {
		if err := validate.Struct(tweet); err != nil {
			errs := strings.Split(err.Error(), "\n")

			c.JSON(http.StatusBadRequest, errs)
			return
		}
	}

	bulk := client.
		Bulk().
		Index(elasticIndex).
		Type(elasticType)
	for i, tweet := range tweets {
		id := uuid.Must(uuid.NewV4()).String()
		date := time.Now()
		tweet.Date = date
		tweets[i].Date = date

		bulk.Add(elastic.NewBulkIndexRequest().Id(id).Doc(tweet))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(200, tweets)
}

func insert(c *gin.Context) {
	var tweet Tweet
	c.Bind(&tweet)

	if err := validate.Struct(tweet); err != nil {
		errs := strings.Split(err.Error(), "\n")

		c.JSON(http.StatusBadRequest, errs)
		return
	}

	id := uuid.Must(uuid.NewV4()).String()
	tweet.Date = time.Now()

	put, err := client.Index().
		Index(elasticIndex).
		Type(elasticType).
		Id(id).
		BodyJson(tweet).
		Do(ctx)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put.Id, put.Index, put.Type)

	// Flush to make sure the documents got written.
	if _, err = client.Flush().Index(elasticIndex).Do(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, tweet)
}

func findOne(c *gin.Context) {
	id := c.Query("id")

	// Get tweet with specified ID
	tweet, err := client.Get().
		Index(elasticIndex).
		Type(elasticType).
		Id(id).
		Do(ctx)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	if tweet.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", tweet.Id, tweet.Version, tweet.Index, tweet.Type)
	}

	c.JSON(http.StatusOK, tweet)
}

func find(c *gin.Context) {
	user := c.Query("user")

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", user)
	searchResult, err := client.Search().
		Index(elasticIndex). // search in index "twitter"
		Query(termQuery).    // specify the query
		Sort("user", true).  // sort by "user" field, ascending
		From(0).Size(10).    // take documents 0-9
		Pretty(true).        // pretty print request and response JSON
		Do(ctx)              // execute
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	c.JSON(http.StatusOK, searchResult)

	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
}

func update(c *gin.Context) {
	id := c.Query("id")

	// increment the number of retweets.
	update, err := client.Update().Index(elasticIndex).Type(elasticType).Id(id).
		Script(elastic.NewScriptInline("ctx._source.retweets += params.num").Lang("painless").Param("num", 1)).
		// Upsert(map[string]interface{}{"retweets": 0}).
		Do(ctx)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	fmt.Printf("New version of tweet %q is now %d\n", update.Id, update.Version)

	c.JSON(http.StatusOK, update)
}

func deleteIndex(c *gin.Context) {
	// Delete an index.
	deleteIndex, err := client.DeleteIndex(elasticIndex).Do(ctx)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}

	c.JSON(http.StatusOK, deleteIndex)
}
