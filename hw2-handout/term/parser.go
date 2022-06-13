package term

import (
	"errors"
	//"fmt"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <start>    ::= <term> | \epsilon
// <term>     ::= ATOM | NUM | VAR | <compound>
// <compound> ::= <functor> LPAR <args> RPAR
// <functor>  ::= ATOM
// <args>     ::= <term> | <term> COMMA <args>
//

// Parser is the interface for the term parser.
// Do not change the definition of this interface.
type Parser interface {
	Parse(string) (*Term, error)
}

// TODO: implement a type that satisfies the Parser interface.
type MyParser struct {
	lex    *lexer
	repeat map[string]*Term
}

func ParserHelper2(f MyParser) (*Term, error, *Token) {
	var term *Term
	token, err := f.lex.next()
	for token.typ != tokenEOF && token.typ != tokenComma && token.typ != tokenRpar { //loops until eof, comma, and rpar is reached
		switch {
		case term == nil && token.typ == tokenAtom: //checks if atom
			term = &Term{Typ: TermAtom, Literal: token.literal}
			// check for existing repeat in map
			// else, add to map
			if isItRepeat, ok := f.repeat[term.String()]; ok {
				term = isItRepeat
			} else {
				f.repeat[term.String()] = term
			}
		case term == nil && token.typ == tokenNumber: //checks if number
			term = &Term{Typ: TermNumber, Literal: token.literal}
			// check for existing repeat in map
			// else, add to map
			if isItRepeat, ok := f.repeat[term.String()]; ok {
				term = isItRepeat
			} else {
				f.repeat[term.String()] = term
			}
		case term == nil && token.typ == tokenVariable: //checks if variable
			term = &Term{Typ: TermVariable, Literal: token.literal}
			// check for existing repeat in map
			// else, add to map
			if isItRepeat, ok := f.repeat[term.String()]; ok {
				term = isItRepeat
			} else {
				f.repeat[term.String()] = term
			}
		case term != nil && term.Typ == TermAtom: //checks if compound
			// expecting left parenthesis
			if token.typ == tokenLpar {
				term = &Term{Typ: TermCompound, Literal: "", Functor: term}
				term2, _, token2 := ParserHelper2(f) //first recursion for first argument
				if term2 != nil {
					term.Args = append(term.Args, term2)
				}
				for token2.typ == tokenComma { //will loop for the rest of the argument
					term2, _, token2 = ParserHelper2(f)
					if term2 != nil {
						term.Args = append(term.Args, term2)
					}
				}
				if token2.typ != tokenRpar || term2 == nil { //token has to be a right parenthesis or else the statement is invalid
					return nil, ErrParser, token2
				}
				// check for existing repeat in map
				// else, add to map
				if isItRepeat, ok := f.repeat[term.String()]; ok {
					term = isItRepeat
				} else {
					f.repeat[term.String()] = term
				}
			} else {
				// unexpected token.typ
				return nil, ErrParser, token
			}
		default:
			return nil, ErrParser, token
		}
		token, err = f.lex.next() //next loop
	}
	if term == nil && token.typ != tokenEOF {
		return nil, ErrParser, token
	}
	return term, err, token
}

func (f MyParser) Parse(str string) (*Term, error) {
	// list of terms
	f.lex = newLexer(str)
	FinalTerm, err, last := ParserHelper2(f)
	if FinalTerm != nil && (FinalTerm.Typ == TermVariable || FinalTerm.Typ == TermNumber || FinalTerm.Typ == TermAtom) && last.typ != tokenEOF {
		return nil, ErrParser
	}
	if FinalTerm != nil && FinalTerm.Typ == TermCompound && (last.typ != tokenRpar && last.typ != tokenEOF) {
		return nil, ErrParser
	}
	return FinalTerm, err
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	Lexer := MyParser{}
	Lexer.repeat = make(map[string]*Term)
	var p Parser = Lexer
	return p
}
