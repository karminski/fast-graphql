package frontend

import (
    "fast-graphql/src/graphql"
	
	"sync"

)



var astCache sync.Map // map[queryHash]Document

func loadAST(request *graphql.Request) (Document, bool)  {
	hash := request.GetQueryHash()
	if doc, ok := astCache.Load(hash); ok {
		return doc.(Document), true
	}
	var doc Document 
	return doc, false
}


func saveAST(request *graphql.Request, doc Document) {
	hash := request.GetQueryHash()
	astCache.Store(hash, doc)
}


func variablesSubstitution() {
	return
}