package backend

import (
    "errors"
    "sync"
    "fmt"

    "fast-graphql/src/luajitter"
    "fast-graphql/src/frontend"


    _ "net/http"
    _ "io"
    _ "net/http/pprof"

)


func callSechemaResolveFunctionByJIT(fieldName string) {


    
}

func NewResolveFunctionMap() map[interface{}]interface{} {
	var fmap map[interface{}]interface{}
	fmap = make(map[interface{}]interface{})
	return fmap
}

func buildSchemaResolveFunctionMap(objectFields ObjectFields, fmap map[interface{}]interface{})  {	
	for fieldName, objectField := range objectFields {
		if objectField.ResolveFunction != nil {
			fmap[fieldName] = objectField.ResolveFunction
		}
		if subFields := getSubObjectFields(objectField); subFields != nil {
			buildSchemaResolveFunctionMap(subFields, fmap)
		}
	} 
}

func packArguments(queryVariablesMap map[string]interface{}) ResolveParams {
	var p ResolveParams
	p.Arguments = queryVariablesMap
	return p
}



func steppingSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (string, error) {
	if selectionSet == nil {
        return "", errors.New("resolveSelectionSet(): empty selectionSet input.")
    }
    // selections  := selectionSet.GetSelections()
    
    // run JIT
    var mutex sync.Mutex
    mutex.Lock()

    code := `
function m0()
    -- init
    local buf = {}
    
    -- build data header
    buf[#buf+1] = "{\"data\":"

    -- get user data by param from [go runtime]
    local user = _G.User 
    
    -- resolve user
    resolveUser(user, buf)

    -- build no error data end
    buf[#buf+1] = ",\"errors\":null,\"jit-result\":true}"

    -- finally, return result to go runtime
    local r = table.concat(buf)
    _G.result = r
end

function resolveUser(user, buf)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"Id\":"
    buf[#buf+1] = user.Id
    buf[#buf+1] = ",\"Name\":"
    buf[#buf+1] = user.Name
    buf[#buf+1] = ",\"Email\":"
    buf[#buf+1] = user.Email
    buf[#buf+1] = ",\"Married\":"
    buf[#buf+1] = user.Married
    buf[#buf+1] = ",\"Height\":"
    buf[#buf+1] = user.Height
    buf[#buf+1] = ",\"Gender\":"
    buf[#buf+1] = user.Gender
    buf[#buf+1] = ",\"Friends\":"
    resolveFriends(user.Friends)
    buf[#buf+1] = ",\"Location\":"
    resolveLocation(user.Location)
    buf[#buf+1] = "}"
end

m0()
   	`
    // emit
    var err error
    var r interface{}
    vm := luajitter.NewState()
    err = vm.DoString(code)
    if err != nil {
        panic(err)
    }
    r, err = vm.GetGlobal("_G.result")
    if err != nil {
        panic(err)
    }
    fmt.Printf("-----------\n")
    fmt.Println(r)

    closeVM(vm)


    defer mutex.Unlock()

    // end
	return "", nil
}

func closeVM(vm *luajitter.LuaState) {
    err := vm.Close()
    if err != nil {
        panic(err)
    }
}