% A list is a 1-D array of numbers.
% A matrix is a 2-D array of numbers, stored in row-major order.

% You may define helper functions here.
myappend([],L,L).
myappend([X|L1],L2,[X|L3]) :- myappend(L1,L2,L3).
% are_adjacent(List, A, B) returns true iff A and B are neighbors in List.
are_adjacent([A, B|_], A, B) :- !.
are_adjacent([B, A|_], A, B) :- !.
are_adjacent([_|List], A, B) :-
    are_adjacent(List, A, B).

% matrix_transpose(Matrix, Answer) returns true iff Answer is the transpose of
% the 2D matrix Matrix.
matrix_transpose([], []).
matrix_transpose([First|Rest], End) :- transpose_helper(First, [First|Rest], End).

transpose_helper([], _, []).
transpose_helper([_|Rest], Matrix, Final) :- first_column(Matrix, X), 
                                             not_first_column(Matrix, MMatrix),
                                             transpose_helper(Rest, MMatrix, Y), 
                                             myappend([X], Y, Final).

first_column([],[]).
first_column([[First|_]|Others], [First|Merge]):- first_column(Others,Merge).

not_first_column([],[]).
not_first_column([[_|Not_First]|Others], [Not_First|Merge]):- not_first_column(Others,Merge).

% are_neighbors(Matrix, A, B) returns true iff A and B are neighbors in the 2D
% matrix Matrix.
%are_neighbors([], A, A).
are_neighbors([FirstRow|Matrix], A, B) :- are_adjacent_helper([FirstRow|Matrix], A, B); 
                                          are_transpose_helper([FirstRow|Matrix], A, B).

are_adjacent_helper([FirstRow|Matrix], A, B) :- are_adjacent(FirstRow, A, B).
    											%are_adjacent_helper(Matrix, A, B).

are_transpose_helper([FirstRow|Matrix], A, B) :- matrix_transpose([FirstRow|Matrix], X),
    											 are_adjacent_helper(X, A, B).
