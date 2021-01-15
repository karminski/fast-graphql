package backend

import (
	"unsafe"


	"fmt"
    "github.com/davecgh/go-spew/spew"
)


// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
    typ  *rtype
    word  unsafe.Pointer
}

// tflag is used by an rtype to signal what extra type information is
// available in the memory directly following the rtype value.
//
// tflag values must be kept in sync with copies in:
//  cmd/compile/internal/gc/reflect.go
//  cmd/link/internal/ld/decodesym.go
//  runtime/type.go
type tflag uint8

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype

// rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
type rtype struct {
    size       uintptr
    ptrdata    uintptr // number of bytes in the type that can contain pointers
    hash       uint32  // hash of type; avoids computation in hash tables
    tflag      tflag   // extra type information flags
    align      uint8   // alignment of variable with this type
    fieldAlign uint8   // alignment of struct field with this type
    kind       uint8   // enumeration for C
    // function for comparing objects of this type
    // (ptr to object A, ptr to object B) -> ==?
    equal     func(unsafe.Pointer, unsafe.Pointer) bool
    gcdata    *byte   // garbage collection data
    str       nameOff // string form
    ptrToThis typeOff // type for pointer to this type, may be zero
}

type flag uintptr
const flagMethod  flag = 1 << 9 // @ go/srcreflect/value.go
const structKind = 25

// structType represents a struct type.
type structType struct {
    rtype
    pkgPath name
    fields  []structField // sorted by offset
}

// name is an encoded type name with optional extra data.
//
// The first byte is a bit field containing:
//
//  1<<0 the name is exported
//  1<<1 tag data follows the name
//  1<<2 pkgPath nameOff follows the name and tag
//
// The next two bytes are the data length:
//
//   l := uint16(data[1])<<8 | uint16(data[2])
//
// Bytes [3:3+l] are the string data.
//
// If tag data follows then bytes 3+l and 3+l+1 are the tag length,
// with the data following.
//
// If the import path follows, then 4 bytes at the end of
// the data form a nameOff. The import path is only set for concrete
// methods that are defined in a different package than their type.
//
// If a name starts with "*", then the exported bit represents
// whether the pointed to type is exported.
type name struct {
    bytes *byte
}

type UnsafeString struct {
    Data unsafe.Pointer
    Len  int
}

func (n name) name() (s string) {
    if n.bytes == nil {
        return
    }
    b := (*[4]byte)(unsafe.Pointer(n.bytes))

    hdr := (*UnsafeString)(unsafe.Pointer(&s))
    hdr.Data = unsafe.Pointer(&b[3])
    hdr.Len = int(b[1])<<8 | int(b[2])
    return s
}

// Struct field
type structField struct {
    name        name    // name is always non-empty
    typ         *rtype  // type of field
    offsetEmbed uintptr // byte offset of field<<1 | isEmbedded
}



func ResolveByFieldName(structData interface{}, name string) interface{} {
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("[INTO] ResolveByFieldName():\n")
    spewo.Dump(structData)
	spewo.Dump(name)

	// unpack Struct
    e := (*emptyInterface)(unsafe.Pointer(&structData))
    fmt.Printf("e:\n")
    spewo.Dump(e)


    startPointer := e.word

    // check flag
    valueMetadataFlag := e.typ.kind
    if valueMetadataFlag == 0 {
        panic("invalied reflect.Value.Type")
    }

    // if valueMetadataFlag&flagMethod == 0 {
    //     // Easy case
    //     fmt.Printf("[Easy case!]\n")
    // }
    if valueMetadataFlag != structKind {
        panic("reflect: Field of non-struct type ")
    }

    // pick up target field
    var targetStructField *structField
    tt := (*structType)(unsafe.Pointer(e.typ))
    for i := range tt.fields {
        tf := &tt.fields[i]
        if tf.name.name() == name {
            targetStructField = tf
            break
        }
    }
    if targetStructField == nil {
    	return nil
    }

    fmt.Printf("targetStructField:\n")
    spewo.Dump(targetStructField)

    // repack target field
    var packed interface{}
    pe := (*emptyInterface)(unsafe.Pointer(&packed))
    pe.typ  = targetStructField.typ
    pe.word = unsafe.Pointer(uintptr(startPointer) + targetStructField.typ.ptrdata)
    fmt.Printf("pe.word:\n")
    spewo.Dump(pe.word)
    return packed
}