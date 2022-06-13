;; You may define helper functions here


(defun reachable-helper (transition liststart final input)
     (cond ((eql liststart NIL) NIL)
        ((if( eql (reachable transition (car liststart) final (cdr input)) T) T))
             (t(reachable-helper transition (cdr liststart) final input)))
 )

(defun reachable (transition start final input)
    ;; TODO: Incomplete function
    ;; The next line should not be in your solution.

    (cond ((eql input NIL)(cond((eql start final) T) (t NIL)))          ;; check if input is NIL, if start==final then output true
          (t (reachable-helper transition (funcall transition start (car input)) final input)))
          ;;((if(reachable (transition (car list) final (cdr input)))    ;; recursive call for next element in input


)

