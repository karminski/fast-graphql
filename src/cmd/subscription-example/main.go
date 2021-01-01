package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "math/rand"
    "fast-graphql/src/backend"
    "github.com/davecgh/go-spew/spew"
    "errors"
    // "os"
)

const (
    Gender_Male   = "MALE"
    Gender_Female = "FEMALE"
)

type User struct {
    Id      int   `json:"id"`
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





var subscriptionObject, _ = backend.NewObject(
    backend.ObjectTemplate{
        Name: "Mutation",
        Fields: backend.ObjectFields{
            "create": &backend.ObjectField{
                Name: "create",
                Type: userType,
                Description: "create new user",
                Arguments: &backend.Arguments{
                    "name": &backend.Argument{
                        Name: "name",
                        Type: backend.String,
                    },
                    "email": &backend.Argument{
                        Name: "email",
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
                    user := User{
                        Id: rand.Intn(1000),
                        Name: p.Arguments["name"].(string),
                        Email: p.Arguments["email"].(string),
                        Married: p.Arguments["married"].(bool),
                        Height: p.Arguments["height"].(float64),
                        Gender: p.Arguments["gender"].(string),
                    }
                    users = append(users, user)
                    return user, nil
                },
            },
            "update": &backend.ObjectField{
                Name: "update",
                Type: userType,
                Description: "update user info",
                Arguments: &backend.Arguments{
                    "id": &backend.Argument{
                        Name: "id",
                        Type: backend.Int,
                    },
                    "name": &backend.Argument{
                        Name: "name",
                        Type: backend.String,
                    },
                    "email": &backend.Argument{
                        Name: "email",
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
                    id, _              := p.Arguments["id"].(int)
                    name, nameOk       := p.Arguments["name"].(string)
                    email, emailOk       := p.Arguments["email"].(string)
                    married, marriedOk := p.Arguments["married"].(bool)
                    height, heightOk   := p.Arguments["height"].(float64)
                    gender, genderOk   := p.Arguments["gender"].(string)
                    // find target user and update
                    target := 0
                    for i, user := range users {
                        if user.Id == id {
                            target = i
                            if nameOk {
                                users[i].Name = name
                            }
                            if emailOk {
                                users[i].Email = email
                            }
                            if marriedOk {
                                users[i].Married = married
                            }
                            if heightOk {
                                users[i].Height = height
                            }
                            if genderOk {
                                users[i].Gender = gender
                            }
                        }
                    }
                    return users[target], nil
                },
            },
            "delete": &backend.ObjectField{
                Name: "delete",
                Type: userType,
                Description: "delete user by id",
                Arguments: &backend.Arguments{
                    "id": &backend.Argument{
                        Name: "id",
                        Type: backend.Int,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    id, _ := p.Arguments["id"].(int)
                    // find target user and update
                    targetUser := User{}
                    for i, user := range users {
                        if user.Id == id {
                            targetUser = user
                            // remove user
                            users = append(users[:i], users[i+1:]...)
                        }
                    }
                    return targetUser, nil
                },
            },
        },
    },
)

var schema, _ = backend.NewSchema(
    backend.SchemaTemplate{
        Query: queryObject,
        Mutation: mutationObject,
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
    if len(result.Errors) > 0 {
        fmt.Printf("\n\n\n")
        fmt.Printf("errors: %v", result.Errors)
    }
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

        var variables map[string]interface{}
        query     := decodedVariables["query"].(string)
        if decodedVariables["variables"] != nil {
            variables = decodedVariables["variables"].(map[string]interface{})
        }
        
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


