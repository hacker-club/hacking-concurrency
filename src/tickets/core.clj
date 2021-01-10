(ns tickets.core
  (:require [tickets.input :as input])
  (:gen-class))

(def STANDARD_SEATS_CAPACITY 45)
(def PREMIUM_SEATS_CAPACITY 15)

(def seats (atom {:standard #{}
                  :premium #{}}))

(defn book [seats customer]
  (when (< (count (:standard @seats)) STANDARD_SEATS_CAPACITY)
    (swap! seats update-in [:standard] conj (:id customer))))

(defn offer-upgrade [seats customer]
  (when (and
         (:upgrades customer)
         (< (count (:premium @seats)) PREMIUM_SEATS_CAPACITY))
    (swap! seats update-in [:standard] disj (:id customer))
    (swap! seats update-in [:premium] conj (:id customer))))

(defn sell []
  (let [customers (input/read-json-file "input.json")
        seats (atom {:standard #{}
                     :premium #{}})]

    (doseq [customer customers]
      (book seats customer)
      (offer-upgrade seats customer))

    (println seats)
    (println "Booked " (count (:standard @seats)) "standard seats")
    (println "Booked " (count (:premium @seats)) "standard seats")))

(defn -main
  "I don't do a whole lot ... yet."
  [& args]
  (sell))
