package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func parse(program string) Object {
	tokens := tokenize(program)
	return build_ast(&tokens)
}

type Env struct {
	mapping map[Symbol]Object
	outer   *Env
}

type Procedure struct {
	body Object
	args Object
	env  *Env
}

type Object interface{}

type List []Object

type Symbol string

type Number float64

func (p *Procedure) call() Object {
	return p.env.eval(p.body) 
}

func build_ast(tokens *[]string) Object {
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

func atom(token string) Object {
	n, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return Number(n)
	}
	return Symbol(token)
}

func tokenize(chars string) []string {
	chars = strings.Replace(chars, "(", " ( ", -1)
	chars = strings.Replace(chars, ")", " ) ", -1)

	return removeEmpty(strings.Split(chars, " "))
}

func mult(args []Object) Object {
	result := args[0].(Number)
	for _, num := range args[1:] {
		result *= num.(Number)
	}
	return result
}

func add(args []Object) Object {
	result := args[0].(Number)
	for _, num := range args[1:] {
		result += num.(Number)
	}
	return result
}

func sub(args []Object) Object {
	result := args[0].(Number)
	for _, num := range args[1:] {
		result -= num.(Number)
	}
	return result
}

func div(args []Object) Object {
	result := args[0].(Number)
	for _, num := range args[1:] {
		result /= num.(Number)
	}
	return result
}

func gt(args []Object) Object {
	x, y := args[0].(Number), args[1].(Number)
	return x > y
}

func lt(args []Object) Object {
	x, y := args[0].(Number), args[1].(Number)
	return x < y
}

func lte(args []Object) Object {
	x, y := args[0].(Number), args[1].(Number)
	return x <= y
}

func gte(args []Object) Object {
	x, y := args[0].(Number), args[1].(Number)
	return x >= y
}

func begin(args []Object) Object {
	return args[len(args)-1]
}

func getStandardEnv() Env {
	e := Env{
		mapping: make(map[Symbol]Object),
	}
	e.mapping["begin"] = begin
	e.mapping["*"] = mult
	e.mapping["/"] = div 
	e.mapping["+"] = add
	e.mapping["-"] = sub
	e.mapping[">"] = gt
	e.mapping[">="] = gte
	e.mapping["<"] = lt
	e.mapping["<="] = lte
	e.mapping["pi"] = Number(3.141592654)
	return e
}

func (e *Env) eval(x Object) Object {
	if val, is_symbol := x.(Symbol); is_symbol {
		return e.mapping[val]
	} else if _, is_list := x.(List); !is_list {
		return x
	} else if l := x.(List); l[0] == Symbol("define") {
		val := e.eval(l[2])
		e.mapping[l[1].(Symbol)] = val
		return val
	} else if l := x.(List); l[0] == Symbol("if") {
		truth := e.eval(l[1]).(bool)
		if truth {
			return e.eval(l[2])
		} else {
			return e.eval(l[3])
		}
	} else if l := x.(List); l[0] == Symbol("lambda") {
		parms, body := l[1], l[2]
		return func(args []Object) Object {
				for i, v := range parms.(List) {
					e.mapping[v.(Symbol)] = args[i]
				}
				return e.eval(body)
			}
	} else {
		proc := e.eval(l[0])
		var args []Object
		for _, v := range l[1:] {
			args = append(args, e.eval(v))
		}
		res := proc.(func([]Object) Object)(args)
		return res
	}
}

func repl(e Env) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("golisp> ")
		text, _ := reader.ReadString('\n')
		fmt.Println(e.eval(parse(text)))
	}
}

func main() {
	e := getStandardEnv()

	repl(e)
}
