package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"

    "fast-graphql/src/frontend"
    "fast-graphql/src/backend"
)

type User struct {
    Id    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {
        Id:    1,
        Name: "Bob",
        Email: "bob@email.com",
    },
    {
        Id:    2,
        Name: "Alice",
        Email: "Alice@email.com",
    },
}

var userType = backend.NewObject(
    backend.ObjectTemplate{
        Name: "User",
        Fields: backend.ObjectFields{
            "id": &backend.ObjectField{
                Name: "id",
                Type: backend.Int
            },
            "name": &backend.ObjectField{
                Name: "name",
                Type: Backend.String
            },
            "email": &backend.ObjectField{
                Name: "email",
                Type: Backend.String
            },
        }
    }
)

var queryType = backend.NewObject(
    backend.ObjectConfig{
        Name: "Query",
        Fields: backend.ObjectFields{
            // field User
            "user": &backend.ObjectField{
                Name: "user",
                Type: userType
                Description: "Get user by id",
                Arguments: &backend.Arguments{
                    "id": &backend.Argument{
                        Name: "id",
                        Type: backend.Int
                    }
                },
                ResolveFunction: func(p graphql.ResolveParams) (interface{}, error) {
                    id, ok := p.Args["id"].(int)
                    if ok {
                        // Find user
                        for _, user := range users {
                            if int(user.Id) == id {
                                return user, nil
                            }
                        }
                    }
                    return nil, nil
                },
            },
            // Field List
            "list": &backend.ObjectField{
                Name: "list",
                Type: backend.NewList(userType)
                Description: "Get user list",
                ResolveFunction: func(params graphql.ResolveParams) (interface{}, error) {
                    return users, nil
                },
            }
        }
    }
)


var schema, _ = backend.NewSchema(
    backend.SchemaTemplate{
        Query: queryType
    }
)


func executeQuery(query string, schema graphql.Schema) *graphql.Result {
    result := graphql.Do(graphql.Params{
        Schema:        schema,
        RequestString: query,
    })
    if len(result.Errors) > 0 {
        fmt.Printf("errors: %v", result.Errors)
    }
    return result
}

func getType(i interface{}) {
    name := reflect.TypeOf(i).Elem().Name()
    fmt.Println("Name is: "+ name)
}

func main() {


    http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
        result := executeQuery(r.URL.Query().Get("query"), schema)
        json.NewEncoder(w).Encode(result)
    })

    fmt.Println("Server is running on port 8080")
    http.ListenAndServe("127.0.0.1:8080", nil)
}


