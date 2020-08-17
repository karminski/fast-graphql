// document.go

package frontend

type Document struct {
	LastLineNum 	  int
	Definitions       []Definition 
	ReturnExpressions []Expression
}