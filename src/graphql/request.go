// request.go
// 
package graphql

import (
	"crypto/md5"

)

type Request struct {
    // GraphQL Query string from client side
    Query string 


    // md5 hash for Request.Query
    QueryHash [16]byte

    SubstituteQuery string

    SubstituteQueryHash [16]byte


    // GraphQL Query variables from client side
    Variables map[string]interface{}

    // substituted variables from arguments scanner
    SubstitutedVariables map[string]interface{}
}

func (r *Request)GenerateQueryHash() {
	r.QueryHash	= md5.Sum([]byte(r.Query))
}

func (r *Request)GenerateSubstituteQueryHash() {
	r.SubstituteQueryHash	= md5.Sum([]byte(r.SubstituteQuery))
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