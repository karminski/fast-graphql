// executor.go
package backend

import (
    "fast-graphql/src/frontend"
    "fmt"
    "errors"
    "log"
    "reflect"
    "encoding/json"

    // "strconv"
    "os"
    "github.com/davecgh/go-spew/spew"

)

const DUMP_FRONTEND = false

type Request struct {
    // GraphQL Schema config for server side
    Schema Schema 

    // GraphQL Query string from client side
    Query string 

    // GraphQL Query variables from client side
    Variables map[string]interface{}
}

type Result struct {
    Data      interface{} `json:"data"`
    Errors []*ErrorInfo   `json:"errors"`
}

type ErrorInfo struct {
    Message   string
    Location *ErrorLocation
}

type ErrorLocation struct {
    Line  int 
    Col   int
}

// GlobalVariables for Query Variables, etc. 
type GlobalVariables struct {
    // asserted query variables from request.Variables by VariableDefinition filtered
    QueryVariablesMap map[string]interface{}

    // executor for operation subscription
    SubscriptionExecutor *SubscriptionExecutor
}

func (result *Result) SetErrorInfo(err error, errorLocation *ErrorLocation) {
    errStr := fmt.Sprintf("%v", err)
    errorInfo := ErrorInfo{errStr, errorLocation}
    result.Errors = append(result.Errors, &errorInfo)
}

func DecodeVariables(inputVariables string) (map[string]interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func DecodeVariables  \033[0m\n")

    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    var decodedVariables map[string]interface{}
    // no variables inputed
    if inputVariables == "" {
        return nil, nil
    }
    err := json.Unmarshal([]byte(inputVariables), &decodedVariables)
    fmt.Printf("\033[33m    [DUMP] decodedVariables:  \033[0m\n")
    spewo.Dump(decodedVariables)
    if err != nil {
        err := "executeQuery(): user input variables decode failed, please check input variables json syntax." 
        return nil, errors.New(err)
    }
    return decodedVariables, nil
}

func NewGlobalVariables() *GlobalVariables {
    g := &GlobalVariables{}
    if ENABLE_SUBSCRIPTION_EXECUTOR {
        g.SubscriptionExecutor = NewSubscriptionExecutor()
    }
    return g
}

func Execute(request Request) (*Result) {

    var document *frontend.Document
    var err       error

    result := Result{} 
    g      := NewGlobalVariables()

    // debugging
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    // process input
    if document, err = frontend.Compile(request.Query); err != nil {
        result.SetErrorInfo(err, nil)
        return &result
    }

    // @todo: THE DOCUMENT NEED VALIDATE!
    
    if DUMP_FRONTEND {
        fmt.Printf("\n")
        fmt.Printf("\033[33m    [DUMP] Document:  \033[0m\n")
        spewo.Dump(document)
        if true {
            result.Data = document
            return &result
        }
        fmt.Printf("\033[33m    [DUMP] Request:  \033[0m\n")
        spewo.Dump(request)
        os.Exit(1)
    }

    // get top layer SelectionSet.Fields and request.Schema.ObjectFields
    var operationDefinition *frontend.OperationDefinition
    if operationDefinition, err = document.GetOperationDefinition(); err != nil {
        result.SetErrorInfo(err, nil)
        return &result
    }

    // fill Query Variables Map
    if g.QueryVariablesMap, err = getQueryVariablesMap(request, operationDefinition.VariableDefinitions); err != nil {
        result.SetErrorInfo(err, nil)
        return &result
    }
    fmt.Printf("\033[33m    [DUMP] g.QueryVariablesMap:  \033[0m\n")
    spewo.Dump(g.QueryVariablesMap)

    selectionSet := operationDefinition.SelectionSet
    // selectionSetFields := getSelectionSetFields(selectionSet)

    // get schema object fields
    var objectFields ObjectFields
    operationType := operationDefinition.OperationType
    if operationType == frontend.OperationTypeQuery && request.Schema.Query != nil {
        objectFields = request.Schema.GetQueryObjectFields()
    } else if operationType == frontend.OperationTypeMutation && request.Schema.Mutation != nil {
        objectFields = request.Schema.GetMutationObjectFields()
    } else if operationType == frontend.OperationTypeSubscription && request.Schema.Subscription != nil {
        objectFields = request.Schema.GetSubscriptionObjectFields()
    } else {
        err = errors.New("Execute(): request.Schema should have Query or Mutation or Subscription field, please check server side Schema definition.")
        result.SetErrorInfo(err, nil)
        return &result
    }
    fmt.Printf("\033[33m    [DUMP] objectFields:  \033[0m\n")
    spewo.Dump(objectFields)

    // execute
    fmt.Println("\n\n\033[33m////////////////////////////////////////// Executor Start ///////////////////////////////////////\033[0m\n")
    var resolvedResult interface{}
    if resolvedResult, err = resolveSelectionSet(g, request, selectionSet, objectFields, nil); err != nil {
        result.SetErrorInfo(err, nil)
        return &result
    }
    fmt.Printf("\033[33m    [DUMP] resolvedResult:  \033[0m\n")
    spewo.Dump(resolvedResult)
    result.Data = resolvedResult
    return &result
}




// get name mapped Fields from SelectionSet
func getSelectionSetFields(selectionSet *frontend.SelectionSet) (map[string]*frontend.Field, error) {
    fields := make(map[string]*frontend.Field)

    selections := selectionSet.GetSelections()

    for _, selection := range selections {
        field := selection.(*frontend.Field)
        fieldName := field.Name.Value
        fields[fieldName] = field
    }
    return fields, nil
}

func getSubObjectFields(objectField *ObjectField) ObjectFields {
    fmt.Printf("\033[31m[INTO] func getSubObjectFields():  \033[0m\n")

    fieldType := objectField.Type
    if _, ok := fieldType.(*List); ok {
        return objectField.Type.(*List).Payload.(*Object).Fields
    } 

    if _, ok := fieldType.(*Object); ok {
        return objectField.Type.(*Object).Fields
    }

    return nil
}

func resolveSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveSelectionSet  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}


    if selectionSet == nil {
        return nil, errors.New("resolveSelectionSet(): empty selectionSet input.")
    }

    selections  := selectionSet.GetSelections()
    finalResult := make(map[string]interface{}, len(selections))

    fmt.Printf("\033[33m    [DUMP] objectFields:  \033[0m\n")
    spewo.Dump(objectFields)

    // resolve SelectionSet.Selections
    var resolvedResult interface{}
    var err            error
    for _, selection := range selections {
        field     := selection.(*frontend.Field)
        fieldName := field.GetFieldNameString()
        // resolve
        if resolvedResult, err = resolveField(g, request, fieldName, field, objectFields, resolvedData); err != nil {
            return nil, err
        }
        finalResult[fieldName] = resolvedResult   
    }
    return finalResult, nil
}

func defaultValueTypeAssertion(value interface{}) (interface{}, error) {
    // notice: the DefaultValue only accept const Value (Variables are not const Value)
    if _, ok := value.(frontend.Variable); ok {
        return nil, errors.New("defaultValueTypeAssertion(): the DefaultValue only accept const Value (Variables are not const Value).")
    } else if ret, ok := value.(frontend.IntValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.FloatValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.StringValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.BooleanValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.NullValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.EnumValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.ListValue); ok {
        return ret.Value, nil
    } else if ret, ok := value.(frontend.ObjectValue); ok {
        return ret.Value, nil
    } else {
        return nil, errors.New("defaultValueTypeAssertion(): illegal default value type '"+reflect.TypeOf(value).Elem().Name()+"'.")
    }
}

func correctJsonUnmarshalIntValue(value interface{}, variableType frontend.Type) (int, error) {
    // Int!
    if  _, ok := variableType.(frontend.NonNullType); ok {
        if val, ok := variableType.(frontend.NonNullType).Type.(*frontend.NamedType); ok {
            if val.Value == Int.Name {
                return int(value.(float64)), nil
            }
        }
    // Int
    } else if val, ok := variableType.(*frontend.NamedType); ok {
        if val.Value == Int.Name {
            return int(value.(float64)), nil
        }
    }
    // not an int at all
    return 0, errors.New("correctJsonUnmarshalIntValue(): not a IntValue.")
}

// build QueryVariables map from user input request.Variables
func getQueryVariablesMap(request Request, variableDefinitions []*frontend.VariableDefinition) (map[string]interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func getQueryVariablesMap  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("\033[33m    [DUMP] variableDefinitions:  \033[0m\n")
    spewo.Dump(variableDefinitions)

    var err error
    queryVariablesMap := make(map[string]interface{}, len(variableDefinitions))
    
    for _, variableDefinition := range variableDefinitions {
        // detect value type & fill
        variableName := variableDefinition.Variable.Value
        variableType := variableDefinition.Type
        if matchedValue, ok := request.Variables[variableName]; ok {
            // convert float64 to int type for json.Unmarshal
            if intValue, err := correctJsonUnmarshalIntValue(matchedValue, variableType); err == nil {
                queryVariablesMap[variableName] = intValue
            } else{
                queryVariablesMap[variableName] = matchedValue
            }
        // check NonNullType
        } else if _, ok := variableType.(frontend.NonNullType); ok {
            typeStr := ""
            if val, ok := variableType.(frontend.NonNullType).Type.(*frontend.NamedType); ok {
                typeStr = val.Value
            } else {
                typeStr = reflect.TypeOf(variableType.(frontend.NonNullType).Type).Elem().Name()
            }
            return nil, errors.New("getQueryVariablesMap(): variable '"+variableName+"' is NonNullType '"+typeStr+"!', query variables not provided.")
        // check DefaultValue
        } else if variableDefinition.DefaultValue != nil {
            if queryVariablesMap[variableName], err = defaultValueTypeAssertion(variableDefinition.DefaultValue); err != nil {
                return nil, err
            }
        }
    }
    
    fmt.Printf("\033[33m    [DUMP] queryVariablesMap:  \033[0m\n")
    spewo.Dump(queryVariablesMap)

    return queryVariablesMap, nil
}

// build Field.Arguments map from GlobalVariables.QueryVariablesMap
func getFieldArgumentsMap(g *GlobalVariables, arguments []*frontend.Argument) (map[string]interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func getFieldArgumentsMap  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fieldArgumentsMap := make(map[string]interface{}, len(arguments))
    
    for _, argument := range arguments {
        // detect argument type & fill
        argumentName  := argument.Name.Value
        argumentValue := argument.Value
        // assert Argument.Value type
        if _, ok := argumentValue.(frontend.Variable); ok {
            // Variable type, resolve referenced value from GlobalVariables.QueryVariablesMap
            if matched, ok := g.QueryVariablesMap[argumentName]; ok {
                fieldArgumentsMap[argumentName] = matched
            } else {
                err := "getFieldArgumentsMap(): Field.Arguments referenced variable $"+argumentName+", but it was NOT defined at OperationDefinition.VariableDefinitions, please check your GraphQL OperationDefinition syntax."
                return nil, errors.New(err)
            }
        } else if val, ok := argumentValue.(frontend.IntValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.FloatValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.StringValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.BooleanValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.NullValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.EnumValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.ListValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else if val, ok := argumentValue.(frontend.ObjectValue); ok {
            fieldArgumentsMap[argumentName] = val.Value
        } else {
            err := "getFieldArgumentsMap(): Field.Arguments.Argument type assert failed, please check your GraphQL Field.Arguments.Argument syntax."
            return nil, errors.New(err)
        }
    }
    
    fmt.Printf("\033[33m    [DUMP] fieldArgumentsMap:  \033[0m\n")
    spewo.Dump(fieldArgumentsMap)

    return fieldArgumentsMap, nil
}

func checkIfInputArgumentsAvaliable(inputArguments map[string]interface{}, targetObjectFieldArguments *Arguments) (bool, error) {
    for argumentName, _ := range inputArguments {
        if _, ok := (*targetObjectFieldArguments)[argumentName]; !ok {
            err := "checkIfInputArgumentsAvaliable(): input argument '"+argumentName+"' does not defined in schema."
            return false, errors.New(err)
        }
    }
    return true, nil
}



func resolveField(g *GlobalVariables, request Request, fieldName string, field *frontend.Field, objectFields ObjectFields, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveField  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
    fmt.Printf("\033[33m    [DUMP] fieldName:  \033[0m\n")
    spewo.Dump(fieldName)
    fmt.Printf("\033[33m    [DUMP] field:  \033[0m\n")
    spewo.Dump(field)
    fmt.Printf("\033[33m    [DUMP] objectFields:  \033[0m\n")
    spewo.Dump(objectFields)

    var err error

    if _, ok := objectFields[fieldName]; !ok {
        err := "resolveField(): input document field name '"+fieldName+"' does not match schema."
        return nil, errors.New(err)
    }
    
    
    // resolve
    if schemaResolveFunctionAvaliable(fieldName, objectFields) {
        // execute schema resolve function
        if resolvedData, err = schemaResolveFunction(g, request, fieldName, field, objectFields, resolvedData); err != nil {
            return nil, err
        }
    } 

    if resolvedData, err = defaultResolveFunction(g, request, field.SelectionSet, objectFields[fieldName], resolvedData); err != nil {
        return nil, err
    }

    return resolvedData, nil
}

func schemaResolveFunctionAvaliable(fieldName string, objectFields ObjectFields) bool {
    if objectFields[fieldName].ResolveFunction != nil {
        return true
    }
    return false
}

func schemaResolveFunction(g *GlobalVariables, request Request, fieldName string, field *frontend.Field, objectFields ObjectFields, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\033[31m[INTO] func schemaResolveFunction():  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}


    // build resolve params for resolve function
    var resolveParams ResolveParams
    var err           error
    resolveParams.Source = resolvedData
    if resolveParams.Arguments, err = getFieldArgumentsMap(g, field.Arguments); err != nil {
        return nil, err
    }

    // get resolve function
    resolveFunction := objectFields[fieldName].ResolveFunction

    // execute
    if resolvedData, err = resolveFunction(resolveParams); err != nil {
        return nil, err
    }

    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    return resolvedData, err
}



func resolvedDataTypeChecker(fieldName string, resolvedData interface{}, expectedType FieldType) (bool, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolvedDataTypeChecker  \033[0m\n")
    errorInfo := func(fieldName string, expected string, but string) error {
        err := "resolvedDataTypeChecker(): schema defined ObjectField '"+fieldName+"' Type is '"+expected+"', but ResolveFunction return type is '"+but+"', please check your schema."
        return errors.New(err)
    }
    resolvedDataType := reflect.TypeOf(resolvedData)
    switch resolvedDataType.Kind() {
        case reflect.Slice:
            if _, ok := expectedType.(*List); ok {
                return true, nil
            }
            return false, errorInfo(fieldName, reflect.TypeOf(expectedType).Elem().Name(), "slice, array or map")
        case reflect.Array:
            if _, ok := expectedType.(*List); ok {
                return true, nil
            }
            return false, errorInfo(fieldName, reflect.TypeOf(expectedType).Elem().Name(), "slice, array or map")
        case reflect.Map:
            if _, ok := expectedType.(*List); ok {
                return true, nil
            }
            return false, errorInfo(fieldName, reflect.TypeOf(expectedType).Elem().Name(), "slice, array or map")
        case reflect.Struct:
            if _, ok := expectedType.(*Object); ok {
                return true, nil
            }
            return false, errorInfo(fieldName, reflect.TypeOf(expectedType).Elem().Name(), "struct")
        default:
            if _, ok := expectedType.(*Scalar); ok {
                return true, nil
            }
    }
    return false, errorInfo(fieldName, reflect.TypeOf(expectedType).Elem().Name(), reflect.TypeOf(resolvedData).Elem().Name())
}


func defaultResolveFunction(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func defaultResolveFunction  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
    fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
    spewo.Dump(objectField)
    fmt.Printf("\033[33m    [DUMP] targetType:  \033[0m\n")
    targetType := objectField.Type
    spewo.Dump(targetType)
    // get resolve target type

    if _, ok := targetType.(*Scalar); ok {
        return resolveScalarData(g, request, selectionSet, objectField, resolvedData)
    }

    if _, ok := targetType.(*List); ok {
        return resolveListData(g, request, selectionSet, objectField, resolvedData)
    } 

    if _, ok := targetType.(*Object); ok {
        return resolveObjectData(g, request, selectionSet, objectField, resolvedData)
    }
    return nil, errors.New("defaultResolveFunction(): can not resolve target field.")
    
}

func resolveScalarData(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveScalarData  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
    spewo.Dump(selectionSet)
    fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
    spewo.Dump(objectField)
    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    // call resolve function
    targetFieldName := objectField.Name
    r0 := getResolvedDataByFieldName(targetFieldName, resolvedData)
    fmt.Printf("\033[33m    [DUMP] r0 result:  \033[0m\n")
    spewo.Dump(resolvedData)
    fmt.Printf("\033[33m    [DUMP] targetFieldName result:  \033[0m\n")
    spewo.Dump(targetFieldName)
    fmt.Printf("\033[33m    [DUMP] getResolvedDataByFieldName result:  \033[0m\n")
    spewo.Dump(r0)
    return r0, nil
    // // convert 
    // resolveFunction := objectField.Type.(*Scalar).ResolveFunction
    // p := ResolveParams{}
    // p.Context = r0
    // r1, _ := resolveFunction(p)
    // fmt.Printf("\033[43;37m    [DUMP] resolveFunction result:  \033[0m\n")
    // spewo.Dump(r1)
    // return r1, nil
}

func resolveListData(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveListData  \033[0m\n")
    // spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    resolvedDataValue := reflect.ValueOf(resolvedData)
    targetObjectFields := objectField.Type.(*List).Payload.(*Object).Fields

    // allocate space for list data returns
    finalResult := make([]interface{}, 0, resolvedDataValue.Len())

    // traverse list
    for i:=0; i<resolvedDataValue.Len(); i++ {
        resolvedDataElement := resolvedDataValue.Index(i).Interface()
        // fmt.Printf("\033[33m    [DUMP] resolvedDataElement:  \033[0m\n")
        // spewo.Dump(resolvedDataElement)
        // fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
        // spewo.Dump(objectField)
        // fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
        // spewo.Dump(selectionSet)
        // execute
        selectionSetResult, _ := resolveSelectionSet(g, request, selectionSet, targetObjectFields, resolvedDataElement)
        finalResult = append(finalResult, selectionSetResult)
    }
    return finalResult, nil
}

func resolveObjectData(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveObjectData  \033[0m\n")

    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
    spewo.Dump(selectionSet)
    fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
    spewo.Dump(objectField)
    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    // check if object type schema need default resolve function to get data
    // @todo: add a check method for situations that can be ignored
    r0 := getResolvedDataByFieldName(objectField.Name, resolvedData)
    fmt.Printf("\033[33m    [DUMP] r0:  \033[0m\n")
    spewo.Dump(r0)
    if r0 != nil {
        resolvedData = r0
    }


    // go
    targetObjectFields := objectField.Type.(*Object).Fields
    selectionSetResult, _ := resolveSelectionSet(g, request, selectionSet, targetObjectFields, resolvedData)
    return selectionSetResult, nil
}

func getResolvedDataByJsonTag(targetFieldName string, resolvedData interface{}) (interface{}) {
    fmt.Printf("\033[31m[INTO] func getResolvedDataByJsonTag  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("\033[33m    [DUMP] targetFieldName:  \033[0m\n")
    spewo.Dump(targetFieldName)

    fmt.Printf("\033[33m    [DUMP] targetFieldName:  \033[0m\n")
    spewo.Dump(resolvedData)

    val := reflect.ValueOf(resolvedData)

    fmt.Printf("\033[33m    [DUMP] val:  \033[0m\n")
    spewo.Dump(val)

    for i := 0; i < val.Type().NumField(); i++ {
        if val.Type().Field(i).Tag.Get("json") == targetFieldName {
            return reflect.Indirect(val).FieldByName(val.Type().Field(i).Name)
        }
    }
    return nil
}

func getResolvedDataByFieldName(targetFieldName string, resolvedData interface{}) (interface{}) {
    fmt.Printf("\033[31m[INTO] func getResolvedDataByFieldName  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("\033[33m    [DUMP] targetFieldName:  \033[0m\n")
    spewo.Dump(targetFieldName)

    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    val := reflect.ValueOf(resolvedData)

    for i := 0; i < val.Type().NumField(); i++ {
        if val.Type().Field(i).Name == targetFieldName {
            return reflect.Indirect(val).FieldByName(val.Type().Field(i).Name).Interface()
        }
    }
    return nil
}


// types

type Type interface {
    GetName() string
}

type FieldType interface {
    GetName() string
}

// List types

type List struct {
    Payload Type `json:payload`
}

func (list *List) GetName() string {
    return fmt.Sprintf("%v", list.Payload)
}

func NewList(i Type) *List {
    list := &List{}

    if i == nil {
        log.Fatal("NewList() input is nil")
        return list
    }

    list.Payload = i
    return list
}

func (list *List)GetDescription() string {
    return fmt.Sprintf("A List of '%v'.", list.Payload)
}

func (list *List) ToString() string {
    if list.Payload != nil {
        return fmt.Sprintf("[%v]", list.Payload)
    }
    return ""
}

// scalar definition

type ScalarTemplate struct {
    Name            string          `json:name`
    Description     string          `json:description`
    ResolveFunction ResolveFunction `json:"-"`
}

type Scalar struct {
    Name            string          `json:name`
    Description     string          `json:description`
    ResolveFunction ResolveFunction `json:"-"`
}

func (scalar *Scalar) GetName() string {
    return scalar.Name
}


func NewScalar(scalarTemplate ScalarTemplate) *Scalar {
    scalar := &Scalar{}

    // check scalar template
    if scalarTemplate.Name == "" {
        err := "scalarTemplate.Name is not defined"
        log.Fatal(err)
    }

    scalar.Name            = scalarTemplate.Name
    scalar.Description     = scalarTemplate.Description
    scalar.ResolveFunction = scalarTemplate.ResolveFunction

    // @todo: Scalar should provide serialize, parse value, parse literal function.
    
    return scalar
}


// scalar types
var Int = NewScalar(ScalarTemplate{
    Name: "Int",
    Description: "GraphQL Int type",
    ResolveFunction: func (p ResolveParams) (interface{}, error) {
        return p.Context.(reflect.Value).Int(), nil
    },
})

var String = NewScalar(ScalarTemplate{
    Name: "String",
    Description: "GraphQL String type",
    ResolveFunction: func (p ResolveParams) (interface{}, error) {
        return p.Context.(reflect.Value).String(), nil
    },
})

var Bool = NewScalar(ScalarTemplate{
    Name: "Bool",
    Description: "GraphQL Bool type",
    ResolveFunction: func (p ResolveParams) (interface{}, error) {
        return p.Context.(reflect.Value).Bool(), nil
    },
})

var Float = NewScalar(ScalarTemplate{
    Name: "Float",
    Description: "GraphQL Float type",
    ResolveFunction: func (p ResolveParams) (interface{}, error) {
        return p.Context.(reflect.Value).Float(), nil
    },
})




// Object Syntax

type ObjectFields map[string]*ObjectField

type ObjectTemplate struct {
    Name   string 
    Fields ObjectFields
}

type Object struct {
    Name   string
    Fields ObjectFields
}

func (object *Object) GetName() string {
    return object.Name
}

func (object *Object) GetFields() ObjectFields {
    return object.Fields
} 

type ObjectField struct {
    Name            string               `json:name`
    Type            FieldType            `json:type`  // maybe call this returnType?
    Description     string               `json:description`
    Arguments       *Arguments           `json:arguments`    
    ResolveFunction ResolveFunction      `json:"-"`
}

func (objectField * ObjectField) GetName() string {
    return objectField.Name
}

type Arguments map[string]*Argument

type Argument struct {
    Name string    `json:name` 
    Type FieldType `json:type`
}

type ResolveFunction func(p ResolveParams) (interface{}, error)

// resolve params for ResolveFunction()
type ResolveParams struct {
    // resolved source data for user defined resolve function 
    Source interface{}

    // context from executor
    Context interface{}

    // arguments map from request
    Arguments map[string]interface{}
}


func NewObject(objectTemplate ObjectTemplate) (*Object, error) {
    object := &Object{}

    // check object input
    if objectTemplate.Name == "" {
        err := errors.New("ObjectTemplate.Name is not defined")
        return nil, err
    }
    
    object.Name = objectTemplate.Name
    object.Fields = objectTemplate.Fields
    return object, nil
}

// Schema Syntax

type SchemaTemplate struct {
    Query        *Object
    Mutation     *Object 
    Subscription *Object 
}

type Schema struct {
    Query        *Object
    Mutation     *Object 
    Subscription *Object 
}


func (schema *Schema) GetQueryObject() *Object {
    return schema.Query
}

func (schema *Schema) GetQueryObjectFields() ObjectFields {
    return schema.Query.Fields
}

func (schema *Schema) GetMutationObject() *Object {
    return schema.Mutation
}

func (schema *Schema) GetMutationObjectFields() ObjectFields {
    return schema.Mutation.Fields
}

func (schema *Schema) GetSubscriptionObject() *Object {
    return schema.Subscription
}

func (schema *Schema) GetSubscriptionObjectFields() ObjectFields {
    return schema.Subscription.Fields
}

func NewSchema(schemaTemplate SchemaTemplate) (Schema, error) {
    schema := Schema{}

    // fill schema
    schema.Query = schemaTemplate.Query
    schema.Mutation = schemaTemplate.Mutation
    schema.Subscription = schemaTemplate.Subscription


    return schema, nil
}