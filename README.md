# Hacking concurrency programming with an example

This repo holds solutions to the hacking concurrency post in [hackerclub](https://hackerclub.io/hacking-concurrency/).

Go code implements a single threaded solution, a locking solution and a queueing solution.

Clojure code implements a transactional based solution.

## Running

Clojure  
Leiningen is a dependency, please install in [here](https://leiningen.org/#install).

```shell
$ cd clj
$ lein run
```

Go
```shell
$ cd go
$ go run tickets.go

$ cd locking
$ go run tickets.go

$ cd queueing
$ go run tickets.go
```