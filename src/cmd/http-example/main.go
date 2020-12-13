package main

import (
    "encoding/json"
    "fmt"
    "net/http"

    "fast-graphql/src/backend"
)

type User struct {
    Id      int64   `json:"id"`
    Name    string  `json:"name"`
    Email   string  `json:"email"`
    Married bool    `json:"married"`
    Height  float64 `json:"height"`
}

var users = []User{
    {
        Id:    1,
        Name: "Bob",
        Email: "bob@email.com",
        Married: false,
        Height: 172.53,
    },
    {
        Id:    2,
        Name: "Alice",
        Email: "Alice@email.com",
        Married: false,
        Height: 175.2,
    },
    {
        Id:    3,
        Name: "Tim",
        Email: "Tim@email.com",
        Married: true,
        Height: 162.3,
    },
    {
        Id:    4,
        Name: "Peter",
        Email: "Peter@email.com",
        Married: false,
        Height: 181.9,
    },
    {
        Id:    5,
        Name: "Juice",
        Email: "Juice@email.com",
        Married: true,
        Height: 132.9,
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
            "married": &backend.ObjectField{
                Name: "married",
                Type: backend.Bool,
            },
            "height": &backend.ObjectField{
                Name: "height",
                Type: backend.Float,
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
                    "id": &backend.Argument{
                        Name: "id",
                        Type: backend.Int,
                    },
                    "name": &backend.Argument{
                        Name: "name",
                        Type: backend.String,
                    },
                    "married": &backend.Argument{
                        Name: "married",
                        Type: backend.Bool,
                    },
                    "height": &backend.Argument{
                        Name: "height",
                        Type: backend.Float,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    fmt.Printf("\033[33m    [INTO] user defined ResolveFunction:  \033[0m\n")

                    id, ok := p.Arguments["id"].(int)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] id:  \033[0m\n")

                        // Find user
                        for _, user := range users {
                            if int(user.Id) == id {
                                return user, nil
                            }
                        }
                    }
                    name, ok := p.Arguments["name"].(string)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] name:  \033[0m\n")

                        // Find user
                        for _, user := range users {
                            if user.Name == name {
                                return user, nil
                            }
                        }
                    }
                    married, ok := p.Arguments["married"].(bool)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] married:  \033[0m\n")

                        // Find user
                        targetUsers := make([]User, 0)
                        for _, user := range users {
                            if bool(user.Married) == married {
                                targetUsers = append(targetUsers, user)
                            }
                        }
                        return targetUsers, nil
                    }
                    height, ok := p.Arguments["height"].(float64)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] height:  \033[0m\n")

                        // Find user
                        for _, user := range users {
                            if float64(user.Height) == height {
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


