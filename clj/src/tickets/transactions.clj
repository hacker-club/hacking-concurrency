(ns tickets.transactions
  (:require
   [clojure.core.async :refer [go <! >! >!! go-loop chan close! onto-chan]]
   [clojure.tools.logging :as log]
   [tickets.input :as input])
  (:gen-class))

(def STANDARD_SEATS_CAPACITY 45)
(def PREMIUM_SEATS_CAPACITY 15)

(defn book [standard customer]
  (if (< (count @standard) STANDARD_SEATS_CAPACITY)
    (alter standard conj (:id customer))
    standard))

(defn upgrade [standard premium customer]
  (if
   (and (:upgrades customer)
        (< (count @premium) PREMIUM_SEATS_CAPACITY)
        (contains? @standard (:id customer)))
    (do
      (alter standard disj (:id customer))
      (alter premium conj (:id customer)))
    [standard premium]))

(defn start-cashier [id queue standard premium]
  (go-loop []
    (let [customer (<! queue)]
      (when customer
        (dosync
         (log/info (format "Customer %d was attended by Cashier %d" (:id customer) id))
         (book standard customer)
         (upgrade standard premium customer))
        (recur)))))

(defn open-box-office [number-of-cashiers]
  (let [customers (input/read-json-file "input.json")
        standard (ref #{})
        premium (ref #{})
        queue (chan)]

    (dotimes [i number-of-cashiers]
      (start-cashier i queue standard premium))

    (doseq [customer customers]
      (>!! queue customer))

    (println "Booked " (count @standard) "standard seats")
    (println "Booked " (count @premium) "premium seats")
    (println @standard)
    (println @premium)))

(defn -main
  [& args]
  (open-box-office 3))
