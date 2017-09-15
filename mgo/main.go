package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	defer db.Close()

	h := Adapt(http.HandlerFunc(handle), withDB(db))

	http.Handle("/comments", context.ClearHandler(h))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Adapter type
type Adapter func(http.Handler) http.Handler

// Adapt func
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Session) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbsession := db.Copy()
			defer dbsession.Close()

			context.Set(r, "database", dbsession)

			h.ServeHTTP(w, r)
		})
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

type comment struct {
	ID     bson.ObjectId `json:"id" bson:"id"`
	Author string        `json:"author" bson:"author"`
	Text   string        `json:"text" bson:"text"`
	When   time.Time     `json:"when" bson:"when"`
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var c comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.ID = bson.NewObjectId()
	c.When = time.Now()

	if err := db.DB("comments").C("comments").Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var comments []*comment
	if err := db.DB("comments").C("comments").Find(nil).Sort("-when").Limit(100).All(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
