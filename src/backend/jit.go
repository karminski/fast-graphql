package backend

import (
    "errors"
    "sync"

    "reflect"

    "fast-graphql/src/luajitter"
    "fast-graphql/src/frontend"

    "github.com/fatih/structs"
    "github.com/karminski/fastreflect"




    _ "net/http"
    _ "io"
    _ "net/http/pprof"

    // "github.com/davecgh/go-spew/spew"


)


func callSechemaResolveFunctionByJIT(fieldName string) {


    
}

func NewResolveFunctionMap() map[string]interface{} {
	var fmap map[string]interface{}
	fmap = make(map[string]interface{})
	return fmap
}

func buildSchemaResolveFunctionMap(objectFields ObjectFields, fmap map[string]interface{})  {	
	for fieldName, objectField := range objectFields {
		if objectField.ResolveFunction != nil {
			fmap[fieldName] = objectField.ResolveFunction
		}
		if subFields := getSubObjectFields(objectField); subFields != nil {
			buildSchemaResolveFunctionMap(subFields, fmap)
		}
	} 
}

func packResolveArguments(queryVariablesMap map[string]interface{}) ResolveParams {
	var p ResolveParams
	p.Arguments = queryVariablesMap
	return p
}

func packResolveSource(source interface{}) ResolveParams {
	var p ResolveParams
	p.Source = source
	return p
}

func callSechemaResolveFunction(fmap map[string]interface{}, field string, p ResolveParams) (interface{}, error) {
	if _, ok := fmap[field]; !ok {
		return nil, errors.New("resolveSelectionSet(): empty selectionSet input.")
	}
	return fmap[field].(ResolveFunction)(p)
}


func loadQueryVariablesMapToJIT(vm *luajitter.LuaState, vmap map[string]interface{}) error {
	var vmapCopy map[interface{}]interface{}
	vmapCopy = make(map[interface{}]interface{}, len(vmap))
	for k, v := range vmap {
		vmapCopy[k] = v
	}
	return vm.SetGlobal("QueryVariablesMap", vmapCopy)
}


// call pointer linked list for JIT context switch
type callPointer struct {
	prev *callPointer
	next *callPointer
	payload interface{}
}

type callPointerlist []callPointer





// main JIT resolve method
func steppingSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (string, error) {
	fmap := NewResolveFunctionMap()
    buildSchemaResolveFunctionMap(objectFields, fmap)

    // init JIT
    var err error
    var r interface{}
    var nowCallPointer *callPointer

    // run JIT
    var mutex sync.Mutex
    mutex.Lock()
    
    vm := luajitter.NewState()

    switchToNextCallPointer := func(args []interface{}) ([]interface{}, error) {
    	nowCallPointer = nowCallPointer.next
		return []interface{}{}, nil
    }

    switchToPrevCallPointer := func(args []interface{}) ([]interface{}, error) {
    	nowCallPointer = nowCallPointer.prev
		return []interface{}{}, nil
    }

    initCallPointer := func(payload interface{}) {
    	var callp callPointer
    	callp.payload = payload
    	nowCallPointer = &callp
    }

	callResolveFuncByArguments := func(args []interface{}) ([]interface{}, error) {
		var field string
		var p ResolveParams
		

		field = args[0].(string)
		p.Arguments = g.QueryVariablesMap

		// call resolve function
		if _, ok := fmap[field]; !ok {
			return nil, errors.New("callResolveFuncByArguments(): unknown field input.")
		}
		r, err := fmap[field].(ResolveFunction)(p)
		mappedData := structs.Map(r) // @WILL_CRASH!, when this data structure are not struct 
		if err != nil {
			return nil, err
		}

		// set to call pointer
		initCallPointer(r)
		
        // set to JIT
        err = vm.SetGlobal(field, mappedData)
        if err != nil {
            panic(err)
        }
		return []interface{}{}, nil
	} 

	callResolveFuncBySource := func(args []interface{}) ([]interface{}, error) {
		var field string
		var p ResolveParams

		field = args[0].(string)
		p.Source = nowCallPointer.payload

		// call resolve function
		if _, ok := fmap[field]; !ok {
			return nil, errors.New("callResolveFuncBySource(): unknown field input.")
		}
		r, err := fmap[field].(ResolveFunction)(p)
		if err != nil {
			return nil, err
		}

		// detect return type for JIT, @todo: move this method to JIT code for just-in-time generate
		rt := reflect.TypeOf(r) 
		switch rt.Kind() {
		case reflect.Slice:
			allElem := fastreflect.SliceAllElements(r)
			// convert to []interface{}
			var convertedResult []interface{}
			convertedResult = make([]interface{}, 0, len(allElem))
			for _, v := range allElem {
				convertedResult = append(convertedResult, structs.Map(v))
			}
			err = vm.SetGlobal(field, convertedResult)
		case reflect.Array:
		default:
			err = vm.SetGlobal(field, r)
		}

		if err != nil {
            panic(err)
        }


		return []interface{}{}, nil
	} 



	if selectionSet == nil {
        return "", errors.New("resolveSelectionSet(): empty selectionSet input.")
    }
    // selections  := selectionSet.GetSelections()
    
    

    code := `

function m0()
    -- init
    local buf = {}
    
    -- build data header
    buf[#buf+1] = "{\"data\":"

    -- get User data by param from [go runtime]
    callResolveFuncByArguments("User")

    -- resolve user
    resolveUser(buf, User)

    -- build no error data end
    buf[#buf+1] = ",\"errors\":null,\"jit-result\":true}"

    -- finally, return result to go runtime
    local r = table.concat(buf)
    _G.result = r
end

function resolveUser(buf, user)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"Id\":"
    buf[#buf+1] = user.Id
    buf[#buf+1] = ",\"Name\":\""
    buf[#buf+1] = user.Name
    buf[#buf+1] = "\",\"Email\":\""
    buf[#buf+1] = user.Email
    buf[#buf+1] = "\",\"Married\":"
    if user.Married then
    	buf[#buf+1] = "true"
    else 
    	buf[#buf+1] = "false"
    end
    buf[#buf+1] = ",\"Height\":"
    buf[#buf+1] = user.Height
    buf[#buf+1] = ",\"Gender\":\""
    buf[#buf+1] = user.Gender
    buf[#buf+1] = "\",\"Friends\":"
    
    -- get User.Friends data by param from [go runtime]
    callResolveFuncBySource("Friends")

    -- resolve Friends
    local Friends = Friends
    resolveFriends(buf, Friends)
    
    buf[#buf+1] = ",\"Location\":"
    resolveLocation(buf, user.Location)
    buf[#buf+1] = "}"
end

function resolveFriends(buf, friends)
    l = #friends
    buf[#buf+1] = "["
    for i=0, l do
        resolveFriend(buf, friends[i])
        if i < l then
            buf[#buf+1] = ","
        end
    end
    buf[#buf+1] = "]"
end

function resolveFriend(buf, friend)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"Id\":"
    buf[#buf+1] = friend.Id
    buf[#buf+1] = ",\"Name\":\""
    buf[#buf+1] = friend.Name
    buf[#buf+1] = "\",\"Email\":\""
    buf[#buf+1] = friend.Email
    buf[#buf+1] = "\",\"Married\":"
    if friend.Married then
    	buf[#buf+1] = "true"
    else 
    	buf[#buf+1] = "false"
    end
    buf[#buf+1] = ",\"Height\":"
    buf[#buf+1] = friend.Height
    buf[#buf+1] = ",\"Gender\":\""
    buf[#buf+1] = friend.Gender
    buf[#buf+1] = "\",\"Location\":"
    resolveLocation(buf, friend.Location)
    buf[#buf+1] = "}"
end

function resolveLocation(buf, location)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"City\":\""
    buf[#buf+1] = location.City
    buf[#buf+1] = "\",\"Country\":\""
    buf[#buf+1] = location.Country
    buf[#buf+1] = "\"}"
end

m0()
   	`
    

    // preprocess for global variables and functions
    loadQueryVariablesMapToJIT(vm, g.QueryVariablesMap)
    err = vm.SetGlobal("callResolveFuncByArguments", callResolveFuncByArguments)
    if err != nil {
        panic(err)
    }
    err = vm.SetGlobal("callResolveFuncBySource", callResolveFuncBySource)
    err = vm.SetGlobal("switchToNextCallPointer", switchToNextCallPointer)
    err = vm.SetGlobal("switchToPrevCallPointer", switchToPrevCallPointer)
    if err != nil {
        panic(err)
    }

    // run JIT
    err = vm.DoString(code)
    if err != nil {
        panic(err)
    }
    r, err = vm.GetGlobal("_G.result")
    if err != nil {
        panic(err)
    }

    closeVM(vm)


    defer mutex.Unlock()

    // end
	return r.(string), nil
}

func closeVM(vm *luajitter.LuaState) {
    err := vm.Close()
    if err != nil {
        panic(err)
    }
}