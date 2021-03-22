package jit

import (
    "strings"
    "strconv"
    "fmt"
)


type Stringifier struct {
    Builder strings.Builder
}


func NewStringifier() *Stringifier {
    var stringifier Stringifier
    stringifier.buildDataHeader()
    return &stringifier
}


func (s *Stringifier)join(strs ...string) {
    for _, str := range strs {
        s.Builder.WriteString(str)
    }
}

// return builded string result 
func (s *Stringifier)Stringify() string {
    return s.Builder.String()
}

// builders
func (s *Stringifier)buildObjectStart() {
    s.Builder.WriteRune('{')
}

func (s *Stringifier)buildObjectEnd() {
    s.Builder.WriteRune('}')
}

func (s *Stringifier)buildArrayStart() {
    s.Builder.WriteRune('[')
}

func (s *Stringifier)buildArrayEnd() {
    s.Builder.WriteRune(']')
}

func (s *Stringifier)buildComma() {
    s.Builder.WriteRune(',')
}

func (s *Stringifier)buildColon() {
    s.Builder.WriteRune(':')
}

// null
func (s *Stringifier)buildNull() {
    s.Builder.WriteString("null")
}

// boolean
func (s *Stringifier)buildBool(b bool) {
    if b {
        s.Builder.WriteString("true")
    } else {
        s.Builder.WriteString("false")
    }
}

// number
func (s *Stringifier)buildInt(i int) {
    s.Builder.WriteString(strconv.Itoa(i))
}

func (s *Stringifier)buildUint32(u uint32) {
    s.Builder.WriteString(strconv.FormatUint(uint64(u), 10))
}

func (s *Stringifier)buildUint64(u uint64) {
    s.Builder.WriteString(strconv.FormatUint(u, 10))
}

func (s *Stringifier)buildInt32(i int32) {
    s.Builder.WriteString(strconv.FormatInt(int64(i), 10))
}

func (s *Stringifier)buildInt64(i int64) {
    s.Builder.WriteString(strconv.FormatInt(i, 10))
}

func (s *Stringifier)buildFloat32(f float32) {
    // @todo: is this method needed?
    return
}

func (s *Stringifier)buildFloat64(f float64) {
    s.Builder.WriteString(strconv.FormatFloat(f, 'E', -1 ,64))
}

// string
func (s *Stringifier)buildString(str string) {
    s.join("\"", str, "\"")
}

func (s *Stringifier)buildEmptyString() {
    s.Builder.WriteString("\"\"")
}

// packed builder method
func (s *Stringifier)buildFieldName(field string) {
    s.buildString(field)
    s.buildColon()
}

func (s *Stringifier)buildStringField(field string, value string) {
    s.buildString(field)
    s.buildColon()
    s.buildString(value)
}

func (s *Stringifier)buildEmptyStringField(field string) {
    s.buildString(field)
    s.buildColon()
    s.buildEmptyString()
}

func (s *Stringifier)buildIntField(field string, value int) {
    s.buildString(field)
    s.buildColon()
    s.buildInt(value)
}

func (s *Stringifier)buildFloat64Field(field string, value float64) {
    s.buildString(field)
    s.buildColon()
    s.buildFloat64(value)
}

func (s *Stringifier)buildBoolField(field string, value bool) {
    s.buildString(field)
    s.buildColon()
    s.buildBool(value)
}

func (s *Stringifier)buildNullField(field string) {
    s.buildString(field)
    s.buildColon()
    s.buildNull()
}

func (s *Stringifier)buildScalar(scalar interface{}) {
    switch scalar.(type){
    case string:
        s.buildString(scalar.(string))
    case int:
        s.buildInt(scalar.(int))
    case float64:
        s.buildFloat64(scalar.(float64))
    case float32:
        s.buildFloat32(scalar.(float32))
    case bool:
        s.buildBool(scalar.(bool))
    case nil:
        s.buildNull()
    }

}


// others
func (s *Stringifier)buildDataHeader() {
    s.Builder.WriteString("{\"data\":")
}

func (s *Stringifier)buildDataBottom() {
    s.Builder.WriteString("}")
}

func (s *Stringifier)buildErrorInfo(err error) {
    s.join(",\"errors\":", fmt.Sprintf("%v", err), "}")
}

func (s *Stringifier)buildNoError() {
    s.join(",\"errors\":null,\"jit-result\":true}")
}