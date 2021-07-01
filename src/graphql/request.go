// request.go
// 
package graphql

import (
	"crypto/md5"

)

type Request struct {
    // GraphQL Query string from client side
    Query string 

    // md5 hash for Query
    QueryHash [16]byte

    // 
    SubstitutedQuery string

    // md5 hash for SubstitutedQuery
    SubstitutedQueryHash [16]byte

    // GraphQL Query variables from client side
    Variables map[string]interface{}

    // substituted variables from arguments scanner
    SubstitutedVariables map[string]interface{}

    // operation name from client side
    OperationName string
}

func (r *Request)GenerateQueryHash() {
	r.QueryHash	= md5.Sum([]byte(r.Query))
}

func (r *Request)GenerateSubstitutedQueryHash() {
	r.SubstitutedQueryHash	= md5.Sum([]byte(r.SubstitutedQuery))
}	

func (r *Request)InitVariables() {
	r.Variables = make(map[string]interface{})
}

func (r *Request)InitSubstitutedVariables() {
	r.SubstitutedVariables = make(map[string]interface{})
}

func (r *Request)GetQueryVariables() (map[string]interface{}) {
	if len(r.SubstitutedVariables) != 0 {
		return r.SubstitutedVariables
	}
	return r.Variables
}

func (r *Request)GetQueryHash() ([16]byte) {
    var emptyHash [16]byte  
    if r.SubstitutedQueryHash != emptyHash {
        return r.SubstitutedQueryHash
    }
    return r.QueryHash
}


func (r *Request)IsVariablesAvaliable() bool {
    if len(r.Variables) == 0 {
        return false
    }
    return true
}