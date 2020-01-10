package Canales

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Phases     = 20
	Goroutines = 10
)

func run(id int, done chan bool, b *Barrier) {
	for i := 0; i < Phases; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		b.Barrier()
		fmt.Printf("%d finished phase %d\n", id, i)
	}
	b.Barrier()
	fmt.Printf("Finished thread %d\n", id)
	done <- true
}

func main() {
	done := make(chan bool, 1)
	barrier := NewBarrier(Goroutines)
	for i := 0; i < Goroutines; i++ {
		go run(i, done, barrier)
	}
	for i := 0; i < Goroutines; i++ {
		<-done
	}
	fmt.Println("End")
}

type Barrier struct {
	arrival   chan int
	departure chan int
	n         int
}

func NewBarrier(value int) *Barrier {
	b := new(Barrier)
	b.arrival = make(chan int, 1)
	b.departure = make(chan int, 1)
	b.n = value
	b.arrival <- value
	return b
}

func (b *Barrier) Barrier() {
	var v int
	v = <-b.arrival
	if v > 1 {
		v--
		b.arrival <- v
	} else {
		b.departure <- b.n
	}
	v = <-b.departure
	if v > 1 {
		v--
		b.departure <- v
	} else {
		b.arrival <- b.n
	}
}
