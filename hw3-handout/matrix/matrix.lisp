; A list is a 1-D array of numbers.
; A matrix is a 2-D array of numbers, stored in row-major order.

; If needed, you may define helper functions here.

; AreAdjacent returns true iff a and b are adjacent in lst.
(defun are-adjacent (lst a b)
    (cond((eql lst NIL) NIL)
        ((if (and(eql (car (cdr lst)) a)(eql (car lst) b)) T))
        ((if (and(eql (car (cdr lst)) b)(eql (car lst) a)) T))
        (t(are-adjacent (cdr lst) a b)))  
)

; Transpose returns the transpose of the 2D matrix mat.

(defun app (a b)
    (cond ((null a) b)
          ( t  (cons (car a) 
                    (app (cdr a) b))   )
    )
) 

(defun transpose-helper (matrix newMatrix)
    (cond((eql (car matrix) NIL) newMatrix)
        ((eql newMatrix NIL) (transpose-helper (mapcar #'cdr matrix) (list (mapcar #'car matrix))))
             (t(transpose-helper (mapcar #'cdr matrix) (app newMatrix (list(mapcar #'car matrix))))))
)

(defun transpose (matrix)
    (cond((eql matrix NIL) NIL)
        (t(transpose-helper matrix NIL)))
)

; AreNeighbors returns true iff a and b are neighbors in the 2D
; matrix mat.

(defun are-neighbors-helper (matrix a b)
    (cond((eql matrix NIL) NIL)
        (t(are-adjacent (car matrix) a b)))  
)

(defun are-neighbors (matrix a b)
    (cond((eql matrix NIL) NIL)
        ((eql (are-neighbors-helper matrix a b) T) T)
        ((eql (are-neighbors-helper (transpose matrix) a b) T) T))
)
