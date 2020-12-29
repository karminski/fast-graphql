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
    var operationDefinition *OperationDefinition
    var hit int8
    hit = 0
    // pickup OperationDefinition
    for _, definition := range document.Definitions {
        if definition.GetDefinitionType() == OperationDefinitionType {
            operationDefinition = definition.(*OperationDefinition)
            hit ++
        }
    }
    // check
    if hit == 0 {
        return nil, errors.New("GetOperationDefinition(): input Document does not have OperationDefinition.")
    }
    if hit > 1 {
        return nil, errors.New("GetOperationDefinition(): multiple OperationDefinition detected, please check your GraphQL syntax.")

    }
    return operationDefinition, nil
}

