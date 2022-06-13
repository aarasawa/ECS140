package sexpr

import (
	"errors"
	"fmt"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <sexpr>       ::= <atom> | <pars> | QUOTE <sexpr>
// <atom>        ::= NUMBER | SYMBOL
// <pars>        ::= LPAR <dotted_list> RPAR | LPAR <proper_list> RPAR
// <dotted_list> ::= <proper_list> <sexpr> DOT <sexpr>
// <proper_list> ::= <sexpr> <proper_list> | \epsilon
//
type Parser interface {
	Parse(string) (*SExpr, error)
}

// TODO: implement a type that satisfies the Parser interface.
type MyParser struct {
	lex       *lexer
	peekTok   *token
	parLayers int
}

// gets the next token either by readig PeekTok or from the lexer
func (p *MyParser) nextToken() (*token, error) {
	if tok := p.peekTok; tok != nil {
		p.peekTok = nil
		return tok, nil
	}
	return p.lex.next()
}

// puts back tok
func (p *MyParser) backToken(tok *token) {
	p.peekTok = tok
}

// main parse function
func (p *MyParser) Parse(input string) (*SExpr, error) {
	// create lexer from input string
	p.lex = newLexer(input)
	p.peekTok = nil

	// empty string case
	tok, err := p.nextToken()
	if err != nil {
		return nil, ErrParser
	}
	if tok.typ == tokenEOF {
		return nil, ErrParser
	}

	// entry to recursive parsing
	p.backToken(tok)
	sexpr, err := p.parseNextSExpr()

	if err != nil {
		return nil, ErrParser
	}
	// error: mismatched parenthesis
	if p.parLayers != 0 {
		return nil, ErrParser
	}
	// error: didn't consume entire input
	if tok, err := p.nextToken(); err != nil || tok.typ != tokenEOF {
		return nil, ErrParser
	}
	return sexpr, nil
}

// parses a prefix of the string (from lexer) into a SExpr
// or returns an error
func (p *MyParser) parseNextSExpr() (*SExpr, error) {
	tok, err := p.nextToken()
	fmt.Println("TOKEN")
	fmt.Println(tok)
	//fmt.Println(tok.typ)
	if err != nil {
		return nil, err
	}
	switch tok.typ {
	case tokenEOF:
		return nil, nil
	// simple cases
	case tokenNumber:
		return mkNumber(tok.num), nil
	case tokenSymbol:
		// "quote" case falls into here
		if tok.literal == "QUOTE" {
			arg, err := p.parseNextSExpr()
			if err != nil {
				return nil, err
			}
			// DOT case
			if arg.atom.typ == tokenDot {
				fmt.Println("DOT_CASE")
				arg, err = p.parseNextSExpr()
				if err != nil {
					return nil, err
				}
				ret := mkConsCell(mkSymbol("QUOTE"), arg)
				fmt.Println(ret.SExprString())
				fmt.Println(ret.atom)
				ret.atom = mkTokenSymbol("QUOTE")
				return ret, nil
			}
		}
		// regular symbol case
		return mkSymbol(tok.literal), nil
	case tokenQuote:
		fmt.Println("NOT_DOT_QUOTE")
		arg, err := p.parseNextSExpr()
		if err != nil {
			return nil, err
		}
		// "'" invalid case
		if arg == nil {
			fmt.Println("INVALID '")
			return nil, ErrParser
		}
		// not DOT case
		ret := mkConsCell(mkSymbol("QUOTE"), mkConsCell(arg, mkNil()))
		fmt.Println(ret.SExprString())
		fmt.Println(ret.atom)
		ret.atom = mkTokenQuote()
		return ret, nil
	// might be the start of a list
	case tokenLpar:
		fmt.Println("LPAR")
		p.parLayers++
		// parse first element
		arg, err := p.parseNextSExpr()
		if err != nil {
			return nil, err
		}
		fmt.Println("FIRST ELE")
		// "(" invalid case
		if arg == nil {
			fmt.Println("INVALID (")
			return nil, ErrParser
		}
		fmt.Println(arg.SExprString())
		//fmt.Println(arg.atom)
		//fmt.Println(arg.atom.typ)
		if arg.atom != nil {
			// "()" case
			if arg.atom.typ == tokenRpar {
				return mkNil(), nil
			}
			// could be (QUOTE . X) case
			if arg.atom.literal == "QUOTE" && arg.atom.typ == tokenSymbol {
				fmt.Println("LPAR_QUOTE")
				// confirm quote closes
				next, err := p.parseNextSExpr()
				if err != nil {
					return mkNil(), err
				}
				if next.atom.typ != tokenRpar {
					fmt.Println("NOT_RPAR")
					return nil, ErrParser
				}
				fmt.Println("LPAR_QUOTE_RPAR")
				return arg, nil
			}
		}
		// array to hold elements of list
		args := []*SExpr{arg}
		// flag for if a dot was found
		foundDot := false
		// parse the elements if any
		for {
			fmt.Println("LOOP")
			arg, err := p.parseNextSExpr()
			if err != nil {
				return nil, ErrParser
			}
			// eof
			if arg == nil {
				fmt.Println("LOOP_EOF")
				return nil, ErrParser
			}
			if arg.atom != nil {
				if arg.atom.typ == tokenRpar {
					fmt.Println("LOOP_RPAR")
					break
				}
				if arg.atom.typ == tokenDot {
					foundDot = true
					continue
				}
			}
			fmt.Println(arg.SExprString())
			// append to args
			args = append(args, arg)
			fmt.Println("ARGS")
			for _, v := range args {
				fmt.Println(v.SExprString())
			}
			// only one SExpr should be after Dot
			if foundDot {
				// after that should be a ")"
				arg, err = p.parseNextSExpr()
				if arg.atom.typ != tokenRpar {
					fmt.Println("NOT_RPAR")
					return nil, ErrParser
				}
				fmt.Println("RPAR")
				break
			}
		}
		// build and return list
		ret := mkNil()
		if foundDot {
			ret = args[len(args)-1]
		}
		for i := len(args) - 1; i >= 0; i-- {
			// skip first loop if dot
			if foundDot {
				foundDot = false
				continue
			}
			ret = mkConsCell(args[i], ret)
			fmt.Println("CONS")
			fmt.Println(ret.SExprString())
		}
		fmt.Println("FINISHED_LIST")
		fmt.Println(ret.SExprString())
		return ret, nil
	case tokenDot:
		fmt.Println("DOT")
		return mkAtom(mkTokenDot()), nil
	case tokenRpar:
		fmt.Println("RPAR")
		p.parLayers--
		return mkAtom(mkTokenRpar()), nil
	default:
		fmt.Println("DEFAULT")
		return nil, ErrParser
	}
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	return &MyParser{
		lex:     nil,
		peekTok: nil,
	}
}
