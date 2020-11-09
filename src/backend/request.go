// request.go
package backend

import (

)


type Fields map[string]*Field

type Field struct {
	Type 	    string
	ResolveFunc func(interface{}) (interface{}, error)
}