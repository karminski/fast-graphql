package frontend

import (
	"sync"
)



var astCache sync.Map // map[queryHash]Document


func loadAST(queryHash [16]byte) (Document, bool)  {
	if f, ok := astCache.Load(queryHash); ok {
		return f.(Document), true
	}
	var doc Document 
	return doc, false
}


func saveAST(queryHash [16]byte, doc Document) {
	astCache.Store(queryHash, doc)
}


func variablesSubstitution() {
	return
}