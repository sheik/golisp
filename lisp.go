package main

import (
	"fmt"
//	"math"
	"strings"
)

func removeEmpty(tokens []string) []string {
	var result []string
	for _, c := range tokens {
		if c != "" {
			result = append(result, c)
		}
	}
	return result
}

func parse(program string) Atom {
	tokens := tokenize(program)
	return build_ast(&tokens)
}

type Atom interface{}

type List []Atom

func build_ast(tokens *[]string) Atom {
	token := pop(tokens)

	switch token {
	case "(":
		var L List
		for (*tokens)[0] != ")" {
			L = append(L, build_ast(tokens))
		}
		pop(tokens)
		return L
	case ")":
		panic("unexpected )")
	default:
		return atom(token)
	}
}

func pop(tokens *[]string) string {
	if len(*tokens) == 0 {
		panic("unexpected EOF while reading")
	}
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]

	return token
}

func atom(token string) Atom {
	return token 
}

func tokenize(chars string) []string {
	chars = strings.Replace(chars, "(", " ( ", -1)
	chars = strings.Replace(chars, ")", " ) ", -1)

	return removeEmpty(strings.Split(chars, " "))
}

type env map[string]func(float64)float64

func getStandardEnv() env {
	var e env
	
	return e
}

func main() {
	program := "(begin (define r 10) (* pi (* r r)))"
	fmt.Printf("%q\n", parse(program))
//	e := getStandardEnv()
//	result := e["abs"](-2.2)
//	fmt.Println(result)
}

