// request.go
// 
package graphql

import (

)

type Request struct {
    // GraphQL Query string from client side
    Query string 

    // md5 hash for Request.Query
    QueryHash [16]byte

    // GraphQL Query variables from client side
    Variables map[string]interface{}
}

