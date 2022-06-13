package unify

import (
	"errors"
	"hw4/disjointset"
	"hw4/term"
)

// ErrUnifier is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrUnifier = errors.New("unifier error")

// UnifyResult is the result of unification. For example, for a variable term
// `s`, `UnifyResult[s]` is the term which `s` is unified with.
type UnifyResult map[*term.Term]*term.Term

// Unifier is the interface for the term unifier.
// Do not change the definition of this interface
type Unifier interface {
	Unify(*term.Term, *term.Term) (UnifyResult, error)
}

// NewUnifier creates a struct of a type that satisfies the Unifier interface.
func NewUnifier() Unifier {
	return &NewParser{
		term_count:   0,
		set:          disjointset.NewDisjointSet(),
		size:         make(map[*term.Term]int),
		term_to_int:  make(map[*term.Term]int),
		int_to_term:  make(map[int]*term.Term),
		schema:       make(map[*term.Term]*term.Term),
		asyclic:      make(map[*term.Term]bool),
		visited:      make(map[*term.Term]bool),
		vars:         make(map[*term.Term][]*term.Term),
		final_answer: make(map[*term.Term]*term.Term),
	}
}

type NewParser struct {
	term_count int

	set disjointset.DisjointSet //contains parent map. Initialize it by pointing to itself

	size map[*term.Term]int

	term_to_int map[*term.Term]int

	int_to_term map[int]*term.Term

	schema map[*term.Term]*term.Term

	asyclic map[*term.Term]bool

	visited map[*term.Term]bool

	vars map[*term.Term][]*term.Term

	final_answer UnifyResult
	//need schema. schema: Term -> Term
	//need boolean flags asyclic and visited. Term -> bool
	//a pointer vars from each representative to a list of all variables for findSolution. Term -> Set of terms //initialized to empty for non variable nodes and Vars[X] = X for variable nodes
}

func (z *NewParser) Unify(s *term.Term, t *term.Term) (UnifyResult, error) {
	_, err := z.UnifClosure(s, t)
	if err != nil {
		return nil, ErrUnifier
	}
	return z.FindSolution(s)
}

//algorithm only distinguishes between variables and non variables
//we should distinguish between variables, numbers, atoms, and compound terms
//create a map between terms and int so we can use the disjointset
func (z *NewParser) UnifClosure(s *term.Term, t *term.Term) (UnifyResult, error) {
	//initialize for term s
	if _, exists := z.schema[s]; !exists {
		z.int_to_term[z.term_count] = s
		z.term_to_int[s] = z.term_count
		z.term_count++
		z.schema[s] = s
		z.size[s] = 1
		if s.Typ == term.TermVariable {
			z.vars[s] = append(z.vars[s], s)
		}
	}
	//initialize for term t
	if _, exists2 := z.schema[t]; !exists2 {
		z.int_to_term[z.term_count] = t
		z.term_to_int[t] = z.term_count
		z.term_count++
		z.schema[t] = t
		z.size[t] = 1
		if t.Typ == term.TermVariable {
			z.vars[t] = append(z.vars[t], t)
		}
	}
	number1 := z.set.FindSet(z.term_to_int[s])
	number2 := z.set.FindSet(z.term_to_int[t])
	term1, term2 := z.int_to_term[number1], z.int_to_term[number2]
	if term1 == term2 {
		return z.schema, nil
	} else {
		f, g := z.schema[term1], z.schema[term2]
		if f.Typ != term.TermVariable && g.Typ != term.TermVariable {
			if f.Typ == term.TermCompound && g.Typ == term.TermCompound && f.Functor == g.Functor && len(f.Args) == len(g.Args) {
				z.Union(term1, term2)
				for i := 0; i < len(f.Args); i++ {
					_, err := z.UnifClosure(f.Args[i], g.Args[i])
					if err != nil {
						return nil, ErrUnifier
					}
				}
			} else {
				return nil, ErrUnifier
			}
		} else {
			z.Union(term1, term2)
		}
	}
	return z.schema, nil
}

func (z *NewParser) Union(s *term.Term, t *term.Term) {
	int_term1, int_term2 := z.term_to_int[s], z.term_to_int[t]
	if z.size[s] >= z.size[t] {
		z.size[s] += z.size[t]
		z.vars[s] = append(z.vars[s], z.vars[t]...)
		if z.schema[s].Typ == term.TermVariable {
			z.schema[s] = z.schema[t]
		}
		z.set.UnionSet(int_term1, int_term2)
	} else {
		z.size[t] += z.size[s]
		z.vars[t] = append(z.vars[t], z.vars[s]...)
		if z.schema[t].Typ == term.TermVariable {
			z.schema[t] = z.schema[s]
		}
		z.set.UnionSet(int_term1, int_term2)
	}
}

func (z *NewParser) FindSolution(s *term.Term) (UnifyResult, error) {
	//initialize for term s if doesn't exist
	if _, exists := z.schema[s]; !exists {
		z.int_to_term[z.term_count] = s
		z.term_to_int[s] = z.term_count
		z.term_count++
		z.schema[s] = s
		z.size[s] = 1
		if s.Typ == term.TermVariable {
			z.vars[s] = append(z.vars[s], s)
		}
	}
	term1 := z.schema[z.int_to_term[z.set.FindSet(z.term_to_int[s])]]
	if z.asyclic[term1] {
		return z.schema, nil
	}
	if z.visited[term1] {
		return nil, ErrUnifier
	}
	if term1.Typ == term.TermCompound {
		z.visited[term1] = true
		for i := 0; i < len(term1.Args); i++ {
			_, err := z.FindSolution(term1.Args[i])
			if err != nil {
				return nil, ErrUnifier
			}
		}
		z.visited[term1] = false
	}
	z.asyclic[term1] = true
	newterm := z.int_to_term[z.set.FindSet(z.term_to_int[s])]
	//newterm := z.int_to_term[z.set.FindSet(z.term_to_int[term1])]
	for j := 0; j < len(z.vars[newterm]); j++ {
		if z.vars[newterm][j] != term1 {
			z.final_answer[z.vars[newterm][j]] = term1
		}
	}
	return z.final_answer, nil
}
