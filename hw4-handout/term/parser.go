package term

import (
	"errors"
 	"strconv"
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

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	return &ParserImpl{
		lex:		nil,
		peekTok:	nil,
		terms:		make(map[string]*Term),
		termID:		make(map[*Term]int),
		termCounter:0,	
	}
}

type ParserImpl struct{
	// Lexer, initialized at each call to Parser
	lex *lexer
	//Look ahead token, initialized at each call to Parse
	peekTok *Token
	//Map from string representing a term to a term
	terms map[string]*Term
	//Map from Term to its ID.
	termID map[*Term]int
	//counter
	termCounter int
}

//nextToken gets the next token either by reading peekTok or
//from the lexer.
func (p *ParserImpl) nextToken() (*Token, error) {
	if tok := p.peekTok; tok != nil {
		p.peekTok = nil
		return tok, nil
	}
	return p.lex.next()
}

//backToken puts back tok.
func (p *ParserImpl) backToken(tok *Token) {
	p.peekTok = tok
}

//Parse a term
func (p *ParserImpl) Parse(input string) (*Term, error) {
	p.lex = newLexer(input)
	p.peekTok = nil

	// If the input is an empty string
	tok, err := p.nextToken()
	if err != nil {
		return nil, ErrParser
	}
	if tok.typ == tokenEOF {
		return nil, nil
	}
	p.backToken(tok)
	term, err := p.parseNextTerm()
	//term, err := p.termNT() //Table-driven parser
	if err != nil {
		return nil, ErrParser
	}
	//Error if we have not consumed all of the input.
	if tok, err := p.nextToken(); err != nil || tok.typ != tokenEOF {
		return nil, ErrParser
	}
	return term, nil
}

//parseNextTerm parses a prefix of the string (via the lexer) into a Term, or
//returns an error
func (p *ParserImpl) parseNextTerm() (*Term, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	switch tok.typ {
	case tokenEOF:
		return nil, nil
	case tokenNumber:
		return p.mkSimpleTerm(TermNumber, tok.literal), nil
	case tokenVariable:
		return p.mkSimpleTerm(TermVariable, tok.literal), nil
	case tokenAtom:
		a := p.mkSimpleTerm(TermAtom, tok.literal)
		nxt, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		if nxt.typ != tokenLpar {
			//Atom is not the functor for a compound term.
			p.backToken(nxt)
			return a, nil
		}
		//Atom might be the functor of a compound term.
		arg, err := p.parseNextTerm()
		if err != nil {
			return nil, err
		}
		//Args of a compound term contains at least one Term.
		args := []*Term{arg}
		nxt, err = p.nextToken()
		if err != nil {
			return nil, err
		}
		//Parse the rest of the arguments, if any.
		for ; nxt.typ == tokenComma; nxt, err = p.nextToken() {
			arg, err = p.parseNextTerm()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		if nxt.typ != tokenRpar {
			return nil, ErrParser
		}
		return p.mkCompoundTerm(a, args), nil
	default:
		return nil, ErrParser
	}
}

//Helper functions to make terms

//mkSimpleTerm makes a simple term.
func (p *ParserImpl) mkSimpleTerm(typ TermType, lit string) *Term {
	key := lit //Use the literal as the key for the simple terms.
	term, ok := p.terms[key]
	if !ok {
		term = &Term{Typ: typ, Literal: lit}
		p.insertTerm(term, key)
	}
	return term
}

//mkCompoundTerm makes a compound term
func (p *ParserImpl) mkCompoundTerm(functor *Term, args []*Term) *Term {
	key := strconv.Itoa(p.termID[functor])
	for _, arg := range args {
		key += ", " + strconv.Itoa(p.termID[arg])
	}
	term, ok := p.terms[key]
	if !ok {
		term = &Term{
			Typ:		TermCompound,
			Functor:	functor,
			Args:		args,
		}
		p.insertTerm(term, key)
	}
	return term
}

//insertTerm inserts term with given key into the terms and termID maps.
func (p *ParserImpl) insertTerm(term *Term, key string) {
	p.terms[key] = term
	p.termID[term] = p.termCounter
	p.termCounter++
}