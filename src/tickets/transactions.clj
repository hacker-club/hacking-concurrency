(ns tickets.transactions
  (:require
   [clojure.core.async :refer [go <! >! >!! go-loop chan close! onto-chan]]
   [clojure.tools.logging :as log]
   [tickets.input :as input])
  (:gen-class))

(def STANDARD_SEATS_CAPACITY 45)
(def PREMIUM_SEATS_CAPACITY 15)

(defn book [seats customer]
  (if (< (count (:standard seats)) STANDARD_SEATS_CAPACITY)
    (update-in seats [:standard] conj (:id customer))
    seats))

(defn upgrade [seats customer]
  (if
   (and (:upgrades customer)
        (< (count (:premium seats)) PREMIUM_SEATS_CAPACITY))
    (do (update-in seats [:standard] disj (:id customer))
        (update-in seats [:premium] conj (:id customer)))
    seats))

(defn start-cashier [id queue seats]
  (go-loop []
    (let [customer (<! queue)]
      (when customer
        (dosync
         (alter seats book customer)
         (alter seats upgrade customer)
         (println (format "Customer %d was attended by Cashier %d" (:id customer) id))
         (log/info (format "Customer %d was attended by Cashier %d" (:id customer) id)))
        (recur)))))

(defn open-box-office [number-of-cashiers]
  (let [customers (input/read-json-file "input.json")
        seats (ref {:standard #{}
                    :premium #{}})
        queue (chan)]

    (dotimes [i number-of-cashiers]
      (start-cashier i queue seats))

    (doseq [customer customers]
      (>!! queue customer))

    (println "Booked " (count (:standard @seats)) "standard seats")
    (println "Booked " (count (:premium @seats)) "premium seats")))

(defn -main
  [& args]
  (log/info "hello world")
  (open-box-office 3))
