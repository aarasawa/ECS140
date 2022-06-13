package nfa

import (
	//"fmt"
	"sync"
)

// A nondeterministic Finite Automaton (NFA) consists of states,
// symbols in an alphabet, and a transition function.

// A state in the NFA is represented as an unsigned integer.
type state uint

var mut sync.Mutex

//var wg sync.WaitGroup
//var reaches = make(chan bool)
//var empty_array []rune
var maps map[bool]bool

// Given the current state and a symbol, the transition function
// of an NFA returns the set of next states the NFA can transition to
// on reading the given symbol.
// This set of next states could be empty.
type TransitionFunction func(st state, sym rune) []state

// Reachable returns true if there exists a sequence of transitions
// from `transitions` such that if the NFA starts at the start state
// `start` it would reach the final state `final` after reading the
// entire sequence of symbols `input`; Reachable returns false otherwise.
func Reachable(
	// `transitions` tells us what our NFA looks like
	transitions TransitionFunction,
	// `start` and `final` tell us where to start, and where we want to end up
	start, final state,
	// `input` is a (possible empty) list of symbols to apply.
	input []rune,
) bool {
	maps = nil
	return ReachableHelper(transitions, start, final, input)
}

func ReachableHelper(
	// `transitions` tells us what our NFA looks like
	transitions TransitionFunction,
	// `start` and `final` tell us where to start, and where we want to end up
	start, final state,
	// `input` is a (possible empty) list of symbols to apply.
	input []rune,
) bool {
	if maps == nil {
		maps = make(map[bool]bool)
	}
	c := make(chan bool)
	if len(input) == 0 {
		close(c)
		if start == final {
			//fmt.Println("true", start, final)
			return true
		} else {
			return false
		}
	}
	var answer bool
	for _, next := range transitions(start, input[0]) {
		go func() {
			solution := Reachable(transitions, next, final, input[1:])
			c <- solution
		}()
		answer = <-c
		mut.Lock()
		maps[answer] = answer
		mut.Unlock()
		if _, ok := maps[true]; ok {
			return true
		}
	}
	close(c)
	return false
}
