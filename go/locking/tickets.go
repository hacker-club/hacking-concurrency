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

type SafeSeats struct {
	mu       sync.Mutex
	Standard map[int]bool
	Premium  map[int]bool
}

func New() *SafeSeats {
	s := make(map[int]bool, StandardSeatsCapacity)
	p := make(map[int]bool, PremiumSeatsCapacity)
	return &SafeSeats{Standard: s, Premium: p}
}

func (s *SafeSeats) Book(c Customer) {
	s.mu.Lock()
	if len(s.Standard) < StandardSeatsCapacity {
		s.Standard[c.ID] = true
	}
	s.mu.Unlock()
}

func (s *SafeSeats) Upgrade(c Customer) {
	s.mu.Lock()
	if c.Upgrades && len(s.Premium) < PremiumSeatsCapacity && s.Standard[c.ID] {
		delete(s.Standard, c.ID)
		s.Premium[c.ID] = true
	}
	s.mu.Unlock()
}

func StartCashier(s *SafeSeats, queue chan Customer, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for customer := range queue {
		s.Book(customer)
		s.Upgrade(customer)
		if *debug {
			fmt.Printf("Customer %d attended by Cashier %d\n", customer.ID, id)
		}
	}
}

func main() {
	flag.Parse()

	var customers []Customer

	file, _ := os.Open("../input.json")
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &customers)

	safeSeats := New()

	queue := make(chan Customer)

	for i := 0; i < NumberOfCashiers; i++ {
		wg.Add(1)
		go StartCashier(safeSeats, queue, &wg, i)
	}

	for _, customer := range customers {
		queue <- customer
	}
	close(queue)
	wg.Wait()

	safeSeats.mu.Lock()
	fmt.Printf("Booked %d Standard seats\n", len(safeSeats.Standard))
	fmt.Printf("Booked %d Premium seats\n", len(safeSeats.Premium))
	fmt.Println(safeSeats)
	safeSeats.mu.Unlock()
}
