(ns mal.step0-repl
  #?(:clj (:gen-class)))

(defn READ [string]
  string)

(defn EVAL [ast env]
  ast)

(defn PRINT [exp]
  exp)

(defn rep [string]
  (PRINT (EVAL (READ string) nil)))

(defn prompt [ps]
  (do (print ps)
      (flush)
      (read-line)))

(defn repl []
  (let [line (prompt "user> ")]
    (if (nil? line)
      nil
      (do (println (rep line))
          (recur)))))

(defn -main [& args]
  (repl))
