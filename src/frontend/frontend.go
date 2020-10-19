// frontend.go
// 

package frontend

import (
    "fmt"
    "reflect"
    "strings"
    "github.com/davecgh/go-spew/spew"
)

func Compile(query string) *Document {
    lexer := NewLexer(query)
    document   := parseDocument(lexer)
    lexer.NextTokenIs(TOKEN_EOF) // set EOF
    fmt.Printf("%+v\n", document)
    ASTVisitor(document)
    fmt.Printf("--------------------------------------\n")
    spew.Dump(document)

    return document
}

/**
 * ASTVisitor, Visit AST(Document) and dump it for debug.
 * @param {[type]} document *Document [description]
 */
func ASTVisitor(document *Document) {
    fmt.Printf("----------------------- \n")
    fmt.Printf("ASTVisitor: \n\n")
    fmt.Printf("Document {\n")
    fmt.Printf("  LastLineNum: %d\n", document.LastLineNum)
    for _, definition := range document.Definitions {
        fmt.Printf("  definition: %v\n", definition)
        examiner(reflect.TypeOf(definition), 0)
        DumpOperationDefinition(definition.(*OperationDefinition), 1)

    }
    fmt.Printf("}\n")

}

func examiner(t reflect.Type, depth int) {
    fmt.Println(strings.Repeat("  ", depth), t.Name(), ":")
    switch t.Kind() {
    case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
        fmt.Println(strings.Repeat("  ", depth+1), "Contained type:")
        examiner(t.Elem(), depth+1)
    case reflect.Struct:
        for i := 0; i < t.NumField(); i++ {
            f := t.Field(i)
            fmt.Println(strings.Repeat("  ", depth+1), f.Name, ": ", f.Type.Kind())
            if f.Tag != "" {
                fmt.Println(strings.Repeat("  ", depth+2), "Tag is", f.Tag)
                fmt.Println(strings.Repeat("  ", depth+2), "tag1 is", f.Tag.Get("tag1"), "tag2 is", f.Tag.Get("tag2"))
            }
        }
    }
}


func DumpOperationDefinition(operationDefinition *OperationDefinition, paddingNum int) {
    pad := strings.Repeat("  ", paddingNum)
    fmt.Printf("%sOperationDefinition { \n", pad)
    if operationDefinition == (&OperationDefinition{}) {
        fmt.Printf("%s  empty", pad)
    } else {
        fmt.Printf("%s  LineNum:               %d\n", pad, operationDefinition.LineNum)
        fmt.Printf("%s  OperationType:       \n", pad)
        fmt.Printf("%s  OperationName:       \n", pad)
        fmt.Printf("%s  VariableDefinitions: \n", pad)
        fmt.Printf("%s  Directives:          \n", pad)
        fmt.Printf("%s  SelectionSet:        \n", pad)
    }
    fmt.Printf("%s}\n", pad)
}




// func DumpOperationType(operationType *OperationType, paddingNum int) {
//     pad := strings.Repeat("  ", paddingNum)
//     fmt.Printf("%OperationType { \n", pad)
//     if operationDefinition == (&OperationType{}) {
//         fmt.Printf("%s  empty", pad)
//     } else {
//         fmt.Printf("%s  LineNum:               %d\n", pad, operationDefinition.LineNum)
//         fmt.Printf("%s  OperationType:       \n", pad)
//         fmt.Printf("%s  OperationName:       \n", pad)
//         fmt.Printf("%s  VariableDefinitions: \n", pad)
//         fmt.Printf("%s  Directives:          \n", pad)
//         fmt.Printf("%s  SelectionSet:        \n", pad)
//     }
//     fmt.Printf("%s}\n", pad)
// }