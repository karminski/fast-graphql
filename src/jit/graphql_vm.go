package jit

import (
	// "fmt"
	"log"

	
)



func newAssembler() *Assembler {
	// @todo: auto generate page size
	buf, e := Alloc(PageSize)
	if e != nil {
		log.Fatalf("alloc: %v", e.Error())
	}
	return &Assembler{buf, 0, GoABI}
}




func Execute(tape *Tape) string {

	// init
	stringifier := NewStringifier()
	assembler := newAssembler()
	defer Release(assembler.Buf)
	var pc int 
	pc = 0


CallBuildObjectStart := func() {
	stringifier.buildObjectStart()
	pc ++
}

CallBuildObjectEnd := func() {
	stringifier.buildObjectEnd()
	pc ++
}

CallBuildArrayStart := func() {
	stringifier.buildArrayStart()
	pc ++
}

CallBuildArrayEnd := func() {
	stringifier.buildArrayEnd()
	pc ++
}

CallBuildComma := func() {
	stringifier.buildComma()
	pc ++
}

CallBuildColon := func() {
	stringifier.buildColon()
	pc ++
}

CallBuildNull := func() {
	stringifier.buildNull()
	pc ++
}

CallBuildBool := func() {
	stringifier.buildBool((*tape)[pc].Operand.(bool))
	pc ++
}

CallBuildInt := func() {
	stringifier.buildInt((*tape)[pc].Operand.(int))
	pc ++
}

CallBuildUint64 := func() {
	stringifier.buildUint64((*tape)[pc].Operand.(uint64))
	pc ++
}

CallBuildUint32 := func() {
	stringifier.buildUint32((*tape)[pc].Operand.(uint32))
	pc ++
}

CallBuildInt64 := func() {
	stringifier.buildInt64((*tape)[pc].Operand.(int64))
	pc ++
}

CallBuildInt32 := func() {
	stringifier.buildInt32((*tape)[pc].Operand.(int32))
	pc ++
}

CallBuildFloat64 := func() {
	stringifier.buildFloat64((*tape)[pc].Operand.(float64))
	pc ++
}

CallBuildFloat32 := func() {
	stringifier.buildFloat32((*tape)[pc].Operand.(float32))
	pc ++
}

CallBuildString := func() {
	stringifier.buildString((*tape)[pc].Operand.(string))
	pc ++
}

CallBuildEmptyString := func() {
	stringifier.buildEmptyString()
	pc ++
}

CallBuildFieldName := func() {	
	stringifier.buildFieldName((*tape)[pc].Operand.(string))
	pc ++
}

	// run tape
		
	for _, inst := range *tape {
		switch inst.Opcode {
		case OP_BUILD_OBJECT_START: 
			assembler.CallFunc(CallBuildObjectStart)
		case OP_BUILD_OBJECT_END: 
			assembler.CallFunc(CallBuildObjectEnd)
		case OP_BUILD_ARRAY_START: 
			assembler.CallFunc(CallBuildArrayStart)
		case OP_BUILD_ARRAY_END: 
			assembler.CallFunc(CallBuildArrayEnd)
		case OP_BUILD_COMMA: 
			assembler.CallFunc(CallBuildComma)
		case OP_BUILD_COLON: 
			assembler.CallFunc(CallBuildColon)
		case OP_BUILD_NULL: 
			assembler.CallFunc(CallBuildNull)
		case OP_BUILD_BOOL: 
			assembler.CallFunc(CallBuildBool)
		case OP_BUILD_INT: 
			assembler.CallFunc(CallBuildInt)
		case OP_BUILD_UINT64: 
			assembler.CallFunc(CallBuildUint64)
		case OP_BUILD_UINT32: 
			assembler.CallFunc(CallBuildUint32)
		case OP_BUILD_INT64: 
			assembler.CallFunc(CallBuildInt64)
		case OP_BUILD_INT32: 
			assembler.CallFunc(CallBuildInt32)
		case OP_BUILD_FLOAT64: 
			assembler.CallFunc(CallBuildFloat64)
		case OP_BUILD_FLOAT32: 
			assembler.CallFunc(CallBuildFloat32)
		case OP_BUILD_STRING: 
			assembler.CallFunc(CallBuildString)
		case OP_BUILD_EMPTY_STRING: 
			assembler.CallFunc(CallBuildEmptyString)
		case OP_BUILD_FIELD_NAME: 
			assembler.CallFunc(CallBuildFieldName)
		}
	}

	// assenbler end
	assembler.Ret()

	// debug
	tape.Dump()
	assembler.Dump()

	// emit
	var f func()
	assembler.BuildTo(&f)

	// execute
	f()

	// get result

	return stringifier.Stringify()
}

