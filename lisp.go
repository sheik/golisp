package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION = "golisp v1.0.2.1"
	VERBOSE = true
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

func recover_error() {
	if r := recover(); r != nil {
		fmt.Println("Parse Error:", r)
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		fmt.Printf("%s", buf)
	}
}

func parse(program string) Object {
	defer recover_error()

	tokens := tokenize(program)
	return build_ast(&tokens)
}

type Env struct {
	mapping map[Symbol]Object
	outer   *Env
}

type Object interface{}

type List []Object

func (n Number) String() string {
	if float64(n) == float64(int64(n)) {
		return fmt.Sprintf("%d", int64(n))
	}
	return fmt.Sprintf("%f", n)
}

func (l List) String() string {
	var s []string
	for _, v := range l {
		s = append(s, fmt.Sprintf("%s", v))
	}
	return "(" + strings.Join(s, " ") + ")"
}

type Symbol string

type Number float64

type Lambda struct {
	env   Env
	parms Object
	body  Object
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

func car(args []Object) Object {
	return args[0].(List)[0]
}

func cdr(args []Object) Object {
	return args[0].(List)[1:]
}

func print(args []Object) Object {
	return args[0]
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
	e.mapping["car"] = car
	e.mapping["cdr"] = cdr
	e.mapping["print"] = print
	e.mapping["pi"] = Number(3.141592654)
	return e
}

func (e *Env) eval(x Object) Object {
	defer recover_error()

	if val, is_symbol := x.(Symbol); is_symbol {
		return e.mapping[val]
	} else if _, is_list := x.(List); !is_list {
		return x
	}

	l := x.(List)

	if l[0] == Symbol("quote") {
		exp := l[1]
		return exp
	} else if l[0] == Symbol("define") {
		val := e.eval(l[2])
		e.mapping[l[1].(Symbol)] = val
		return val
	} else if l[0] == Symbol("if") {
		truth := e.eval(l[1]).(bool)
		if truth {
			return e.eval(l[2])
		} else {
			return e.eval(l[3])
		}
	} else if l[0] == Symbol("lambda") {
		parms, body := l[1], l[2]
		newenv := Env{}
		newenv.mapping = make(map[Symbol]Object)
		for k, v := range e.mapping {
			newenv.mapping[k] = v
		}
		return Lambda{newenv, parms, body}
	} else {
		proc := e.eval(l[0])

		if ln, is_lambda := proc.(Lambda); is_lambda {
			env := Env{}
			env.mapping = make(map[Symbol]Object)
			for k, v := range e.mapping {
				env.mapping[k] = v
			}

			for i, v := range ln.parms.(List) {
				val := env.eval(l[i+1])
				env.mapping[v.(Symbol)] = val
			}

			return env.eval(ln.body)
		}

		var args []Object
		for _, v := range l[1:] {
			args = append(args, e.eval(v))
		}
		res := proc.(func([]Object) Object)(args)
		return res
	}
}

func repl(e Env, profile_code bool) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("golisp> ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}
		if profile_code {
			start := time.Now()
			val := e.eval(parse(text))
			elapsed := time.Since(start)
			fmt.Println(val)
			fmt.Println("Execution took", elapsed)
		} else {
			fmt.Println(e.eval(parse(text)))
		}
	}
}

func main() {
	profile_code := flag.Bool("profile", false, "profile code execution time")
	show_version := flag.Bool("version", false, "show program version and exit")
	flag.Parse()

	if *show_version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	e := getStandardEnv()
	repl(e, *profile_code)
}
