package frontend

import (
	"sync"
	"crypto/md5"
	
)



var astCache sync.Map // map[queryHash]Document

func getQueryHash(query string) [16]byte {
	return md5.Sum([]byte(query))
}

func loadAST(queryHash [16]byte) (Document, bool)  {
	if doc, ok := astCache.Load(queryHash); ok {
		return doc.(Document), true
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