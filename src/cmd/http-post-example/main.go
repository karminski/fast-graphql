package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"

    "fast-graphql/src/backend"
    "github.com/davecgh/go-spew/spew"

)

const (
    Gender_Male   = "MALE"
    Gender_Female = "FEMALE"
)

type User struct {
    Id      int64   `json:"id"`
    Name    string  `json:"name"`
    Email   string  `json:"email"`
    Married bool    `json:"married"`
    Height  float64 `json:"height"`
    Gender  string  `json:"gender"`
}

var users = []User{
    {
        Id:    1,
        Name: "Bob",
        Email: "bob@email.com",
        Married: false,
        Height: 172.53,
        Gender: Gender_Male,
    },
    {
        Id:    2,
        Name: "Alice",
        Email: "Alice@email.com",
        Married: false,
        Height: 175.2,
        Gender: Gender_Female,
    },
    {
        Id:    3,
        Name: "Tim",
        Email: "Tim@email.com",
        Married: true,
        Height: 162.3,
        Gender: Gender_Male,
    },
    {
        Id:    4,
        Name: "Peter",
        Email: "Peter@email.com",
        Married: false,
        Height: 181.9,
        Gender: Gender_Male,
    },
    {
        Id:    5,
        Name: "Juice",
        Email: "Juice@email.com",
        Married: true,
        Height: 132.9,
        Gender: Gender_Female,
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
            "gender": &backend.ObjectField{
                Name: "gender",
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
                    "gender": &backend.Argument{
                        Name: "gender",
                        Type: backend.String,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

                    fmt.Printf("\033[33m    [INTO] user defined ResolveFunction:  \033[0m\n")

                    spewo.Dump(p.Arguments)
                    id, ok := p.Arguments["id"].(float64)
                    spewo.Dump(id)

                    if ok {
                        intId := int(id)
                        fmt.Printf("\033[33m    [INTO] id:  \033[0m\n")

                        // Find user
                        for _, user := range users {
                            if int(user.Id) == intId {
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
                        for _, user := range users {
                            if bool(user.Married) == married {
                                return user, nil
                            }
                        }
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
                    gender, ok := p.Arguments["gender"].(string)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] gender:  \033[0m\n")

                        // Find gender
                        for _, user := range users {
                            if user.Gender == gender {
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


func executeQuery(query string, variables map[string]interface{}, schema backend.Schema) *backend.Result {
    var result *backend.Result 
    // execute
    result = backend.Execute(backend.Request{
        Schema: schema,
        Query:  query,
        Variables: variables,
    })
    // if len(result.Errors) > 0 {
    //     fmt.Printf("errors: %v", result.Errors)
    // }
    return result
}

func main() {
    http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("\n\n\033[33m////////////////////////////////////////// Request Start ///////////////////////////////////////\033[0m\n")

        spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

        // HTTP Post method
        var decodedVariables map[string]interface{}
        body, _ := ioutil.ReadAll(r.Body)
        spewo.Dump(body)
        json.Unmarshal([]byte(body), &decodedVariables)

        // HTTP Get method 
        query     := decodedVariables["query"].(string)
        variables := decodedVariables["variables"].(map[string]interface{})
        
        // execute
        result    := executeQuery(query, variables, schema)

        // return
        w.Header().Set("content-type","text/json")
        json.NewEncoder(w).Encode(result)
    })
    fmt.Printf("START.\n")

    fmt.Println("Server is running on port 8081")
    http.ListenAndServe("127.0.0.1:8081", nil)

    fmt.Printf("EXIT. \n")
}


