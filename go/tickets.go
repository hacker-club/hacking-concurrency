package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	StandardSeatsCapacity int = 45
	PremiumSeatsCapacity  int = 15
)

type Customer struct {
	ID       int  `json:"id"`
	Upgrades bool `json:"upgrades"`
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

func main() {
	var customers []Customer

	file, _ := os.Open("input.json")
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &customers)

	seats := New()

	for _, customer := range customers {
		seats.Book(customer)
		seats.Upgrade(customer)
	}

	fmt.Printf("Booked %d Standard seats\n", len(seats.Standard))
	fmt.Printf("Booked %d Premium seats\n", len(seats.Premium))
	fmt.Println(seats)
}
