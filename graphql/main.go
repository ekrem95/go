package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User ...
type User struct {
	ID   int    `json:"id"   gorm:"type:serial;primary key"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
	Age  int    `json:"age"  gorm:"not null"`
}

// Graphql ...
type Graphql struct{}

var graph Graphql

func (Graphql) open() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=root dbname=graphql password=pass sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (Graphql) insert(u *User) error {
	db, err := graph.open()
	if err != nil {
		return err
	}
	defer db.Close()

	db.Create(&User{Name: u.Name, Age: u.Age})
	return nil
}

func (Graphql) find() ([]User, error) {
	db, err := graph.open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var users []User
	db.Find(&users)

	return users, nil
}

func (Graphql) findOne(id int) (User, error) {
	user := User{}

	db, err := graph.open()
	if err != nil {
		return user, err
	}
	defer db.Close()

	db.First(&user, "id = ?", id)
	return user, nil
}

func (Graphql) update(id, age int) (User, error) {
	user := User{}

	db, err := graph.open()
	if err != nil {
		return user, err
	}
	defer db.Close()

	db.First(&user, "id = ?", id)
	if user.ID == 0 {
		return user, nil
	}
	user.Age = age
	db.Save(&user)

	return user, nil
}

func init() {
	db, err := graph.open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})
}

// define custom GraphQL ObjectType `userType` for struct `User`
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// Mutation ...
var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// http://localhost:8080/graphql?query=mutation+_{addUser(name:"Steve",age:50){id,name,age}}
		"addUser": &graphql.Field{
			Type:        userType,
			Description: "Add new user",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"age": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name, _ := params.Args["name"].(string)
				age, _ := params.Args["age"].(int)

				newUser := User{Name: name, Age: age}
				graph.insert(&newUser)

				return newUser, nil
			},
		},
		// http://localhost:8080/graphql?query=mutation+_{updateUser(id:1,age:27){id,name,age}}
		"updateUser": &graphql.Field{
			Type:        userType,
			Description: "Update an existing user's age",
			Args: graphql.FieldConfigArgument{
				"age": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				age, _ := params.Args["age"].(int)
				id, _ := params.Args["id"].(int)

				user, _ := graph.update(id, age)
				if user.ID == 0 {
					return nil, nil
				}
				return user, nil
			},
		},
	},
})

// Query ...
var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		// http://localhost:8080/graphql?query={user(id:1){id,name,age}}
		"user": &graphql.Field{
			Type:        userType,
			Description: "Get a user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int)

				var user User
				if ok {
					user, _ = graph.findOne(id)
				}
				if user.ID == 0 {
					return nil, nil
				}
				return user, nil
			},
		},
		// http://localhost:8080/graphql?query={userList{id,name,age}}
		"userList": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "List of users",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				users, _ := graph.find()
				return users, nil
			},
		},
	},
})

// define schema, with our Query and Mutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    Query,
	Mutation: Mutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	port := 8080
	fmt.Printf("server is running on port %d", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
