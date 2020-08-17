// lexer.go
// 

package frontend

import (
    "strings"
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

    // literal

    TOKEN_NUMBER              // number literal
    TOKEN_IDENTIFIER          // identifier
    // Name
    // Int
    // Float
    // String
    // BlockString
    
    // keywords
    TOKEN_QUERY                 // query                
    TOKEN_FRAGMENT              // fragment
    TOKEN_MUTATION              // mutation
    TOKEN_TRUE                  // true
    TOKEN_FALSE                 // false
)

var keywords = map[string]int{
    "query":    TOKEN_QUERY,
    "fragment": TOKEN_FRAGMENT,
    "mutation": TOKEN_MUTATION,
}

// regex match patterns
var regexNumber     = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
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
    nowLineNum, nowTokenType, nowToken := lexer.NextToken()
    // syntax error
    if tokenType != nowTokenType {
        err := fmt.Sprintf("line %d: syntax error near '%s'.", lexer.line, nowToken) 
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
    lineNum, tokenType, token := lexer.NextToken()
    lexer.lineNum              = nowLineNum
    lexer.nextTokenType        = tokenType
    lexer.nextTokenLineNum     = lineNum
    lexer.nextToken            = token
    return tokenType
}

func (lexer *Lexer) nextDocumentIs(s string) bool {
    return strings.HasPrefix(lexer.document, s)
}

func (lexer *Lexer) skipDocument(n int) {
    lexer.document = lexer.document[n:]
}

func (lexer *Lexer) skipWhiteSpace() {
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
    // matching
    for len(lexer.document) > 0 {
        if lexer.nextDocumentIs("\r\n") || lexer.nextDocumentIs("\n\r") {
            lexer.skipDocument(2)
            lexer.lineNum += 1
        } else if isNewLine(lexer.document[0]) {
            lexer.skipDocument(1)
            lexer.lineNum += 1
        } else if matchWhiteSpace(lexer.document[0]) {
            lexer.skipDocument(1)
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
    painc("unreachable!")
}

func (lexer *Lexer) scanNumber() string {
    return lexer.scan(regexNumber)
}

func (lexer *Lexer) scanIdentifier() string {
    return lexer.scan(regexIdentifier)
}

// NextToken(), main lexer method
func (lexer *Lexer) NextToken() (lineNum int, tokenType int, token string) {
    // next token already loaded
    if lexer.nextTokenLineNum > 0 {
        lineNum                = lexer.nextTokenLineNum
        tokenType              = lexer.nextTokenType
        token                  = lexer.nextToken
        lexer.lineNum          = lexer.nextTokenLineNum
        lexer.nextTokenLineNum = 0
    }
    // skip spaces
    lexer.skipWhiteSpace()
    // finish
    if len(lexer.document) == 0 {
        return lexer.lineNum, TOKEN_EOF, "EOF"
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
    if lexer.document[0] == '.' || isDigit(lexer.document[0]) {
        token := lexer.scanNumber()
        return lexer.line, TOKEN_NUMBER, token
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