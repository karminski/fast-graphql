// fast-graphql/main.go

package main

import (
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "fast-graphql/src/frontend"
)


const version = `fast-graphql 0.0.1`
const usage = `
fast-graphql tester usage: %s [options] [filename]
Avaliable options are:
    -o name  output to file 'name' (default is "fastgraphqlc.out")
`

func main() {
    // input params
    _v := flag.Bool("v", false, "")
    // _o := flag.String("o", "fastgraphqlc.out", "")
    flag.Usage = printUsage
    flag.Parse()

    /* Print Version */
    if *_v {
        fmt.Println(version)
        return
    }
    
    if len(flag.Args()) != 1 {
        printUsage()
        fmt.Printf("%v", flag.Args())
        return
    }
    
    // load target file
    filename := flag.Args()[0]
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    fileName := "@" + filename
    fileContent := string(data)

    // compile
    fmt.Println("\n\n-----------------------------------------------------------------------------------\n\n")
    fmt.Println("fast-graphql client tester\n\n")
    fmt.Println("- [PARSING PHRASE] --------------------------------------------------------------------\n\n")
    fmt.Printf("Now parsing file: %s\n", fileName)

    document := frontend.Compile(fileContent)

    fmt.Println("\n\n")
    fmt.Println("- [EXECUTE PHRASE] --------------------------------------------------------------------\n\n")

    // execute 
    
    return
}


func printUsage() {
    fmt.Printf(usage, os.Args[0])
}