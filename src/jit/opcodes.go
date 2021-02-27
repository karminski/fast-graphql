package jit

import (
	"fmt"
)

const (
	// symbol
	OP_BUILD_OBJECT_START = iota
	OP_BUILD_OBJECT_END
	OP_BUILD_ARRAY_START
	OP_BUILD_ARRAY_END
	OP_BUILD_COMMA
	OP_BUILD_COLON
	
	// null
	OP_BUILD_NULL
	
	// boolean
	OP_BUILD_BOOL       
	
	// number
	OP_BUILD_INT        
	OP_BUILD_UINT64     
	OP_BUILD_UINT32     
	OP_BUILD_INT64      
	OP_BUILD_INT32      
	OP_BUILD_FLOAT64    
	OP_BUILD_FLOAT32    
	
	// string
	OP_BUILD_STRING       
	OP_BUILD_EMPTY_STRING
	
	// combined
	OP_BUILD_FIELD_NAME 
)

type Opcode     int
type JITOperand interface{}



type JITInstruction struct {
	Opcode  Opcode
	Operand JITOperand
}

type Tape []JITInstruction


func NewTape() *Tape {
	return new(Tape)
}

func (tape *Tape) Record(opcode Opcode, operand JITOperand) {
	inst := JITInstruction{opcode, operand}

	*tape = append(*tape, inst)
}

func (tape *Tape) Dump() {
	for _, inst := range *tape {
		fmt.Printf("%d, %v\n", inst.Opcode, inst.Operand)
	}
}