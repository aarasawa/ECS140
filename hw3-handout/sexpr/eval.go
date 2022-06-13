package sexpr

import (
	"errors"
	"fmt"
	"math/big"
	// You will need to use this package in your implementation.
)

// ErrEval is the error value returned by the Evaluator if the contains
// an invalid token.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrEval = errors.New("eval error")
var length_count = 0
var recurs_flag = false
var nil_flag = false
var num_flag = false

func (expr *SExpr) Eval() (*SExpr, error) {
	//var newex SExpr
	fmt.Println("before")
	fmt.Println(expr.SExprString())
	// fmt.Println(expr.atom == nil, expr.car.isNil(), expr.cdr.isNil(), expr.isNil())
	//fmt.Println(expr.car.atom.literal, expr.cdr.SExprString())
	switch {
	case recurs_flag:
		fmt.Println("recurs_flag")
		fmt.Println(expr.atom, expr.car, expr.cdr, expr.isNumber())
		if expr.isNumber() {
			recurs_flag = false
			return expr, nil
		}
		if expr.isNil() {
			return expr, nil
		}
		length_count++
		fmt.Println("pass isNum", length_count, expr.isNil())
		//fmt.Println(expr.car.SExprString(), expr.cdr.SExprString(), expr.car.isNumber())
		if expr.cdr.isNil() && expr.car.isNumber() {
			fmt.Println("recurs nil and car is num")
			recurs_flag = false
			//fmt.Println(expr.car)
			return expr.car, nil
		} else {
			expr.cdr.Eval()
		}
	case expr.isNil() || (expr.atom == nil && expr.car.isNil() && expr.cdr.isNil()):
		fmt.Println("nils")
		return expr, nil
	case (expr.atom != nil && expr.cdr == nil && expr.car == nil):
		//symbol
		if expr.atom.typ == tokenSymbol {
			fmt.Println("SYMBOL")
			return nil, ErrEval
		} else if expr.atom.typ == tokenNumber {
			fmt.Println("NUMBER")
			num_flag = true
			if nil_flag && num_flag {
				nil_flag = false
				num_flag = false
				return nil, ErrEval
			}
			return expr, nil
		}
	case expr.cdr.isNil():
		fmt.Println("CDR NIL", expr.SExprString())
		nil_flag = true
		expr, err := (expr.car).Eval()
		if err != nil {
			return expr, ErrEval
		}
		return expr, err
	case (expr.car.atom.literal == "CAR"):
		fmt.Println("CAR")
		if expr.cdr.car.atom.literal == "NIL" {
			expr.cdr.car.atom = nil
		}
		expr, err := (expr.cdr).Eval()
		fmt.Println(expr.SExprString())
		return expr.car, err
	case (expr.car.atom.literal == "CDR"):
		fmt.Println("CDR")
		fmt.Println(expr.car.SExprString(), expr.cdr.SExprString())
		if expr.cdr.car.atom.literal == "NIL" {
			expr.cdr.car.atom = nil
		}
		expr, err := (expr.cdr).Eval()
		fmt.Println(expr.SExprString())
		return expr.cdr, err
	case (expr.car.atom.literal == "LENGTH"):
		fmt.Println("LENGTH")
		expr, err := (expr.cdr).Eval()
		if err != nil {
			return expr, ErrEval
		}
		fmt.Println(expr.SExprString(), expr.isNil(), length_count)
		new_expr := mkNumber(big.NewInt(int64(length_count)))
		if expr.isNil() {
			return new_expr, nil
		} else {
			fmt.Println(expr.car.SExprString(), expr.cdr.SExprString())
			length_count++
			if expr.cdr.isNil() {
				return expr.car, err
			} else {
				recurs_flag = true
				expr1, err1 := (expr.cdr).Eval()
				if err1 != nil {
					return expr1, ErrEval
				}
				fmt.Println("get here?")
				fmt.Println(expr1)
				new_expr := mkNumber(big.NewInt(int64(length_count)))
				return new_expr, err1
			}
		}
	case (expr.car.atom.literal == "QUOTE") || (expr.atom.typ == tokenQuote) || (expr.atom.literal == "QUOTE"):
		fmt.Println("QUOTE.S")
		length_count = 0
		fmt.Println(expr.SExprString())
		return expr.cdr.car, nil
	}
	return nil, nil
}


