(ns tickets.input
  (:require [clojure.string :as s]
            [clojure.java.io :as io]
            [clojure.data.json :as json]))

(defn key->keyword [key-string]
  (-> key-string
      (s/replace #"([a-z])([A-Z])" "$1-$2")
      (s/replace #"([A-Z]+)([A-Z])" "$1-$2")
      (s/lower-case)
      (keyword)))

(defn parse-temp [key value]
  (if
   (= key :temp)
    (keyword value)
    value))

(defn read-json-message [is]
  (-> is
      (io/reader)
      (json/read :key-fn key->keyword
                 :value-fn parse-temp)))

(defn read-json-file [path]
  (-> path
      (io/resource)
      (read-json-message)))
