package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "math/rand"
    "fast-graphql/src/backend"
    "fast-graphql/src/frontend"
    "github.com/davecgh/go-spew/spew"
    "errors"
    "strconv"
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

var friendType, _ = backend.NewObject(
    backend.ObjectTemplate{
        Name: "Friends",
        Fields: backend.ObjectFields{
            "Id": &backend.ObjectField{
                Name: "Id",
                Type: backend.Int,
            },
            "Name": &backend.ObjectField{
                Name: "Name",
                Type: backend.String,
            },
            "Email": &backend.ObjectField{
                Name: "Email",
                Type: backend.String,
            },
            "Married": &backend.ObjectField{
                Name: "Married",
                Type: backend.Bool,
            },
            "Height": &backend.ObjectField{
                Name: "Height",
                Type: backend.Float,
            },
            "Gender": &backend.ObjectField{
                Name: "Gender",
                Type: backend.String,
            },
            "Location": &backend.ObjectField{
                Name: "Location",
                Type: locationType,
            },
        },
    },
)


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
            "Id": &backend.ObjectField{
                Name: "Id",
                Type: backend.Int,
            },
            "Name": &backend.ObjectField{
                Name: "Name",
                Type: backend.String,
            },
            "Email": &backend.ObjectField{
                Name: "Email",
                Type: backend.String,
            },
            "Married": &backend.ObjectField{
                Name: "Married",
                Type: backend.Bool,
            },
            "Height": &backend.ObjectField{
                Name: "Height",
                Type: backend.Float,
            },
            "Gender": &backend.ObjectField{
                Name: "Gender",
                Type: backend.String,
            },
            "Friends": &backend.ObjectField{
                Name: "Friends",
                Type: backend.NewList(friendType),
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
                    fmt.Printf("\033[31m[INTO] func Friends.ResolveFunction()  \033[0m\n")
                    spewo.Dump(p)
                   // convert source input
                    var user User 
                    var ok   bool
                    if user, ok = p.Source.(User); !ok {
                        return nil, errors.New("func Friends.Resolve() can not resolve p.Source.")
                    }
                    friends := user.Friends
                    matchedUsers := []User{}
                    // get friends
                    userIdIndex := make(map[int]int, len(users))
                    // build user id index
                    for i, user := range users {
                        userIdIndex[user.Id] = i
                    }
                    // match user
                    for _, friendsId := range friends {
                        var user User
                        if i, ok    := userIdIndex[friendsId]; ok {
                            user = users[i]
                        }
                        matchedUsers = append(matchedUsers, user)
                    }
                    return matchedUsers, nil
                },
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
            "User": &backend.ObjectField{
                Name: "User",
                Type: userType,
                Description: "Get user by id",
                Arguments: &backend.Arguments{
                    "Id": &backend.Argument{
                        Name: "id",
                        Type: backend.Int,
                    },
                    "Name": &backend.Argument{
                        Name: "name",
                        Type: backend.String,
                    },
                    "Married": &backend.Argument{
                        Name: "married",
                        Type: backend.Bool,
                    },
                    "Height": &backend.Argument{
                        Name: "height",
                        Type: backend.Float,
                    },
                    "Gender": &backend.Argument{
                        Name: "gender",
                        Type: backend.String,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

                    fmt.Printf("\033[33m    [INTO] user defined ResolveFunction:  \033[0m\n")

                    id, ok := p.Arguments["Id"].(int)
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
                    name, ok := p.Arguments["Name"].(string)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] name:  \033[0m\n")

                        // Find user
                        for _, user := range users {
                            if user.Name == name {
                                return user, nil
                            }
                        }
                    }
                    married, ok := p.Arguments["Married"].(bool)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] married:  \033[0m\n")
                        // Find user
                        for _, user := range users {
                            if bool(user.Married) == married {
                                return user, nil
                            }
                        }
                    }
                    height, ok := p.Arguments["Height"].(float64)
                    if ok {
                        fmt.Printf("\033[33m    [INTO] height:  \033[0m\n")

                        // Find user
                        for _, user := range users {
                            if float64(user.Height) == height {
                                return user, nil
                            }
                        }
                    }
                    gender, ok := p.Arguments["Gender"].(string)
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
                    "Name": &backend.Argument{
                        Name: "Name",
                        Type: backend.String,
                    },
                    "Email": &backend.Argument{
                        Name: "Email",
                        Type: backend.String,
                    },
                    "Married": &backend.Argument{
                        Name: "Married",
                        Type: backend.Bool,
                    },
                    "Height": &backend.Argument{
                        Name: "Height",
                        Type: backend.Float,
                    },
                    "Gender": &backend.Argument{
                        Name: "Gender",
                        Type: backend.String,
                    },
                    "Friends": &backend.Argument{
                        Name: "Friends",
                        Type: backend.NewList(backend.Int),
                    },
                    "Location": &backend.Argument{
                        Name: "Location",
                        Type: locationType,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
                    fmt.Printf("\033[33m    [INTO] create ResolveFunction():  \033[0m\n")
                    fmt.Printf("\033[33m    [DUMP] p.Arguments:  \033[0m\n")
                    spewo.Dump(p.Arguments)
                    // fill friends
                    assertedFriends := p.Arguments["Friends"].([]frontend.Value)
                    friends := make([]int, len(assertedFriends))
                    for i, n := range assertedFriends {
                        friends[i] = n.(frontend.IntValue).Value
                    }
                    // fill location
                    var location Location
                    assertedLocation := p.Arguments["Location"].([]*frontend.ObjectField)
                    for _, n := range assertedLocation {
                        if n.Name.Value == "City" {
                            location.City = n.Value.(frontend.StringValue).Value
                        }
                        if n.Name.Value == "Country" {
                            location.Country = n.Value.(frontend.StringValue).Value
                        }
                    }
                    // fill user
                    user := User{
                        Id: rand.Intn(1000),
                        Name: p.Arguments["Name"].(string),
                        Email: p.Arguments["Email"].(string),
                        Married: p.Arguments["Married"].(bool),
                        Height: p.Arguments["Height"].(float64),
                        Gender: p.Arguments["Gender"].(string),
                        Friends: friends,
                        Location: location,
                    }
                    users = append(users, user)
                    fmt.Printf("\033[33m    [DUMP] user:  \033[0m\n")
                    spewo.Dump(user)
                    return user, nil
                },
            },
            "update": &backend.ObjectField{
                Name: "update",
                Type: userType,
                Description: "update user info",
                Arguments: &backend.Arguments{
                    "Id": &backend.Argument{
                        Name: "Id",
                        Type: backend.Int,
                    },
                    "Name": &backend.Argument{
                        Name: "Name",
                        Type: backend.String,
                    },
                    "Email": &backend.Argument{
                        Name: "Email",
                        Type: backend.String,
                    },
                    "Married": &backend.Argument{
                        Name: "Married",
                        Type: backend.Bool,
                    },
                    "Height": &backend.Argument{
                        Name: "Height",
                        Type: backend.Float,
                    },
                    "Gender": &backend.Argument{
                        Name: "Gender",
                        Type: backend.String,
                    },
                    "Friends": &backend.Argument{
                        Name: "Friends",
                        Type: backend.NewList(backend.Int),
                    },
                    "Location": &backend.Argument{
                        Name: "Location",
                        Type: locationType,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {
                    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
                    fmt.Printf("\033[33m    [INTO] update ResolveFunction():  \033[0m\n")
                    // base info
                    id, _                := p.Arguments["Id"].(int)
                    name, nameOk         := p.Arguments["Name"].(string)
                    email, emailOk       := p.Arguments["Email"].(string)
                    married, marriedOk   := p.Arguments["Married"].(bool)
                    height, heightOk     := p.Arguments["Height"].(float64)
                    gender, genderOk     := p.Arguments["Gender"].(string)
                    // extend info
                    // fill friends
                    var friends []int
                    friendsOk := false
                    if assertedFriends, ok := p.Arguments["Friends"].([]frontend.Value); ok {
                        for _, n := range assertedFriends {
                            friends = append(friends, n.(frontend.IntValue).Value)
                        }
                        friendsOk = true
                    }
                    // fill location
                    var location Location
                    locationOk := false
                    if assertedLocation, ok := p.Arguments["Location"].([]*frontend.ObjectField); ok {
                        for _, n := range assertedLocation {
                            if n.Name.Value == "City" {
                                location.City = n.Value.(frontend.StringValue).Value
                            }
                            if n.Name.Value == "Country" {
                                location.Country = n.Value.(frontend.StringValue).Value
                            }
                        }
                        locationOk = true
                    }
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
                            if friendsOk {
                                users[i].Friends = friends
                            }
                            if locationOk {
                                users[i].Location = location
                            }
                        }
                    }
                    fmt.Printf("\033[33m    [DUMP] users[target]:  \033[0m\n")
                    spewo.Dump(users[target])

                    return users[target], nil
                },
            },
            "delete": &backend.ObjectField{
                Name: "delete",
                Type: userType,
                Description: "delete user by id",
                Arguments: &backend.Arguments{
                    "Id": &backend.Argument{
                        Name: "Id",
                        Type: backend.Int,
                    },
                },
                ResolveFunction: func(p backend.ResolveParams) (interface{}, error) {

                    fmt.Printf("\033[33m    [INTO] delete ResolveFunction():  \033[0m\n")

                    id, _ := p.Arguments["Id"].(int)
                    // find target user and update
                    targetUser := User{}
                    hit := false
                    for i, user := range users {
                        if user.Id == id {
                            hit = true
                            targetUser = user
                            // remove user
                            users = append(users[:i], users[i+1:]...)
                        }
                    }
                    if !hit {
                        return nil, errors.New("ResolveFunction(): target user(Id:"+ strconv.Itoa(id)+") not found, can not delete.")
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


