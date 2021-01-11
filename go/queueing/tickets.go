package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var wg sync.WaitGroup
var debug = flag.Bool("debug", false, "Should print debug output")

const (
	StandardSeatsCapacity int = 45
	PremiumSeatsCapacity  int = 15
	NumberOfCashiers      int = 3
)

type Customer struct {
	ID       int  `json:"id"`
	Upgrades bool `json:"Upgrades"`
}

type Seats struct {
	Standard map[int]bool
	Premium  map[int]bool
}

func New() *Seats {
	s := make(map[int]bool, StandardSeatsCapacity)
	p := make(map[int]bool, PremiumSeatsCapacity)
	return &Seats{Standard: s, Premium: p}
}

func (s *Seats) Book(c Customer) {
	if len(s.Standard) < StandardSeatsCapacity {
		s.Standard[c.ID] = true
	}
}

func (s *Seats) Upgrade(c Customer) {
	if c.Upgrades && len(s.Premium) < PremiumSeatsCapacity && s.Standard[c.ID] {
		delete(s.Standard, c.ID)
		s.Premium[c.ID] = true
	}
}

func StartCashier(s *Seats, queue chan Customer, updates chan Customer, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for customer := range queue {
		updates <- customer
		if *debug {
			fmt.Printf("Customer %d attended by Cashier %d\n", customer.ID, id)
		}
	}
}

func StartUpdater(s *Seats, updates chan Customer) {
	for customer := range updates {
		if *debug {
			fmt.Println("Processing customer", customer)
		}
		s.Book(customer)
		s.Upgrade(customer)
	}
}

func main() {
	var customers []Customer

	flag.Parse()

	file, _ := os.Open("../input.json")
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &customers)

	seats := New()

	queue := make(chan Customer)
	updates := make(chan Customer)

	for i := 0; i < NumberOfCashiers; i++ {
		wg.Add(1)
		go StartCashier(seats, queue, updates, &wg, i)
	}

	go StartUpdater(seats, updates)

	for _, customer := range customers {
		queue <- customer
	}
	close(queue)

	wg.Wait()

	fmt.Printf("Booked %d Standard seats\n", len(seats.Standard))
	fmt.Printf("Booked %d Premium seats\n", len(seats.Premium))
	fmt.Println(seats)
}
