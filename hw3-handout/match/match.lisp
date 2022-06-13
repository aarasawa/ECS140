; You may define helper functions here

(defun match-helper (pattern assertion)
    (cond((not (eql (car pattern) '!)) NIL)
        ((or (eql (car (cdr pattern)) (car assertion)) (eql (car (cdr pattern)) '?)) (match (cdr pattern) assertion))
        ((not (or (eql (car (cdr pattern)) (car assertion)) (eql (car (cdr pattern)) '?))) (match pattern assertion)))
    
)

(defun match (pattern assertion)
  (cond((eql pattern assertion) T)
      ((eql assertion NIL) NIL)
      ((eql (car pattern) (car assertion)) (match (cdr pattern) (cdr assertion)))
      ((eql (car pattern) '?) (match (cdr pattern) (cdr assertion)))
      ((and (eql (car pattern) '!) (eql (car (cdr pattern)) NIL)) T)
      (t(match-helper pattern (cdr assertion))))
)
