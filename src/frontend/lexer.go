// lexer.go

package frontend

import (
    "strings"
    "regexp"
    "fmt"
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
    TOKEN_EOF           : "EOF",
    TOKEN_NOT_NULL      : "!",
    TOKEN_VAR_PREFIX    : "$",
    TOKEN_LEFT_PAREN    : "(",
    TOKEN_RIGHT_PAREN   : ")",    
    TOKEN_LEFT_BRACKET  : "[",
    TOKEN_RIGHT_BRACKET : "]",    
    TOKEN_LEFT_BRACE    : "{",
    TOKEN_RIGHT_BRACE   : "}",    
    TOKEN_COLON         : ":",
    TOKEN_DOTS          : "...",
    TOKEN_EQUAL         : "=",
    TOKEN_AT            : "@",
    TOKEN_AND           : "&",
    TOKEN_VERTICAL_BAR  : "|",
	TOKEN_QUOTE         : "\"",
	TOKEN_DUOQUOTE 		: "\"\"",
	TOKEN_TRIQUOTE      : "\"\"\"",
	TOKEN_HEXQUOTE 		: "\"\"\"\"\"\"",
    TOKEN_NUMBER        : "number",
    TOKEN_IDENTIFIER    : "identifier",
    TOKEN_QUERY         : "query",
    TOKEN_FRAGMENT      : "fragment",
    TOKEN_MUTATION      : "mutation",
    TOKEN_SUBSCRIPTION  : "subscription",
    TOKEN_TYPE          : "type",
    TOKEN_INTERFACE     : "interface",
    TOKEN_UNION         : "union",
    TOKEN_SCHEMA        : "schema",
    TOKEN_ENUM          : "enum",
    TOKEN_INPUT         : "input",
    TOKEN_DIRECTIVE     : "directive",
    TOKEN_IMPLEMENTS    : "implements",
    TOKEN_EXTEND        : "extend",
    TOKEN_SCALAR        : "scalar",
    TOKEN_TRUE          : "true",
    TOKEN_FALSE         : "false",
    TOKEN_NULL          : "null",
    TOKEN_ON            : "on",
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

// regex match patterns
var regexNumber     = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[-]?[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var regexIdentifier = regexp.MustCompile(`^[_\d\w]+`)

// lexer struct
type Lexer struct {
    document            string // graphql document
    lineNum             int    // current line number
    nextToken           string 
    nextTokenType       int 
    nextTokenLineNum    int
}

func NewLexer(document string) *Lexer {
    return &Lexer{document, 1, "", 0, 0} // start at line 1 in default.
}

func (lexer *Lexer) GetLineNum() int {
    return lexer.lineNum
}

func (lexer *Lexer) NextTokenIs(tokenType int) (lineNum int, token string) {

    nowLineNum, nowTokenType, nowToken := lexer.GetNextToken()
    fmt.Printf("    lexer.NextTokenIs( expect:'%v'==>'%v' -> Got:'%v'==>'%v' )\n", tokenType, tokenNameMap[tokenType], nowTokenType, tokenNameMap[nowTokenType])
    // syntax error
    if tokenType != nowTokenType {
        fmt.Println("\n\n\033[05m\033[41;37m                    OOOOOOOOOPS! TOKEN EXCEPT FAILED                    \033[0m\n")
        err := fmt.Sprintf("line %d: syntax error near '%s'.", lexer.GetLineNum(), nowToken) 
        fmt.Println("- [Lexer Dump] -------------------------------------------------------\n")
        fmt.Printf("document:\n\n\033[33m%v\033[0m\n\n", lexer.document)
        fmt.Printf("lineNum:          \033[33m%v\033[0m\n", lexer.lineNum)
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
    // fmt.Printf("\033[36mLookAhead().lexer.nextTokenLineNum: %v\033[0m\n", lexer.nextTokenLineNum)
    if lexer.nextTokenLineNum > 0 {
    // fmt.Printf("\033[36mLookAhead().lexer.nextTokenType: %v : %v\033[0m\n", lexer.nextTokenType, tokenNameMap[lexer.nextTokenType])

        return lexer.nextTokenType
    }
    // set it
    nowLineNum                := lexer.lineNum
    // fmt.Printf("\033[36mLookAhead().lineNum: %v \033[0m\n", lexer.lineNum)

    lineNum, tokenType, token := lexer.GetNextToken()
    // fmt.Printf("\033[36mLookAhead().lineNum: %v \033[0m\n", lexer.lineNum)
    
    lexer.lineNum              = nowLineNum
    lexer.nextTokenLineNum     = lineNum
    lexer.nextTokenType        = tokenType
    lexer.nextToken            = token
    // fmt.Printf("\033[36mLookAhead().tokenType: %v : %v\033[0m\n", tokenType, tokenNameMap[tokenType])
    return tokenType
}

func (lexer *Lexer) nextDocumentIs(s string) bool {
    return strings.HasPrefix(lexer.document, s)
}

func (lexer *Lexer) skipDocument(n int) {
    lexer.document = lexer.document[n:]
}

func (lexer *Lexer) skipIgnored() {
    // target pattern
    isNewLine := func(c byte) bool {
        return c == '\r' || c == '\n'
    }
    isWhiteSpace := func(c byte) bool {
        switch c {
        case '\t', '\n', '\v', '\f', '\r', ' ':
            return true
        }
        return false
    }
    isComma := func(c byte) bool {
        switch c {
        case ',':
            return true
        }
        return false
    }
    isComment := func(c byte) bool {
        switch c {
        case '#':
            return true
        }
        return false
    }
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

func (lexer *Lexer) scanNumber() string {
    return lexer.scan(regexNumber)
}

func (lexer *Lexer) scanIdentifier() string {
    return lexer.scan(regexIdentifier)
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
    fmt.Printf("now lexer: %c\n", lexer.document[0])
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




func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
    return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

