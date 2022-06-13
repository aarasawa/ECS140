/* All novels published either during the year 1953 or during the year 1996*/
year_1953_1996_novels(Book) :-
    %% remove fail and add body/other cases for this predicate
    novel(Book, X), member(X, [1953, 1996]).

/* List of all novels published during the period 1800 to 1900 (not inclusive)*/
period_1800_1900_novels(Book) :-
    %% remove fail and add body/other cases for this predicate
    novel(Book, X), X > 1800, X < 1900.

/* Characters who are fans of LOTR */
lotr_fans(Fan) :-
    %% remove fail and add body/other cases for this predicate
    fan(Fan, X), member(the_lord_of_the_rings, X).

/* Authors of the novels that heckles is fan of. */
heckles_idols(Author) :-
    %% remove fail and add body/other cases for this predicate
    fan(heckles, X), author(Author, Y),
    member(B, X), member(B, Y).

/* Characters who are fans of any of Robert Heinlein's novels */
heinlein_fans(Fan) :-
    %% remove fail and add body/other cases for this predicate
    author(robert_heinlein, X), fan(Fan, Y),
    member(C, X), member(C, Y).

/* Novels common between either of Phoebe, Ross, and Monica */
mutual_novels(Book) :-
    %% remove fail and add body/other cases for this predicate
    fan(phoebe, X), fan(ross, Y),
    member(Book, X), member(Book, Y).
mutual_novels(Book) :-
    %% other case for this predicate
    fan(phoebe, X), fan(monica, Z),
    member(Book, X), member(Book, Z).
mutual_novels(Book) :-
    %% other case for this predicate
    fan(ross, Y), fan(monica, Z),
    member(Book, Y), member(Book, Z).
