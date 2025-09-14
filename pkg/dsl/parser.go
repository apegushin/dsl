package dsl

import "fmt"

type AST struct{}

func Parse(input string) {
	_, tokenChan := lex(input)
	for t := range tokenChan {
		fmt.Println(t)
	}

	fmt.Println()
}

//func Parse(tokens <-chan token) (AST, error) {
//	return AST{}, nil
//}
