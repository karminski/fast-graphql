// definition.go

package frontend



type Definition interface{}

type EmptyDefinition struct{}

type OperationTypeDefinition struct {
	LineNum 	  int 
	TokenName 	  string // anonymous operation if OperationName is empty.
	OperationType int
}


type FieldDefinition struct {
	LineNum 	int 
	TokenName 	string
}