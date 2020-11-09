// fast-graphql-backend/main.go

package main

import (
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "fast-graphql/src/frontend"
    "fast-graphql/src/backend"
    "github.com/davecgh/go-spew/spew"
)


const version = `fast-graphql-backend 0.0.1`
const usage = `
fast-graphql-frontend tester usage: %s [options] [filename]
Avaliable options are:
    -o name  output to file 'name' (default is "fastgraphqlc.out")
`

func main() {
    // pre config
    spewo := spew.ConfigState{
        Indent: "    ", 
        DisablePointerAddresses: true,
        DisableCapacities: true,
        DisableVariableType: true}

    // query string
    queryString := `
        query t1 {
            user{
                id
                name
                email
            }
        }
    `

    // query config
    fields := backend.Fields{
        "user": &backend.Field{
            Type: "string",
            ResolveFunc: func(i interface{}) (interface{}, error) {
                return {Id: "a2xrbHNka2ljdmlpaWFqbg==", Name: "Apple", Email: "apple@email.com"}, nil
            }
        }
    }

    // compile
    fmt.Printf("Now parsing: %s\n", queryString)
    ast := frontend.Compile(queryString)

    // dump ast
    spewo.Dump(ast)

    // execute
    r := backend.Execute(ast, fields)

    spewo.Dump(r)
    return
}


func printUsage() {
    fmt.Printf(usage, os.Args[0])
}