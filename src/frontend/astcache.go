package frontend

import (
    "fast-graphql/src/graphql"
	
	"sync"

)



var astCache sync.Map // map[queryHash]Document

func loadAST(request *graphql.Request) (Document, bool)  {
	if doc, ok := astCache.Load(request.QueryHash); ok {
		return doc.(Document), true
	}
	var doc Document 
	return doc, false
}


func saveAST(request *graphql.Request, doc Document) {
	astCache.Store(request.QueryHash, doc)
}


func variablesSubstitution() {
	return
}