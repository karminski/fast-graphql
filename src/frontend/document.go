// document.go

package frontend

import (
    "errors"
)

type Document struct {
    LastLineNum       int
    Definitions       []Definition 
    // ReturnExpressions []Expression
}


func (document *Document) GetDefinitions() []Definition {
    return document.Definitions
}

func (document *Document) GetOperationDefinition() (*OperationDefinition, error) {
    for _, definition := range document.Definitions {
        if definition.GetDefinitionType() == OperationDefinitionType {
            return definition.(*OperationDefinition), nil
        }
    }
    return nil, errors.New("GetOperationDefinition(): input Document does not have OperationDefinition.")
}