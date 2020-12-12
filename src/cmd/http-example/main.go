package main

import (
    "encoding/json"
    "fmt"
    "net/http"

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
    {
        Id:    3,
        Name: "Tim",
        Email: "Tim@email.com",
    },
    {
        Id:    4,
        Name: "Peter",
        Email: "Peter@email.com",
    },
    {
        Id:    5,
        Name: "Juice",
        Email: "Juice@email.com",
    },
}

var userType, _ = backend.NewObject(
    backend.ObjectTemplate{
        Name: "User",
        Fields: backend.ObjectFields{
            "id": &backend.ObjectField{
                Name: "id",
                Type: backend.Int,
            },
            "name": &backend.ObjectField{
                Name: "name",
                Type: backend.String,
            },
            "email": &backend.ObjectField{
                Name: "email",
                Type: backend.String,
            },
        },
    },
)

var queryType, _ = backend.NewObject(
    backend.ObjectTemplate{
        Name: "Query",
        Fields: backend.ObjectFields{
            // field User
            "user": &backend.ObjectField{
                Name: "user",
                Type: userType,
                Description: "Get user by id",
                Arguments: &backend.Arguments{
                    // "id": &backend.Argument{
                    //     Name: "id",
                    //     Type: backend.Int,
                    // },
                    "name": &backend.Argument{
                        Name: "name",
                        Type: backend.String,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    id, ok := p.Arguments["id"].(int)
                    if ok {
                        // Find user
                        for _, user := range users {
                            if int(user.Id) == id {
                                return user, nil
                            }
                        }
                    }
                    name, ok := p.Arguments["name"].(string)
                    if ok {
                        // Find user
                        for _, user := range users {
                            if user.Name == name {
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
                Type: backend.NewList(userType),
                Description: "Get user list",
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    return users, nil
                },
            },
        },
    },
)


var schema, _ = backend.NewSchema(
    backend.SchemaTemplate{
        Query: queryType,
    },
)


func executeQuery(query string, schema backend.Schema) *backend.Result {
    result := backend.Execute(backend.Request{
        Schema: schema,
        Query:  query,
    })
    // if len(result.Errors) > 0 {
    //     fmt.Printf("errors: %v", result.Errors)
    // }
    return result
}

func main() {
    http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
        result := executeQuery(r.URL.Query().Get("query"), schema)
        json.NewEncoder(w).Encode(result)
    })
    fmt.Printf("START.\n")

    fmt.Println("Server is running on port 8081")
    http.ListenAndServe("127.0.0.1:8081", nil)

    fmt.Printf("EXIT. \n")
}


