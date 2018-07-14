package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// User ...
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Users ...
var Users []User

func uniqueID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func init() {
	user1 := User{ID: "1", Name: "John", Age: 16}
	user2 := User{ID: "2", Name: "Andrew", Age: 26}
	Users = append(Users, user1, user2)
}

// define custom GraphQL ObjectType `userType` for struct `User`
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
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

				// generate new id
				id := uniqueID()

				newUser := User{
					ID:   id,
					Name: name,
					Age:  age,
				}

				Users = append(Users, newUser)

				return newUser, nil
			},
		},
		// http://localhost:8080/graphql?query=mutation+_{updateUser(id:"1",age:27){id,name,age}}
		"updateUser": &graphql.Field{
			Type:        userType,
			Description: "Update an existing user's age",
			Args: graphql.FieldConfigArgument{
				"age": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				age, _ := params.Args["age"].(int)
				id, _ := params.Args["id"].(string)
				user := User{}

				for i := 0; i < len(Users); i++ {
					if Users[i].ID == id {
						Users[i].Age = age
						user = Users[i]
						break
					}
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
		// http://localhost:8080/graphql?query={user(id:"1"){id,name,age}}
		"user": &graphql.Field{
			Type:        userType,
			Description: "Get a user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(string)
				if ok {
					for _, user := range Users {
						if user.ID == id {
							return user, nil
						}
					}
				}

				return User{}, nil
			},
		},
		// http://localhost:8080/graphql?query={userList{id,name,age}}
		"userList": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "List of users",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Users, nil
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
