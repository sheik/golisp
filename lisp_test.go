package main

import "testing"

func TestEval(t *testing.T) {
	cases := []struct {
		in   string
		want Object
	}{
		{"(define x 10)", Number(10)},
		{"(define r 10.3)", Number(10.3)},
		{"(+ 2 2)", Number(4)},
		{"(* 2 2)", Number(4)},
		{"(* 2 x)", Number(20)},
		{"(define x (+ 10 (* pi (* r r))))", Number(343.2915646628601)},
		{"(if (> x 2) x 0)", Number(343.2915646628601)},
		{"(if (> 2 x) x 0)", Number(0)},
		{"(* 2 (+ 3 (- 10 8)))", Number(10)},
		{"(+ 2 2 2 2 2 2)", Number(12)},
		{"(* 2 2 2)", Number(8)},
		{"(circle-area 3)", Number(28.274333886)},
		{"(fact 5)", Number(120)},
		{"(fact 3)", Number(6)},
		{"(fact 8)", Number(40320)},
		{"(begin (define r 10) (* pi (* r r)))", Number(314.1592654)},
		{"(fib 0)", Number(1)},
		{"(fib 1)", Number(1)},
		{"(fib 2)", Number(2)},
		{"(fib 3)", Number(3)},
		{"(fib 4)", Number(5)},
		{"(fib 5)", Number(8)},
		{"(fib 6)", Number(13)},
		{"(fib 7)", Number(21)},
		{"(fib 8)", Number(34)},
		{"(fib 9)", Number(55)},
		{"(fib 10)", Number(89)},
		{"(car (quote (5 4 3 2 1)))", Number(5)},
		{"(second (quote (1 2 3 4 5)))", Number(2)},
	}
	e := getStandardEnv()

	// define some functions
	e.eval(parse("(define circle-area (lambda (r) (* pi (* r r))))"))
	e.eval(parse("(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))"))
	e.eval(parse("(define fib (lambda (n) (if (<= n 1) 1 (+ (fib (- n 1)) (fib (- n 2)))))))"))
	e.eval(parse("(define second (lambda (x) (car (cdr x))))"))

	for _, c := range cases {
		got := e.eval(parse(c.in))
		if got != c.want {
			t.Errorf("eval(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
