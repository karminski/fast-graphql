// lexer.go

package frontend

import (
    "strings"
    "regexp"
    "fmt"
    "bytes"
    "strconv"
    "errors"
)

// token const
const (
    TOKEN_EOF         = iota  // end-of-file
    TOKEN_NOT_NULL            // !
    TOKEN_VAR_PREFIX          // $
    TOKEN_LEFT_PAREN          // (
    TOKEN_RIGHT_PAREN         // )
    TOKEN_LEFT_BRACKET        // [
    TOKEN_RIGHT_BRACKET       // ]
    TOKEN_LEFT_BRACE          // {
    TOKEN_RIGHT_BRACE         // }
    TOKEN_COLON               // :
    TOKEN_DOTS                // ...
    TOKEN_EQUAL               // =
    TOKEN_AT                  // @
    TOKEN_AND                 // &
    TOKEN_VERTICAL_BAR        // |
	TOKEN_QUOTE         	  // "
	TOKEN_DUOQUOTE 			  // ""
	TOKEN_TRIQUOTE      	  // """
	TOKEN_HEXQUOTE 			  // """"""
    TOKEN_ESCAPE_CHARACTER    // \

    // literal

    TOKEN_NUMBER              // number literal
    TOKEN_IDENTIFIER          // identifier
    // Name
    // TOKEN_NAME                // all name, include OperationName, TypeName, VariableName, ArgumentName etc.
    // Int
    // TOKEN_INT                 // value int 
    // TOKEN_FLOAT               // value float

    // Float
    // String
    // BlockString
    
    // keywords
    TOKEN_QUERY                 // query                
    TOKEN_FRAGMENT              // fragment
    TOKEN_MUTATION              // mutation
    TOKEN_SUBSCRIPTION          // subscription
    TOKEN_TYPE                  // type
    TOKEN_INTERFACE             // interface
    TOKEN_UNION                 // union
    TOKEN_SCHEMA                // schema
    TOKEN_ENUM                  // enum
    TOKEN_INPUT                 // input
    TOKEN_DIRECTIVE             // directive
    TOKEN_IMPLEMENTS            // implements
    TOKEN_EXTEND                // extend
    TOKEN_SCALAR                // scalar
    TOKEN_TRUE                  // true
    TOKEN_FALSE                 // false
    TOKEN_NULL                  // null
    TOKEN_ON                    // on
)

var tokenNameMap = map[int]string{
    TOKEN_EOF              : "EOF",
    TOKEN_NOT_NULL         : "!",
    TOKEN_VAR_PREFIX       : "$",
    TOKEN_LEFT_PAREN       : "(",
    TOKEN_RIGHT_PAREN      : ")",    
    TOKEN_LEFT_BRACKET     : "[",
    TOKEN_RIGHT_BRACKET    : "]",    
    TOKEN_LEFT_BRACE       : "{",
    TOKEN_RIGHT_BRACE      : "}",    
    TOKEN_COLON            : ":",
    TOKEN_DOTS             : "...",
    TOKEN_EQUAL            : "=",
    TOKEN_AT               : "@",
    TOKEN_AND              : "&",
    TOKEN_VERTICAL_BAR     : "|",
	TOKEN_QUOTE            : "\"",
	TOKEN_DUOQUOTE 		   : "\"\"",
	TOKEN_TRIQUOTE         : "\"\"\"",
	TOKEN_HEXQUOTE 		   : "\"\"\"\"\"\"",
    TOKEN_ESCAPE_CHARACTER : "\\",
    TOKEN_NUMBER           : "number",
    TOKEN_IDENTIFIER       : "identifier",
    TOKEN_QUERY            : "query",
    TOKEN_FRAGMENT         : "fragment",
    TOKEN_MUTATION         : "mutation",
    TOKEN_SUBSCRIPTION     : "subscription",
    TOKEN_TYPE             : "type",
    TOKEN_INTERFACE        : "interface",
    TOKEN_UNION            : "union",
    TOKEN_SCHEMA           : "schema",
    TOKEN_ENUM             : "enum",
    TOKEN_INPUT            : "input",
    TOKEN_DIRECTIVE        : "directive",
    TOKEN_IMPLEMENTS       : "implements",
    TOKEN_EXTEND           : "extend",
    TOKEN_SCALAR           : "scalar",
    TOKEN_TRUE             : "true",
    TOKEN_FALSE            : "false",
    TOKEN_NULL             : "null",
    TOKEN_ON               : "on",
}

var keywords = map[string]int{
    "query"        : TOKEN_QUERY,
    "fragment"     : TOKEN_FRAGMENT,
    "mutation"     : TOKEN_MUTATION,
    "subscription" : TOKEN_SUBSCRIPTION,
    "type"         : TOKEN_TYPE,
    "interface"    : TOKEN_INTERFACE,
    "union"        : TOKEN_UNION,
    "schema"       : TOKEN_SCHEMA,
    "enum"         : TOKEN_ENUM,
    "input"        : TOKEN_INPUT,
    "directive"    : TOKEN_DIRECTIVE,
    "implements"   : TOKEN_IMPLEMENTS,
    "extend"       : TOKEN_EXTEND,
    "scalar"       : TOKEN_SCALAR,
    "true"         : TOKEN_TRUE,
    "false"        : TOKEN_FALSE,
    "null"         : TOKEN_NULL,
    "on"           : TOKEN_ON,
}

var avaliableNumberParts = map[byte]bool{
    '0': true,
    '1': true,
    '2': true,
    '3': true,
    '4': true,
    '5': true,
    '6': true,
    '7': true,
    '8': true,
    '9': true,
    '-': true,
    'x': true,
    'X': true, 
    '.': true,
    '+': true,
    'p': true, 
    'P': true, 
    'a': true,
    'b': true, 
    'c': true, 
    'd': true, 
    'e': true, 
    'f': true, 
    'A': true,
    'B': true, 
    'C': true, 
    'D': true, 
    'E': true, 
    'F': true,
}

// regex match patterns
var regexNumber        = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[-]?[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reDecEscapeSeq     = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq     = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

// lexer struct
type Lexer struct {
    document            string // graphql document
    lineNum             int    // current line number
    nextToken           string 
    nextTokenType       int 
    nextTokenLineNum    int
    pos                 int    // now position of document 
}

func NewLexer(document string) *Lexer {
    return &Lexer{document, 1, "", 0, 0, 0} // start at line 1 in default.
}

func (lexer *Lexer) GetLineNum() int {
    return lexer.lineNum
}

func (lexer *Lexer) GetPos() int {
    return lexer.pos
}

func (lexer *Lexer) NextTokenIs(tokenType int) (lineNum int, token string) {

    nowLineNum, nowTokenType, nowToken := lexer.GetNextToken()
    // syntax error
    if tokenType != nowTokenType {
        fmt.Println("\n\n\033[05m\033[41;37m                    OOOOOOOOOPS! TOKEN EXCEPT FAILED                    \033[0m\n")
        err := fmt.Sprintf("line %d: syntax error near '%s'.", lexer.GetLineNum(), nowToken) 
        fmt.Println("- [Lexer Dump] -------------------------------------------------------\n")
        fmt.Printf("document:\n\n\033[33m%v\033[0m\n\n", lexer.document)
        fmt.Printf("lineNum:          \033[33m%v\033[0m\n", lexer.lineNum)
        fmt.Printf("NextTokenIs:      \033[33m%v\033[0m\n", tokenNameMap[tokenType])
        fmt.Printf("nowToken:         \033[33m%v\033[0m\n", nowToken)
        fmt.Printf("nextToken:        \033[33m%v\033[0m\n", lexer.nextToken)
        fmt.Printf("nextTokenType:    \033[33m%v\033[0m\n", lexer.nextTokenType)
        fmt.Printf("nextTokenLineNum: \033[33m%v\033[0m\n", lexer.nextTokenLineNum)
        fmt.Printf("\n---------------------------------------------------------------------\n\n")
        panic(err)
    }
    return nowLineNum, nowToken
}

func (lexer *Lexer) LookAhead() int {
    // lexer.nextToken* already setted
    if lexer.nextTokenLineNum > 0 {
        return lexer.nextTokenType
    }
    // set it
    nowLineNum                := lexer.lineNum
    lineNum, tokenType, token := lexer.GetNextToken()
    
    lexer.lineNum              = nowLineNum
    lexer.nextTokenLineNum     = lineNum
    lexer.nextTokenType        = tokenType
    lexer.nextToken            = token
    return tokenType
}

func (lexer *Lexer) nextDocumentIs(s string) bool {
	return len(lexer.document) >= len(s) && lexer.document[0:len(s)] == s
}

func (lexer *Lexer) skipDocument(n int) {
    lexer.pos += n
    lexer.document = lexer.document[n:]
}

// target pattern
func isNewLine(c byte) bool {
    return c == '\r' || c == '\n'
}

func isWhiteSpace(c byte) bool {
    switch c {
    case '\t', '\n', '\v', '\f', '\r', ' ':
        return true
    }
    return false
}

func isComma(c byte) bool {
    switch c {
    case ',':
        return true
    }
    return false
}

func isComment(c byte) bool {
    switch c {
    case '#':
        return true
    }
    return false
}

func (lexer *Lexer) skipIgnored() {
    // matching
    for len(lexer.document) > 0 {
        if lexer.nextDocumentIs("\r\n") || lexer.nextDocumentIs("\n\r") {
            lexer.skipDocument(2)
            lexer.lineNum += 1
        } else if isNewLine(lexer.document[0]) {
            lexer.skipDocument(1)
            lexer.lineNum += 1
        } else if isWhiteSpace(lexer.document[0]) {
            lexer.skipDocument(1)
        } else if isComma(lexer.document[0]) {
            lexer.skipDocument(1)
        } else if isComment(lexer.document[0]) {
            lexer.skipDocument(1)
            for !isNewLine(lexer.document[0]) {
                lexer.skipDocument(1)
            }
        } else {
            break
        } 
    }
}

// use regex scan for number, identifier 
func (lexer *Lexer) scan(regexp *regexp.Regexp) string {
    if token := regexp.FindString(lexer.document); token != "" {
        lexer.skipDocument(len(token))
        return token
    }
    panic("unreachable!")
    return ""
}

// return content before token
func (lexer *Lexer) scanBeforeToken(token string) string {
    s := strings.Split(lexer.document, token)
    if len(s) < 2 {
        panic("unreachable!")
        return ""
    }
    lexer.skipDocument(len(s[0]))
    return s[0]
}

// NOTE: this method can skip escape character 
func (lexer *Lexer) scanBeforeByte(b byte) (string, error) {
    docLen    := len(lexer.document)
    var r string
    var err error
    i := 0
    for ; i < docLen; i++ {
        // hit target
        if lexer.document[i] == b {
            r = lexer.document[:i]
            // convert escape character
            if r, err = lexer.escape(r); err != nil {
                return "", err
            }
            lexer.skipDocument(i)
            return r, nil
        }
        // skip escape character
        if lexer.document[i] == '\\' {
            i +=2
        }
    }
    return "", errors.New("Can not find target byte.")
}

func (lexer *Lexer) scanNumber() string {
    docLen := len(lexer.document)
    for i := 0; i< docLen; i++ {
        c := lexer.document[i] 
        if _, ok := avaliableNumberParts[c]; ok {
            continue
        } else {
            target := lexer.document[:i]
            lexer.skipDocument(i)
            return target
        }
    }
    panic("unreachable!")
    return ""
}

func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
    return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}


func (lexer *Lexer) scanIdentifier() string {
    docLen := len(lexer.document)
    for i := 0; i< docLen; i++ {
        c := lexer.document[i] 
        if c == '_' || isLetter(c) || isDigit(c) {
            continue
        } else {
            target := lexer.document[:i]
            lexer.skipDocument(i)
            return target
        }
    }
    panic("unreachable!")
    return ""
}


func (lexer *Lexer) GetNextToken() (lineNum int, tokenType int, token string) {
    // next token already loaded
    if lexer.nextTokenLineNum > 0 {
        lineNum                = lexer.nextTokenLineNum
        tokenType              = lexer.nextTokenType
        token                  = lexer.nextToken
        lexer.lineNum          = lexer.nextTokenLineNum
        lexer.nextTokenLineNum = 0
        return
    }
    return lexer.MatchToken()

}


func (lexer *Lexer) MatchToken() (lineNum int, tokenType int, token string) {
    // skip spaces
    lexer.skipIgnored()
    // finish
    if len(lexer.document) == 0 {
        return lexer.lineNum, TOKEN_EOF, tokenNameMap[TOKEN_EOF]
    }
    // check token
    switch lexer.document[0] {
    case '!' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_NOT_NULL, "!"
    case '$' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_VAR_PREFIX, "$"
    case '(' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_LEFT_PAREN, "("
    case ')' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_RIGHT_PAREN, ")"
    case '[' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_LEFT_BRACKET, "["
    case ']' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_RIGHT_BRACKET, "]"
    case '{' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_LEFT_BRACE, "{"
    case '}' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_RIGHT_BRACE, "}"
    case ':' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_COLON, ":"
    case '.' :
        if lexer.nextDocumentIs("...") {
            lexer.skipDocument(3)
            return lexer.lineNum, TOKEN_DOTS, "..."
        }
    case '=' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_EQUAL, "="
    case '@' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_AT, "@"
    case '&' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_AND, "&"
    case '|' :
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_VERTICAL_BAR, "|"
    case '"' :
        if lexer.nextDocumentIs("\"\"\"\"\"\"") {
            lexer.skipDocument(6)
            return lexer.lineNum, TOKEN_HEXQUOTE, "\"\"\"\"\"\""
        }
        if lexer.nextDocumentIs("\"\"\"") {
            lexer.skipDocument(3)
            return lexer.lineNum, TOKEN_TRIQUOTE, "\"\"\""
        }
        if lexer.nextDocumentIs("\"\"") {
            lexer.skipDocument(2)
            return lexer.lineNum, TOKEN_DUOQUOTE, "\"\""
        }
        lexer.skipDocument(1)
        return lexer.lineNum, TOKEN_QUOTE, "\""
    }

    // check multiple character token
    if lexer.document[0] == '_' || isLetter(lexer.document[0]) {
        token := lexer.scanIdentifier()
        if tokenType, isMatch := keywords[token]; isMatch {
            return lexer.lineNum, tokenType, token
        } else {
            return lexer.lineNum, TOKEN_IDENTIFIER, token
        }
    }
    if lexer.document[0] == '.' || lexer.document[0] == '-' || isDigit(lexer.document[0]) {
        token := lexer.scanNumber()
        return lexer.GetLineNum(), TOKEN_NUMBER, token
    }

    // unexpected symbol
    err := fmt.Sprintf("line %d: unexpected symbol near '%q'.", lexer.lineNum, lexer.document[0])
    panic(err)
    return 
}

// this escape method copy from https://github.com/zxh0/lua.go/blob/master/compiler/lexer/lexer.go
// @copyright zxh0
// @license MIT
func (lexer *Lexer) escape(str string) (string, error) {
    var buf bytes.Buffer

    for len(str) > 0 {
        if str[0] != '\\' {
            buf.WriteByte(str[0])
            str = str[1:]
            continue
        }

        if len(str) == 1 {
            return "", errors.New("unfinished string")
        }

        switch str[1] {
        case 'a':
            buf.WriteByte('\a')
            str = str[2:]
            continue
        case 'b':
            buf.WriteByte('\b')
            str = str[2:]
            continue
        case 'f':
            buf.WriteByte('\f')
            str = str[2:]
            continue
        case 'n', '\n':
            buf.WriteByte('\n')
            str = str[2:]
            continue
        case 'r':
            buf.WriteByte('\r')
            str = str[2:]
            continue
        case 't':
            buf.WriteByte('\t')
            str = str[2:]
            continue
        case 'v':
            buf.WriteByte('\v')
            str = str[2:]
            continue
        case '"':
            buf.WriteByte('"')
            str = str[2:]
            continue
        case '\'':
            buf.WriteByte('\'')
            str = str[2:]
            continue
        case '\\':
            buf.WriteByte('\\')
            str = str[2:]
            continue
        case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
            if found := reDecEscapeSeq.FindString(str); found != "" {
                d, _ := strconv.ParseInt(found[1:], 10, 32)
                if d <= 0xFF {
                    buf.WriteByte(byte(d))
                    str = str[len(found):]
                    continue
                }
                return "", errors.New(fmt.Sprintf("decimal escape too large near '%s'", found))
            }
        case 'x': // \xXX
            if found := reHexEscapeSeq.FindString(str); found != "" {
                d, _ := strconv.ParseInt(found[2:], 16, 32)
                buf.WriteByte(byte(d))
                str = str[len(found):]
                continue
            }
        case 'u': // \u{XXX}
            if found := reUnicodeEscapeSeq.FindString(str); found != "" {
                d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
                if err == nil && d <= 0x10FFFF {
                    buf.WriteRune(rune(d))
                    str = str[len(found):]
                    continue
                }
                return "", errors.New(fmt.Sprintf("UTF-8 value too large near '%s'", found))
            }
        case 'z':
            str = str[2:]
            for len(str) > 0 && isWhiteSpace(str[0]) { // todo
                str = str[1:]
            }
            continue
        }
        return "", errors.New(fmt.Sprintf("invalid escape sequence near '\\%c'", str[1]))
    }

    return buf.String(), nil
}