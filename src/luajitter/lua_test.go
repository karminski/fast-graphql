package luajitter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func closeVM(t *testing.T, vm *LuaState) {
	err := vm.Close()
	require.Nil(t, err)
}
func TestSimpleGlobal(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.SetGlobal("wow", 2)
	require.Nil(t, err)

	val, err := vm.GetGlobal("wow")
	require.Nil(t, err)
	require.NotNil(t, val)

	number, ok := val.(float64)
	require.True(t, ok)
	require.Equal(t, 2.0, number)
}

func TestInitGlobal(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.InitGlobal("test.test2.test3", "value")
	require.Nil(t, err)

	val, err := vm.GetGlobal("test.test2.test3")
	require.Nil(t, err)
	require.NotNil(t, val)
	strVal, ok := val.(string)
	require.True(t, ok)
	require.Equal(t, "value", strVal)

	tableObj, err := vm.GetGlobal("test.test2")
	require.Nil(t, err)
	require.NotNil(t, tableObj)
	table, ok := tableObj.(*LocalLuaTable)
	require.True(t, ok)
	require.NotNil(t, table)
	require.Equal(t, 5, int(table.value.valueType))

	err = table.Close()
	require.Nil(t, err)
}

func TestIsYieldable(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.DoString(`
	function canyield()
		return coroutine.isyieldable()
	end
`)
	require.Nil(t, err)

	funcObj, err := vm.GetGlobal("canyield")
	require.Nil(t, err)
	f := funcObj.(*LocalLuaFunction)
	outVal, err := f.Call()
	require.Nil(t, err)
	require.Equal(t, false, outVal[0])

	err = f.Close()
	require.Nil(t, err)
}

func TestInitGlobalNumIndex(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.InitGlobal("test.test2.1", "This ")
	require.Nil(t, err)

	err = vm.SetGlobal("test.test2.2", "is")
	require.Nil(t, err)

	err = vm.SetGlobal("test.test2.3", " good")
	require.Nil(t, err)

	err = vm.DoString("function getVal() return table.concat(test.test2) end")
	require.Nil(t, err)

	funcObj, err := vm.GetGlobal("getVal")
	require.Nil(t, err)

	f := funcObj.(*LocalLuaFunction)
	outVal, err := f.Call()
	require.Nil(t, err)
	require.Len(t, outVal, 1)
	require.Equal(t, "This is good", outVal[0])
	err = f.Close()
	require.Nil(t, err)
}

const fibo string = `
function fib(val)
	if val < 2 then 
		return val
	end

	return fib(val-2) + fib(val-1)
end

print(fib(5))
`

func TestDoStringAndCall(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.DoString(fibo)
	require.Nil(t, err)

	val, err := vm.GetGlobal("fib")
	require.Nil(t, err)
	require.NotNil(t, val)

	fibFunc, ok := val.(*LocalLuaFunction)
	require.True(t, ok)
	require.NotNil(t, fibFunc)

	out, err := fibFunc.Call(7)
	require.Nil(t, err)
	require.NotNil(t, out)
	require.Len(t, out, 1)
	require.NotNil(t, out[0])

	outNumber, ok := out[0].(float64)
	require.True(t, ok)
	require.Equal(t, 13.0, outNumber)

	err = fibFunc.Close()
	require.Nil(t, err)
}

const multiRet string = `
function multiCall()
	return 9,"testing",false
end
`

func TestMultiRetCall(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	require := require.New(t)

	err := vm.DoString(multiRet)
	require.Nil(err)

	val, err := vm.GetGlobal("multiCall")
	require.Nil(err)
	require.NotNil(val)

	multiCallFunc, ok := val.(*LocalLuaFunction)
	require.True(ok)
	require.NotNil(multiCallFunc)

	out, err := multiCallFunc.Call()
	require.Nil(err)
	require.NotNil(out)
	require.Len(out, 3)
	require.Equal(9.0, out[0])
	require.Equal("testing", out[1])
	require.Equal(false, out[2])

	err = multiCallFunc.Close()
	require.Nil(err)
}

func TestDoStringAndCallNil(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.DoString("function retNil() return nil end")
	require.Nil(t, err)

	val, err := vm.GetGlobal("retNil")
	require.Nil(t, err)
	require.NotNil(t, val)

	fibFunc, ok := val.(*LocalLuaFunction)
	require.True(t, ok)
	require.NotNil(t, fibFunc)

	out, err := fibFunc.Call()
	require.Nil(t, err)
	require.NotNil(t, out)
	require.Len(t, out, 1)
	require.Nil(t, out[0])

	err = fibFunc.Close()
	require.Nil(t, err)
}

func TestDoStringWithError(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.DoString(`error("some error")`)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "some error")
}

func TestDoCallWithError(t *testing.T) {
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(t, 0, outlyingAllocs())
	}()

	err := vm.DoString(`function errcall(msg) error(msg) end`)
	require.Nil(t, err)

	val, err := vm.GetGlobal("errcall")
	require.Nil(t, err)
	require.NotNil(t, val)

	fibFunc, ok := val.(*LocalLuaFunction)
	require.True(t, ok)
	require.NotNil(t, fibFunc)

	out, err := fibFunc.Call("another error")
	require.NotNil(t, err)
	require.Len(t, out, 0)
	require.Contains(t, err.Error(), "another error")

	err = fibFunc.Close()
	require.Nil(t, err)
}

var callbackArgs []interface{}

func SomeErrorCallback(args []interface{}) ([]interface{}, error) {
	callbackArgs = args
	return []interface{}{
		"test",
		5,
		true,
		SomeErrorCallback,
	}, errors.New("WOW ERROR")
}
func TestDoErrorCallback(t *testing.T) {
	require := require.New(t)
	clearAllocs()
	vm := NewState()
	defer func() {
		closeVM(t, vm)
		require.Equal(0, outlyingAllocs())
	}()

	err := vm.InitGlobal("test.error_callback", SomeErrorCallback)
	require.Nil(err)

	err = vm.DoString(`
function doErrorCallback()
	return test.error_callback(5, "bleh", nil, {})
end
`)
	require.Nil(err)

	errorFuncObj, err := vm.GetGlobal("doErrorCallback")
	require.Nil(err)
	require.NotNil(errorFuncObj)

	errorF, ok := errorFuncObj.(*LocalLuaFunction)
	require.True(ok)
	require.NotNil(errorF)

	retVals, err := errorF.Call()
	require.Len(retVals, 0)

	require.NotNil(err)
	require.Equal("WOW ERROR", err.Error())

	require.Len(callbackArgs, 4)
	require.Equal(5.0, callbackArgs[0])
	require.Equal("bleh", callbackArgs[1])
	require.Nil(callbackArgs[2])
	require.IsType(&LocalLuaTable{}, callbackArgs[3])

	data := callbackArgs[3].(*LocalLuaTable)
	err = data.Close()
	require.Nil(err)

	err = errorF.Close()
	require.Nil(err)
}

func SomeCallback(args []interface{}) ([]interface{}, error) {
	callbackArgs = args
	return []interface{}{
		"test",
		5,
		true,
		SomeCallback,
	}, nil
}
func TestDoCallback(t *testing.T) {
	require := require.New(t)
	clearAllocs()
	vm := NewState()

	err := vm.InitGlobal("test.callback", SomeCallback)
	require.Nil(err)

	err = vm.DoString(`
function doCallback()
	return test.callback(5, "bleh", nil, {})
end
`)
	require.Nil(err)

	funcObj, err := vm.GetGlobal("doCallback")
	require.Nil(err)
	require.NotNil(funcObj)

	f, ok := funcObj.(*LocalLuaFunction)
	require.True(ok)
	require.NotNil(f)

	retVals, err := f.Call()
	require.NotNil(retVals)
	require.Nil(err)
	require.Len(retVals, 4)
	require.Equal("test", retVals[0])
	require.Equal(5.0, retVals[1])
	require.Equal(true, retVals[2])
	require.IsType(&LocalLuaData{}, retVals[3])

	data := retVals[3].(LocalData)
	err = data.Close()
	require.Nil(err)

	require.Len(callbackArgs, 4)
	require.Equal(5.0, callbackArgs[0])
	require.Equal("bleh", callbackArgs[1])
	require.Nil(callbackArgs[2])
	require.IsType(&LocalLuaTable{}, callbackArgs[3])

	tableRet := callbackArgs[3].(*LocalLuaTable)
	err = tableRet.Close()
	require.Nil(err)

	err = f.Close()
	require.Nil(err)

	closeVM(t, vm)

	require.Equal(0, outlyingAllocs())
}

func TestPermanentValues(t *testing.T) {
	require := require.New(t)
	clearAllocs()
	vm := NewState()

	err := vm.SetGlobal("sometable", map[interface{}]interface{}{})
	require.Nil(err)

	someTable, err := vm.GetGlobal("sometable")
	require.Nil(err)

	err = vm.DoString(`
function somefunc()
end
`)
	require.Nil(err)

	someFunc, err := vm.GetGlobal("somefunc")
	require.Nil(err)

	err = vm.SetGlobal("mytable", map[interface{}]interface{}{
		"t": someTable,
		"f": someFunc,
	})
	require.Nil(err)

	tableVal, err := vm.GetGlobal("mytable")
	require.Nil(err)
	require.IsType(&LocalLuaTable{}, tableVal)
	table := tableVal.(*LocalLuaTable)
	unrolled, err := table.Unroll()
	require.Nil(err)

	require.IsType(make(map[interface{}]interface{}), unrolled["t"])
	require.IsType(&LocalLuaFunction{}, unrolled["f"])

	err = unrolled["f"].(*LocalLuaFunction).Close()
	require.Nil(err)

	err = table.Close()
	require.Nil(err)
	err = someTable.(*LocalLuaTable).Close()
	require.Nil(err)
	err = someFunc.(*LocalLuaFunction).Close()
	require.Nil(err)
	closeVM(t, vm)

	require.Equal(0, outlyingAllocs())
}

func TestTwoTierUnroll(t *testing.T) {
	require := require.New(t)
	clearAllocs()
	vm := NewState()

	err := vm.SetGlobal("mytable", map[interface{}]interface{}{
		"testValue": "wow!",
		"testNumber": 5,
		"innerTable": map[interface{}]interface{}{
			"innerValue": "neat",
			"someBool": true,
		},
	})
	require.Nil(err)

	tableVal, err := vm.GetGlobal("mytable")
	require.Nil(err)
	require.IsType(&LocalLuaTable{}, tableVal)
	table := tableVal.(*LocalLuaTable)
	unrolled, err := table.Unroll()
	require.Nil(err)

	require.Equal("wow!", unrolled["testValue"])
	require.Equal(5.0, unrolled["testNumber"])
	require.IsType(make(map[interface{}]interface{}), unrolled["innerTable"])
	innerTable := unrolled["innerTable"].(map[interface{}]interface{})
	require.Equal("neat", innerTable["innerValue"])
	require.Equal(true, innerTable["someBool"])

	err = table.Close()
	require.Nil(err)
	closeVM(t, vm)

	require.Equal(0, outlyingAllocs())
}

func TestRollUnroll(t *testing.T) {
	require := require.New(t)
	clearAllocs()
	vm := NewState()

	err := vm.SetGlobal("mytable", map[interface{}]interface{}{
		"testValue": "wow!",
		"testNumber": 5,
		"innerTable": map[interface{}]interface{}{
			"innerValue": "neat",
			"someBool": true,
		},
	})
	require.Nil(err)

	strValue, err := vm.GetGlobal("mytable.testValue")
	require.Nil(err)
	require.IsType("", strValue)
	nowStr := strValue.(string)
	require.Equal("wow!", nowStr)

	numValue, err := vm.GetGlobal("mytable.testNumber")
	require.Nil(err)
	require.IsType(0.0, numValue)
	nowNum := numValue.(float64)
	require.Equal(5.0, nowNum)

	tableValue, err := vm.GetGlobal("mytable.innerTable")
	require.Nil(err)
	require.IsType(&LocalLuaTable{}, tableValue)
	nowTable := tableValue.(*LocalLuaTable)
	table, err := nowTable.Unroll()
	require.Nil(err)

	nowStr = table["innerValue"].(string)
	require.Equal("neat", nowStr)
	nowBool := table["someBool"].(bool)
	require.Equal(true, nowBool)

	err = nowTable.Close()
	require.Nil(err)
	closeVM(t, vm)

	require.Equal(0, outlyingAllocs())
}

func BenchmarkFib35(b *testing.B) {
	clearAllocs()
	vm := NewState()
	defer vm.Close()

	err := vm.DoString(`
function fib(val)
	if val < 2 then 
		return val
	end

	return fib(val-2) + fib(val-1)
end
`)
	if err != nil {
		panic(err)
	}

	funcObj, err := vm.GetGlobal("fib")
	if err != nil {
		panic(err)
	}

	f := funcObj.(*LocalLuaFunction)
	out, err := f.Call(35)
	if err != nil {
		panic(err)
	}
	fmt.Println(out[0])
}

func BenchmarkRandom(b *testing.B) {
	clearAllocs()
	vm := NewState()
	defer vm.Close()

	err := vm.DoString(`
	local rand = math.random
	math.randomseed(os.time())
	rand(); rand(); rand()
function uuid()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and rand(0, 0xf) or rand(8, 0xb)
        return string.format('%x', v)
    end)
end
`)
	if err != nil {
		panic(err)
	}
	
	uuidFuncObj, err := vm.GetGlobal("uuid")
	if err != nil {
		panic(err)
	}

	uuidFunc := uuidFuncObj.(*LocalLuaFunction)
	uuid, err := uuidFunc.Call()
	fmt.Println("ITERATION ",b.N, uuid[0])
}


var cbCount = 0

func AddCallback(args []interface{}) ([]interface{}, error) {
	cbCount++
	if len(args) != 2 {
		return nil, errors.New("incorrect arguments passed to Add")
	}

	l, ok := args[0].(float64)
	if !ok {
		return nil, errors.New("argument 1 to Add was not a number")
	}
	r, ok := args[1].(float64)
	if !ok {
		return nil, errors.New("argument 2 to Add was not a number")
	}

	return []interface{}{
		l + r,
	}, nil
}

func BenchmarkRollUnroll(b *testing.B) {
	clearAllocs()
	vm := NewState()
	defer vm.Close()

	in := make(map[interface{}]interface{})
	for i := 0; i < 100; i++ {
		in[i+1] = map[interface{}]interface{}{
			"field1": "a",
			"field2": "b",
			"field3": 3,
		}
	}

	err := vm.SetGlobal("bigmap", in)
	if err != nil {
		panic(err)
	}

	table, err := vm.GetGlobal("bigmap")
	if err != nil {
		panic(err)
	}

	_, err = table.(*LocalLuaTable).Unroll()
	if err != nil {
		panic(err)
	}

	err = table.(LocalData).Close()
	if err != nil {
		panic(err)
	}
}

func BenchmarkCallbackFib35(b *testing.B) {
	clearAllocs()
	vm := NewState()
	defer vm.Close()

	err := vm.SetGlobal("_Add", AddCallback)
	if err != nil {
		panic(err)
	}

	err = vm.DoString(`
function fib(val)
	if val < 2 then 
		return val
	end

	return _Add(fib(val-2), fib(val-1))
end
`)
	if err != nil {
		panic(err)
	}

	funcObj, err := vm.GetGlobal("fib")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	f := funcObj.(*LocalLuaFunction)
	out, err := f.Call(35)
	if err != nil {
		panic(err)
	}
	fmt.Println(cbCount)
	fmt.Println(out[0])
}
