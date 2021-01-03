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

const (
    Country_A = "Country A"
    Country_B = "Country B"
    Country_C = "Country C"
    Country_D = "Country D"
    Country_E = "Country E"
)

const (
    City_X1 = "City X1"
    City_X2 = "City X2"
    City_X3 = "City X3"
    City_X4 = "City X4"
    City_X5 = "City X5"
)

type User struct {
    Id      int       `json:"id"`
    Name    string    `json:"name"`
    Email   string    `json:"email"`
    Married bool      `json:"married"`
    Height  float64   `json:"height"`
    Gender  string    `json:"gender"`
    Friends []int     `json:"friends"`
    Location Location `json:"location"`

}

type Location struct {
    Country string `json:"country"`
    City    string `json:"city"`
}

var users = []User{
    {
        Id:    1,
        Name: "Bob",
        Email: "bob@email.com",
        Married: false,
        Height: 172.53,
        Gender: Gender_Male,
        Friends: []int{2,3,4},
        Location: Location{Country_A, City_X1},
    },
    {
        Id:    2,
        Name: "Alice",
        Email: "Alice@email.com",
        Married: false,
        Height: 175.2,
        Gender: Gender_Female,
        Friends: []int{1},
        Location: Location{Country_B, City_X2},
    },
    {
        Id:    3,
        Name: "Tim",
        Email: "Tim@email.com",
        Married: true,
        Height: 162.3,
        Gender: Gender_Male,
        Friends: []int{1,4},
        Location: Location{Country_C, City_X3},
    },
    {
        Id:    4,
        Name: "Peter",
        Email: "Peter@email.com",
        Married: false,
        Height: 181.9,
        Gender: Gender_Male,
        Friends: []int{1,3},
        Location: Location{Country_D, City_X4},
    },
    {
        Id:    5,
        Name: "Juice",
        Email: "Juice@email.com",
        Married: true,
        Height: 132.9,
        Gender: Gender_Female,
        Friends: []int{},
        Location: Location{Country_E, City_X5},
    },
}

var locationType, _ = backend.NewObject(
    backend.ObjectTemplate{
        Name: "Location",
        Fields: backend.ObjectFields{
            "Country": &backend.ObjectField{
                Name: "Country",
                Type: backend.String,
            },
            "City": &backend.ObjectField{
                Name: "City",
                Type: backend.String,
            },
        },
    },
)

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
            "Friends": &backend.ObjectField{
                Name: "Friends",
                Type: backend.NewList(backend.Int),
                // ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                //     spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
                //     fmt.Printf("\033[31m[INTO] func Friends.ResolveFunction()  \033[0m\n")
                //     friends, ok := p.Source.([]int)
                //     fmt.Printf("\033[33m    [DUMP] p.Arguments:  \033[0m\n")
                //     spewo.Dump(p.Arguments)
                //     os.Exit(1)
                //     matchedUsers := []User{}
                //     if ok {
                //         userIdIndex := make(map[int]int, len(users))
                //         // build user id index
                //         for i, user := range users {
                //             userIdIndex[user.Id] = i
                //         }
                //         // match user
                //         for _, friendsId := range friends {
                //             i    := userIdIndex[friendsId]
                //             user := users[i]
                //             matchedUsers = append(matchedUsers, user)
                //         }
                //     }
                //     return matchedUsers, nil
                // },
            },
            "Location": &backend.ObjectField{
                Name: "Location",
                Type: locationType,
            },
        },
    },
)


var queryObject, _ = backend.NewObject(
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

                    id, ok := p.Arguments["id"].(int)
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
                    return nil, errors.New("ResolveFunction(): target data not found.")
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



var mutationObject, _ = backend.NewObject(
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
                    email, emailOk     := p.Arguments["email"].(string)
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


