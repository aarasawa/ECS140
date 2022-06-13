reachable(Nfa, StartState, FinalState, Input) :-
    [H | T] = Input, 
    transition(Nfa, StartState, H, Tstates),
    [A | _] = Tstates,
    reachable(Nfa, A, FinalState, T).

reachable(Nfa, StartState, FinalState, Input) :-
    [H | T] = Input, 
    transition(Nfa, StartState, H, Tstates),
    [_ | B] = Tstates, [A | _] = B,
    reachable(Nfa, A, FinalState, T).

reachable(Nfa, StartState, FinalState, []) :- 
    Nfa \= [], 
    FinalState = StartState.
