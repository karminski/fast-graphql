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
	var operand  JITOperand


CallBuildObjectStart := func() {
	stringifier.buildObjectStart()
}

CallBuildObjectEnd := func() {
	stringifier.buildObjectEnd()
}

CallBuildArrayStart := func() {
	stringifier.buildArrayStart()
}

CallBuildArrayEnd := func() {
	stringifier.buildArrayEnd()
}

CallBuildComma := func() {
	stringifier.buildComma()
}

CallBuildColon := func() {
	stringifier.buildColon()
}

CallBuildNull := func() {
	stringifier.buildNull()
}

CallBuildBool := func() {
	stringifier.buildBool(operand.(bool))
}

CallBuildInt := func() {
	stringifier.buildInt(operand.(int))
}

CallBuildUint64 := func() {
	stringifier.buildUint64(operand.(uint64))
}

CallBuildUint32 := func() {
	stringifier.buildUint32(operand.(uint32))
}

CallBuildInt64 := func() {
	stringifier.buildInt64(operand.(int64))
}

CallBuildInt32 := func() {
	stringifier.buildInt32(operand.(int32))
}

CallBuildFloat64 := func() {
	stringifier.buildFloat64(operand.(float64))
}

CallBuildFloat32 := func() {
	stringifier.buildFloat32(operand.(float32))
}

CallBuildString := func() {
	stringifier.buildString(operand.(string))
}

CallBuildEmptyString := func() {
	stringifier.buildEmptyString()
}

CallBuildFieldName := func() {
	stringifier.buildFieldName(operand.(string))
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
			operand = inst.Operand
			assembler.CallFunc(CallBuildBool)
		case OP_BUILD_INT: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildInt)
		case OP_BUILD_UINT64: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildUint64)
		case OP_BUILD_UINT32: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildUint32)
		case OP_BUILD_INT64: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildInt64)
		case OP_BUILD_INT32: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildInt32)
		case OP_BUILD_FLOAT64: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildFloat64)
		case OP_BUILD_FLOAT32: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildFloat32)
		case OP_BUILD_STRING: 
			operand = inst.Operand
			assembler.CallFunc(CallBuildString)
		case OP_BUILD_EMPTY_STRING: 
			assembler.CallFunc(CallBuildEmptyString)
		case OP_BUILD_FIELD_NAME: 
			operand = inst.Operand
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

