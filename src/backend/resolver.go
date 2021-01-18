package backend

import (
	"unsafe"

	// "reflect"
	// "fmt"
    // "github.com/davecgh/go-spew/spew"
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

const (
	invalidKind uint8 = iota
	boolKind
	intKind
	int8Kind
	int16Kind
	int32Kind
	int64Kind
	uintKind
	uint8Kind
	uint16Kind
	uint32Kind
	uint64Kind
	uintptrKind
	float32Kind
	float64Kind
	complex64Kind
	complex128Kind
	arrayKind
	chanKind
	funcKind
	interfaceKind
	mapKind
	ptrKind
	sliceKind
	stringKind
	structKind
	unsafePointerKind
)

// structType represents a struct type.
type structType struct {
    rtype
    pkgPath name
    fields  []structField // sorted by offset
}

// sliceType represents a slice type.
type sliceType struct {
	rtype
	elem *rtype // slice element type
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

// Slice is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
//
// Unlike reflect.SliceHeader, its Data field is sufficient to guarantee the
// data it references will not be garbage collected.
type UnsafeSlice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
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

func (f *structField) offset() uintptr {
	return f.offsetEmbed >> 1
}


// add returns p+x.
//
// The whySafe string is ignored, so that the function still inlines
// as efficiently as p+x, but all call sites should use the string to
// record why the addition is safe, which is to say why the addition
// does not cause x to advance to the very end of p's allocation
// and therefore point incorrectly at the next block in memory.
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// arrayAt returns the i-th element of p,
// an array whose elements are eltSize bytes wide.
// The array pointed at by p must have at least i+1 elements:
// it is invalid (but impossible to check here) to pass i >= len,
// because then the result will point outside the array.
// whySafe must explain why i < len. (Passing "i < len" is fine;
// the benefit is to surface this assumption at the call site.)
func arrayAt(p unsafe.Pointer, i int, eltSize uintptr, whySafe string) unsafe.Pointer {
	return add(p, uintptr(i)*eltSize, "i < len")
}



func ResolveByFieldName(structData interface{}, name string) interface{} {
	// unpack Struct
    e := (*emptyInterface)(unsafe.Pointer(&structData))
    startPointer := e.word

    // check flag
    valueMetadataFlag := e.typ.kind
    if valueMetadataFlag == 0 {
        panic("invalied reflect.Value.Type")
    }
    
    // struct input please
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
        }
    }
    if targetStructField == nil {
    	return nil
    }

    // repack target field
    var packed interface{}
    pe := (*emptyInterface)(unsafe.Pointer(&packed))
    pe.typ  = targetStructField.typ
    pe.word = add(startPointer, targetStructField.offset(), "same as non-reflect &v.field")
    return packed
}


func ResolveSliceAllElements(sliceData interface{}) []interface{} {
	// unpack Struct
    e := (*emptyInterface)(unsafe.Pointer(&sliceData))

    // check flag
    valueMetadataFlag := e.typ.kind
    if valueMetadataFlag == 0 {
        panic("invalied reflect.Value.Type")
    }

    // slice input please
    if valueMetadataFlag != sliceKind {
        panic("resolver: non-slice type ")
    }
    
    s := (*UnsafeSlice)(e.word)

    // check slice len and init a container for return 
    if s.Len < 1 { // an empty slice
    	return nil
    }

    // pickup and return
    allElements := make([]interface{}, s.Len)
    tt := (*sliceType)(unsafe.Pointer(e.typ))
    typ := tt.elem
    for i := 0; i < s.Len; i++ {
        element := arrayAt(s.Data, i, typ.size, "i < s.Len")
        var packed interface{}
    	pe := (*emptyInterface)(unsafe.Pointer(&packed))
    	pe.typ  = typ
    	pe.word = element
    	allElements[i] = packed
    }
    return allElements
}