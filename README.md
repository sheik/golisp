# golisp
LISP interpreter written in Go

It currently supports a subset of LISP, but important features already work:

	golisp> (define r 10)
	10
	golisp> (* pi (* r r))
	314.1592654
	golisp> (define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))
	0x403320
	golisp> (fact 3)
	6
	golisp> (fact 5)
	120
	golisp> (fact 120)
	6.689502913449124e+198
	golisp> (define circle-area (lambda (r) (* pi (* r r))))
	0x403320
	golisp> (circle-area (fact 5))
	45238.9342176
	golisp> 

